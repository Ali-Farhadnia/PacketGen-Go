# PacketGen-Go

PacketGen-Go is a powerful command-line Ethernet packet generator written in Go. It supports custom packet crafting, multiple protocols, various traffic patterns, and real-time monitoring of network traffic generation.

## Features

- **Packet Crafting**
  - Custom MAC addresses
  - Configurable EtherType
  - Variable payload sizes
  - VLAN tagging support

- **Protocol Support**
  - IPv4
  - IPv6
  - ARP
  - TCP
  - UDP
  - VLAN
  - Custom protocols via EtherType

- **Traffic Generation**
  - Adjustable packet rates
  - Configurable packet counts
  - Continuous mode
  - Multiple traffic patterns (sequential, random, burst)

- **Interface Management**
  - List available network interfaces
  - Interface selection
  - Interface details (MAC, IP addresses)

- **Real-time Monitoring**
  - Packets sent counter
  - Current transmission rate
  - Bandwidth utilization
  - Error tracking

## Prerequisites

Before installing PacketGen-Go, ensure you have the following prerequisites:

1. Go 1.16 or higher
2. libpcap development headers

### Installing libpcap

For Ubuntu/Debian:
```bash
sudo apt-get install libpcap-dev
```

For RedHat/Fedora/CentOS:
```bash
sudo dnf install libpcap-devel
```

For macOS:
```bash
brew install libpcap
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Ali-Farhadnia/PacketGen-Go.git
cd PacketGen-Go
```

2. Install dependencies:
```bash
go get github.com/google/gopacket
```

3. Build the project:
```bash
go build
```

## Usage

### List Available Interfaces
```bash
sudo ./PacketGen-Go --list
```

### Basic Packet Generation
```bash
sudo ./PacketGen-Go -i eth0 -src-mac 00:11:22:33:44:55 -dst-mac 66:77:88:99:AA:BB -rate 1000 -count 10000
```

### Command Line Options

```
-i           : Network interface name (required)
-list        : List available network interfaces
-src-mac     : Source MAC address
-dst-mac     : Destination MAC address
-ether-type  : EtherType value (default: IPv4)
-payload-size: Payload size in bytes (default: 46)
-protocol    : Protocol (ipv4, ipv6, arp, tcp, udp)
-rate        : Packets per second (default: 1000)
-count       : Number of packets to send (0 for continuous)
-vlan        : VLAN ID (-1 to disable)
-mode        : Traffic mode (sequential, random, burst)
```

### Examples

Generate VLAN tagged traffic:
```bash
sudo ./PacketGen-Go -i eth0 -src-mac 00:11:22:33:44:55 -dst-mac 66:77:88:99:AA:BB -vlan 100
```

Generate random traffic:
```bash
sudo ./PacketGen-Go -i eth0 -src-mac 00:11:22:33:44:55 -dst-mac 66:77:88:99:AA:BB -mode random
```

Continuous traffic generation:
```bash
sudo ./PacketGen-Go -i eth0 -src-mac 00:11:22:33:44:55 -dst-mac 66:77:88:99:AA:BB -count 0
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [gopacket](https://github.com/google/gopacket) library for packet manipulation
- libpcap for packet capture and transmission

## Security Considerations

Please note that this tool can generate significant network traffic. Use responsibly and only on networks you own or have permission to test. Misuse of this tool could violate laws or network policies.
