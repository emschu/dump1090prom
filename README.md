# Dump1090Prom – Aircraft data exporter for Prometheus

A Prometheus exporter for `readsb`/`dump1090` aircraft data.
This exporter converts dump1090 data into [Prometheus](https://prometheus.io/) metrics, allowing you to monitor and visualize collected aircraft data in detail.
The repository includes a [Grafana dashboard](dashboards/dump1090prom-grafana-dashboard.json) for visualizing the metrics.

Tested with:
- [flightaware/dump1090](https://github.com/flightaware/dump1090)
- [wiedehopf/readsb](https://github.com/wiedehopf/readsb)

Version: 1.0.4

## Features

- Exports aircraft data as Prometheus metrics.
- Supports multiple data sources:
  - Local directory containing `aircraft.json` and `receiver.json`.
  - HTTP/HTTPS base URL exposing `aircraft.json` and `receiver.json`.
- Enriches aircraft data with airline information (embedded from Wikipedia).
- Calculates the distance between the receiver and aircraft when the receiver position is known (from `receiver.json` or `-lat`/`-lon` override).
- Includes a ready-to-use Grafana dashboard.

## Installation

### Prerequisites

- `readsb` or `dump1090` (with `rtlsdr` or compatible hardware) generating `aircraft.json` and `receiver.json`.

### Option 1: Download a release file

Download the latest release from the [releases page](https://github.com/emschu/dump1090prom/releases).

### Option 2: Install via go

```bash
go install github.com/emschu/dump1090prom@latest
```

### Option 3: Building from source

```bash
git clone https://github.com/emschu/dump1090prom.git
cd dump1090prom
go build
```

### Running permanently (Linux)

If you use Linux with systemd, you can use [this systemd unit file](./dev/dump1090prom.service) and copy it
to `/etc/systemd/system/dump1090prom.service`. You may want to adjust the `port` and add custom parameters in `ExecStart`.
Look for the Configuration and CLI section below for more information.

Run the following commands to enable the exporter service permanently in a Linux system using systemd:

```bash
sudo systemctl daemon-reload
sudo systemctl enable dump1090prom 
sudo systemctl start dump1090prom 
```

#### Build for other platforms

```bash
# e.g. for Raspberry Pi
env GOOS=linux GOARCH=arm go build -o dump1090prom
```

## Configuration

The exporter is configured via command-line flags:

- `-base-url` `string`
  Base URL to the directory where aircraft.json and receiver.json are available
- `-base-path` `string`
  Local filesystem directory containing aircraft.json and receiver.json, e.g. /var/www/html/data
- `-lat` `float`
  Override receiver latitude (used for distance calculation)
- `-lon` `float`
  Override receiver longitude (used for distance calculation)
- `-host` `host` (default: `127.0.0.1`)
  Override host where /metrics is exposed
- `-port` `port` (default: `8080`)
  Override port where /metrics is exposed
- `-verbose` `bool` (default: `false`)
  Enable verbose logging
- `-distance-calc` `bool` (default: `true`)
  Enable distance calculation to aircraft
- `-airline-label` `bool` (default: `true`)
  Enable airline labelling
- `-expose-files` `bool` (default: `true`)
  Expose original aircraft.json and receiver.json files at /aircraft.json and /receiver.json
- `-rolling-map-size` `int` (default: `1000`)
  Default size of the rolling map for caching
- `-collector` `string` (default: `dump1090prom`)
  Constant 'collector' label value of this instance
- `-labels` `string`
  Global labels to add to all metrics (e.g., 'location=home,environment=testing')

Notes:
- When `-lat` and `-lon` are not provided, the exporter will try to read receiver position from `receiver.json`.
- Distance calculation is performed only when a receiver position is known.
- Provide either `-base-url` or `-base-path` (not both).

## Usage Examples
### Starting `dump1090prom`

Use a local or remotely accessible dump1090 instance:

```bash
./dump1090prom -base-url http://your-dump1090-server:8080/data
# for example for dump1090
./dump1090prom -base-url http://127.0.0.1:8080/data
```

Or without HTTP:

```bash
./dump1090prom -base-path /path/to/dump1090/data
# for example for readsb
./dump1090prom -base-path /run/readsb
```

Show verbose logs:
```bash
./dump1090prom -base-url http://your-dump1090-server:8080/data -verbose
```

Expose metrics on a different port:
```bash
./dump1090prom -base-url http://your-dump1090-server:8080/data -port 9091
```

Add global labels to all metrics:
```bash
./dump1090prom -base-url http://your-dump1090-server:8080/data -labels "environment=home"
```

Typically, `dump1090prom` uses the `<lat>` and `<lon>` values defined in your `receiver.json`. To override the receiver position, use the `-lat` and `-lon` flags.

```bash
./dump1090prom -base-url http://your-dump1090-server:8080/data -lat <lat> -lon <lon>
```

## Prometheus Configuration

Note: Replace `<host_ip>` and port with the IP address of your `dump1090prom` server.

```yaml
    scrape_configs:
      - job_name: 'dump1090prom'
        static_configs:
          - targets: ['<host_ip>:8080']
        metrics_path: /metrics
        scrape_interval: 1s
```

## Grafana Dashboard

The project includes an example Grafana dashboard in the `dashboards` directory. 

To use it, import the [`dashboards/dump1090prom-grafana-dashboard.json`](dashboards/dump1090prom-grafana-dashboard.json) JSON definition into your Grafana instance.

![Grafana Dashboard 1](./example/grafana_dashboard_1.png)
![Grafana Dashboard 2](./example/grafana_dashboard_4.png)

## Metrics

### Accessing metrics

```
http://localhost:8080/metrics
http://localhost:8080/aircraft.json
http://localhost:8080/receiver.json
```

#### Metric names

All metrics are prefixed with `dump1090prom_`.

| Metric Name | Type | Description |
| ----------- | ---- | ----------- |
| `aircraft_flight_info` | Gauge | Metadata about the flight and aircraft (includes airline, country, squawk, etc.). |
| `aircraft_altitude_baro_feet` | Gauge | Barometric altitude in feet. |
| `aircraft_ground_speed_knots` | Gauge | Ground speed in knots. |
| `aircraft_distance_from_position_meters` | Gauge | Distance from the receiver in meters. |
| `aircraft_count` | Gauge | Total number of aircraft currently seen. |
| `total_messages` | Counter | Total number of messages received. |

Additional metrics are available for ADS-B version, vertical rates, navigation data, and more. See [metric.go](./metric.go) for a full list or visit the `/metrics` endpoint.

## Development

### Testing

```bash
go test ./...
```

### Development Environment

The `dev` directory contains resources for setting up a development environment:

- `dev/prometheus`: Prometheus configuration and Podman Compose setup

This project is licensed under the terms of the [GNU Affero General Public License v3.0](./LICENSE).
Wikipedia airline data is licensed under [CC-BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
