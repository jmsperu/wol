# wol

Wake-on-LAN CLI tool. Send magic packets to wake devices on your network. Save devices by name, check their online status, and batch wake multiple devices. Single binary, no dependencies.

Works on **Windows**, **macOS**, and **Linux**.

## Features

- **Send magic packets** -- wake devices by MAC address
- **Named devices** -- save devices with friendly names for quick access
- **Batch wake** -- wake multiple devices in one command
- **SecureOn password** -- supports WOL password authentication
- **Status checks** -- verify which saved devices are online (TCP probe)
- **Wait mode** -- send packet and poll until the device comes online
- **Custom broadcast** -- specify broadcast address and UDP port

## Install

### Binary download

Grab the latest binary from [Releases](https://github.com/jmsperu/wol/releases) and place it in your `PATH`.

### Go install

```bash
go install github.com/jmsperu/wol@latest
```

### Build from source

```bash
git clone https://github.com/jmsperu/wol.git
cd wol
make build
```

Cross-compile for all platforms:

```bash
make build-all    # outputs to dist/
```

## Quick start

```bash
wol AA:BB:CC:DD:EE:FF                  # wake by MAC address
wol add myserver AA:BB:CC:DD:EE:FF     # save a device
wol myserver                           # wake by name
```

## Commands

### Wake by MAC or name

```bash
wol AA:BB:CC:DD:EE:FF                         # wake by MAC address
wol myserver                                   # wake a saved device
wol AA:BB:CC:DD:EE:FF -b 192.168.1.255         # specify broadcast address
wol AA:BB:CC:DD:EE:FF --port 7                 # use UDP port 7
wol AA:BB:CC:DD:EE:FF -p AA:BB:CC:DD:EE:FF     # with SecureOn password
```

### `wol wake`

Wake one or more devices (by name or MAC). Supports batch waking.

```bash
wol wake myserver
wol wake myserver mynas mydesktop              # batch wake
wol wake AA:BB:CC:DD:EE:FF                     # by MAC
wol wake myserver -w                           # wake and wait until online
```

### `wol add`

Save a device for quick access.

```bash
wol add myserver AA:BB:CC:DD:EE:FF
wol add myserver AA:BB:CC:DD:EE:FF -i 192.168.1.100                  # with IP for status checks
wol add myserver AA:BB:CC:DD:EE:FF -i 192.168.1.100 -b 192.168.1.255 # with broadcast
wol add myserver AA:BB:CC:DD:EE:FF -p AA:BB:CC:DD:EE:FF              # with SecureOn password
```

### `wol list` (alias: `ls`)

List all saved devices.

```bash
wol list
```

### `wol remove` (aliases: `rm`, `delete`)

Remove a saved device.

```bash
wol remove myserver
```

### `wol status`

Check which saved devices are online (probes TCP ports 22, 80, 443, 3389).

```bash
wol status
```

## Flags reference

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--broadcast` | `-b` | `255.255.255.255` | Broadcast address |
| `--password` | `-p` | | SecureOn password (6 hex bytes, e.g. `AA:BB:CC:DD:EE:FF`) |
| `--port` | | `9` | UDP port |
| `--wait` | `-w` | `false` | Wait and check if device comes online after wake |

### `wol add` flags

| Flag | Short | Description |
|------|-------|-------------|
| `--ip` | `-i` | IP address for status checks |
| `--broadcast` | `-b` | Broadcast address for this device |
| `--password` | `-p` | SecureOn password for this device |

## Config file

Devices are saved to `~/.wol.yml`:

```yaml
devices:
  myserver:
    mac: "AA:BB:CC:DD:EE:FF"
    ip: "192.168.1.100"
    broadcast: "192.168.1.255"
  mynas:
    mac: "11:22:33:44:55:66"
    ip: "192.168.1.200"
```

## License

MIT
