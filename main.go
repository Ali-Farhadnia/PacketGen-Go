package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type PacketConfig struct {
	Interface   string
	SrcMAC      string
	DstMAC      string
	EtherType   uint16
	PayloadSize int
	Protocol    string
	Rate        int
	Count       int
	Continuous  bool
	TrafficMode string
	VLAN        *int
	IPv4        *layers.IPv4
	IPv6        *layers.IPv6
	TCP         *layers.TCP
	UDP         *layers.UDP
}

type Statistics struct {
	PacketsSent uint64
	BytesSent   uint64
	Errors      uint64
	StartTime   time.Time
	mutex       sync.Mutex
}

func (s *Statistics) update(bytes uint64, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err != nil {
		s.Errors++
		return
	}

	s.PacketsSent++
	s.BytesSent += bytes
}

func (s *Statistics) display() {
	duration := time.Since(s.StartTime).Seconds()
	pps := float64(s.PacketsSent) / duration
	bps := float64(s.BytesSent) * 8 / duration

	fmt.Printf("\rPackets: %d | Rate: %.2f pps | Bandwidth: %.2f Mbps | Errors: %d",
		s.PacketsSent, pps, bps/1000000, s.Errors)
}

func listInterfaces() {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available network interfaces:")
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		addresses := make([]string, 0)
		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}

		fmt.Printf("[%d] %s\n", iface.Index, iface.Name)
		fmt.Printf("    MAC: %s\n", iface.HardwareAddr)
		fmt.Printf("    Addresses: %s\n", strings.Join(addresses, ", "))
	}
}

func createPacket(config PacketConfig) ([]byte, error) {
	srcMAC, err := net.ParseMAC(config.SrcMAC)
	if err != nil {
		return nil, fmt.Errorf("invalid source MAC: %v", err)
	}

	dstMAC, err := net.ParseMAC(config.DstMAC)
	if err != nil {
		return nil, fmt.Errorf("invalid destination MAC: %v", err)
	}

	// Create ethernet layer
	eth := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetType(config.EtherType),
	}

	// Create buffer and serialize layers
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	payloadData := make([]byte, config.PayloadSize)

	// Add VLAN tag if specified
	if config.VLAN != nil {
		dot1q := &layers.Dot1Q{
			Priority:       0,
			DropEligible:   false,
			VLANIdentifier: uint16(*config.VLAN),
			Type:           eth.EthernetType,
		}
		eth.EthernetType = layers.EthernetTypeDot1Q

		if err := gopacket.SerializeLayers(buffer, opts,
			&eth, dot1q, gopacket.Payload(payloadData)); err != nil {
			return nil, err
		}
	} else {
		if err := gopacket.SerializeLayers(buffer, opts,
			&eth, gopacket.Payload(payloadData)); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func generateTraffic(handle *pcap.Handle, config PacketConfig, stats *Statistics) {
	packet, err := createPacket(config)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second / time.Duration(config.Rate))
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		switch config.TrafficMode {
		case "sequential":
			err = handle.WritePacketData(packet)
		case "random":
			// Modify payload randomly
			randPacket := make([]byte, len(packet))
			copy(randPacket, packet)
			rand.Read(randPacket[14:]) // Skip ethernet header
			err = handle.WritePacketData(randPacket)
		case "burst":
			// Send burst of 10 packets
			for i := 0; i < 10; i++ {
				err = handle.WritePacketData(packet)
				stats.update(uint64(len(packet)), err)
			}
			continue
		}

		stats.update(uint64(len(packet)), err)

		if !config.Continuous {
			count++
			if count >= config.Count {
				break
			}
		}
	}
}

func main() {
	// Parse command line flags
	interfaceName := flag.String("i", "", "Network interface name")
	listInterfacesFlag := flag.Bool("list", false, "List available network interfaces")
	srcMAC := flag.String("src-mac", "", "Source MAC address")
	dstMAC := flag.String("dst-mac", "", "Destination MAC address")
	etherType := flag.Int("ether-type", int(layers.EthernetTypeIPv4), "EtherType value")
	payloadSize := flag.Int("payload-size", 46, "Payload size in bytes")
	protocol := flag.String("protocol", "ipv4", "Protocol (ipv4, ipv6, arp, tcp, udp)")
	rate := flag.Int("rate", 1000, "Packets per second")
	count := flag.Int("count", 1000, "Number of packets to send (0 for continuous)")
	vlan := flag.Int("vlan", -1, "VLAN ID (-1 to disable)")
	trafficMode := flag.String("mode", "sequential", "Traffic mode (sequential, random, burst)")

	flag.Parse()

	if *listInterfacesFlag {
		listInterfaces()
		os.Exit(0)
	}

	if *interfaceName == "" {
		log.Fatal("Interface name is required")
	}

	// Create packet configuration
	config := PacketConfig{
		Interface:   *interfaceName,
		SrcMAC:      *srcMAC,
		DstMAC:      *dstMAC,
		EtherType:   uint16(*etherType),
		PayloadSize: *payloadSize,
		Protocol:    *protocol,
		Rate:        *rate,
		Count:       *count,
		Continuous:  *count == 0,
		TrafficMode: *trafficMode,
	}

	if *vlan != -1 {
		config.VLAN = vlan
	}

	// Open network interface
	handle, err := pcap.OpenLive(config.Interface, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Initialize statistics
	stats := &Statistics{StartTime: time.Now()}

	// Start statistics display goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			stats.display()
		}
	}()

	// Generate traffic
	generateTraffic(handle, config, stats)

	fmt.Println("\nTraffic generation completed")
}
