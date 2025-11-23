# transmission-gluetun-port update

Small go binary to periodically update the port of transmission based on gluten port forwarding status.

Based on [transmission-nat-pmp](https://github.com/jordanpotter/transmission-nat-pmp).

## Usage

```bash
transmission-gluetun-port -h
```

## Available environment variables

| Name                    | Description                                                                        | Default                     |
|-------------------------|------------------------------------------------------------------------------------|-----------------------------|
| `TRANSMISSION_USER`     | Transmission user                                                                  | -                           |
| `TRANSMISSION_PASSWORD` | Transmission password                                                              | -                           |
| `TRANSMISSION_PROTOCOL` | Transmission api protocol: `http`, `https`                                         | `http`                      |
| `GLUETUN_PROTOCOL`      | Gluetun api protocol: `http`, `https`                                              | `http`                      |
| `GLUETUN_HOSTNAME`      | Gluetun api hostname                                                               | `127.0.0.1`                 |
| `GLUETUN_PORT`          | Gluetun api port                                                                   | `8000`                      |
| `GLUETUN_ENDPOINT`      | Gluetun api current port endpoint                                                  | `/v1/openvpn/portforwarded` |
| `GLUETUN_AUTH_TYPE`     | Gluetun auth type: `basic`, `apikey`                                               | `none`                      |
| `GLUETUN_AUTH_USERNAME` | Gluetun basic auth username                                                        | -                           |
| `GLUETUN_AUTH_PASSWORD` | Gluetun basic auth password                                                        | -                           |
| `GLUETUN_AUTH_API_KEY`  | Gluetun auth api key                                                               | -                           |
| `INITIAL_DELAY`         | Initial delay ([format](https://pkg.go.dev/time#ParseDuration))                    | `5s`                        |
| `CHECK_INTERVAL`        | Update interval ([format](https://pkg.go.dev/time#ParseDuration))                  | `1m`                        |
| `ERROR_INTERVAL`        | Update interval in case of error ([format](https://pkg.go.dev/time#ParseDuration)) | `5s`                        |
