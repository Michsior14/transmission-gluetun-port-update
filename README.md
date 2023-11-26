# transmission-gluetun-port update

Small go binary to periodically update the port of transmission based on gluten port forwarding status.

Based on [transmission-nat-pmp](https://github.com/jordanpotter/transmission-nat-pmp).

## Usage

```bash
transmission-gluetun-port -h
```

## Available environment variables

| Name | Description | Default |
|------|-------------|---------|
| `TRANSMISSION_USER` | Transmission user | - |
| `TRANSMISSION_PASSWORD` | Transmission password | - |
| `GLUETUN_HOST` | Gluetun api host | `127.0.0.1` |
| `GLUETUN_PORT` | Gluetun api port | `8000` |
| `CHECK_INTERVAL` | Update interval ([format](https://pkg.go.dev/time#ParseDuration)) | `1m` |
| `ERROR_INTERVAL` | Update interval in case of error ([format](https://pkg.go.dev/time#ParseDuration)) | `5s` |
