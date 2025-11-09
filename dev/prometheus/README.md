# Local Monitoring Development Environment

This directory contains the necessary files to run Prometheus and Grafana in a local development environment using Podman.

## Files

- `prometheus.yml`: Prometheus configuration file
- `podman-compose.yml`: Podman Compose file to run Prometheus and Grafana

## Prerequisites

- [Podman](https://podman.io/getting-started/installation)
- [Podman Compose](https://github.com/containers/podman-compose)

## Usage

1. Make sure the `dump1090prom` application is running on the host machine on port 8080.

2. Start the monitoring stack using Podman Compose:

```bash
cd dev/prometheus
podman-compose up -d
```

3. Access the monitoring interfaces:
   - Prometheus web interface: http://localhost:9090
   - Grafana dashboard: http://localhost:3000 (default login: admin/admin)

4. To stop the monitoring stack:

```bash
podman-compose down
```

## Configuration

### Prometheus Configuration

The `prometheus.yml` file contains the Prometheus configuration. By default, it is configured to:

- Scrape metrics from the dump1090prom application at http://host.containers.internal:8080/metrics every second
- Scrape Prometheus's own metrics at http://localhost:9090

You can modify this file to change the scrape configuration or add additional targets.

### Grafana Configuration

Grafana is configured with the following default settings:

- Default admin password: `admin` (you'll be prompted to change it on first login)
- User sign-up is disabled
- Server domain is set to localhost

### Volume Mounts

The Podman Compose configuration includes the following volume mounts:

- Prometheus:
  - `./prometheus.yml:/etc/prometheus/prometheus.yml`: Mounts the local Prometheus configuration file into the container
  - `prometheus_data:/prometheus`: Mounts a named volume for persistent storage of Prometheus data
- Grafana:
  - `grafana_data:/var/lib/grafana`: Mounts a named volume for persistent storage of Grafana data

### Setting up Prometheus as a Data Source in Grafana

After starting the stack, follow these steps to configure Prometheus as a data source in Grafana:

1. Log in to Grafana at http://localhost:3000 (default credentials: admin/admin)
2. Go to Configuration > Data Sources
3. Click "Add data source" and select "Prometheus"
4. Set the URL to `http://prometheus:9090` (using the service name within the shared network)
5. Click "Save & Test" to verify the connection

## Customization

To customize the configuration:

1. Edit the `prometheus.yml` file or modify the `podman-compose.yml` file
2. Restart the stack:

```bash
podman-compose restart
```

## Troubleshooting

If you encounter issues with the host.containers.internal DNS name not resolving, make sure your Podman version supports the host-gateway feature. Alternatively, you can replace host.containers.internal with the actual IP address of your host machine.