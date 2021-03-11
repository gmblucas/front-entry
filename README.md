# front-entry

Lightweight HTTP(s) reverse proxy

## Usage

```bash
$ front-entry config.toml
```

## Config

```toml
[proxy]
"rpc.gambinode.com" = "http://localhost:7072"
"ws.gambinode.com" = "http://localhost:7074"
"monitor.gambinode.com" = "http://localhost:8080"
"netdata.gambinode.com" = "http://localhost:19999"
"152.228.141.68" = "http://localhost:8080"

[tls]
"certfile" = "fullchain.pem"
"keyfile" = "privkey.pem"
```

A go routine watches for config file modifications & reload at runtime the map of source <-> destination
