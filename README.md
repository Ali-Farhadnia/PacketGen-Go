# PacketGen-Go

PacketGen-Go is a simple Go-based packet generator tool that allows you to send custom Ethernet packets over a network interface. It is designed for testing and debugging network configurations, performance, and behavior.

## Features

- Send custom Ethernet packets with specified source and destination MAC addresses.
- Configure EtherType (e.g., IPv4, ARP, etc.).
- Set a custom payload for the packets.
- Control the packet rate (packets per second).
- Option to send packets continuously or for a specified count.

## Installation

1. **Install Go**: Ensure you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

2. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/PacketGen-Go.git
   cd PacketGen-Go
   ```

3. **Install dependencies**:
   ```bash
   go get github.com/google/gopacket
   ```

4. **Build the project**:
   ```bash
   go build
   ```

## Usage

Run the `PacketGen-Go` tool with the following command-line options:

```bash
./PacketGen-Go -i <interface> -c <count> -r <rate> -srcmac <source-mac> -dstmac <destination-mac> -ethertype <ethertype> -payload <payload> -cont
```

### Command-Line Options

| Option       | Description                                      | Default Value          |
|--------------|--------------------------------------------------|------------------------|
| `-i`         | Network interface to send packets on             | `eth0`                 |
| `-c`         | Number of packets to send                        | `10`                   |
| `-r`         | Packets per second                               | `1000`                 |
| `-cont`      | Continuous mode (send packets indefinitely)      | `false`                |
| `-srcmac`    | Source MAC address                               | `00:00:00:00:00:01`    |
| `-dstmac`    | Destination MAC address                          | `00:00:00:00:00:02`    |
| `-ethertype` | EtherType (e.g., `0x0800` for IPv4)              | `0x0800`               |
| `-payload`   | Packet payload                                   | `"Hello, World!"`      |

### Example

Send 100 packets at a rate of 500 packets per second on interface `eth0`:

```bash
./PacketGen-Go -i eth0 -c 100 -r 500 -srcmac 00:11:22:33:44:55 -dstmac 66:77:88:99:AA:BB -ethertype 0x0800 -payload "Test Packet"
```

Send packets continuously:

```bash
./PacketGen-Go -i eth0 -cont -r 1000 -srcmac 00:11:22:33:44:55 -dstmac 66:77:88:99:AA:BB -ethertype 0x0800 -payload "Continuous Test"
```

## Contributing

Contributions are welcome! If you find a bug or have a feature request, please open an issue. If you'd like to contribute code, feel free to submit a pull request.

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a pull request.

```

### Notes:
- Replace `your-username` in the clone URL with your actual GitHub username.
- Add a `LICENSE` file if you want to include one (e.g., MIT License).
- Customize the `README.md` further if you have additional features or instructions.

Let me know if you need further assistance!
