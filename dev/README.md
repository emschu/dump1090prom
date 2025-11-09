# Development

## Prometheus Development Docker Setup

Look in [`prometheus/README.md`](./prometheus/README.md).


## Systemd example Service


#### The following is an example of a systemd service on a raspberry pi
```
[Unit]
Description=Prometheus Exporter for dump1090/readsb data
Documentation=https://github.com/emschu/dump1090prom
After=network.target
Wants=network.target

[Service]
Type=simple
User=pi
Group=pi

ExecStart=/usr/local/bin/dump1090prom --base-path /run/readsb --verbose --port 8081

NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/run/readsb

Restart=on-failure
RestartSec=5

StandardOutput=journal
StandardError=journal
SyslogIdentifier=dump1090prom

[Install]
WantedBy=multi-user.target
```

## Experimental PromQL Queries
```promql
rate((dump1090prom_aircraft_altitude_baro_feet or dump1090prom_aircraft_altitude_geom_feet)[5])

rate(dump1090prom_aircraft_altitude_baro_feet[1m])

group(dump1090prom_aircraft_flight_info == 1) by (flight)
group(dump1090prom_aircraft_flight_info == 1) by (flight, hex)
count by (flight) (dump1090prom_aircraft_flight_info)
dump1090prom_aircraft_true_airspeed_knots or dump1090prom_aircraft_ground_speed_knots
group(dump1090prom_aircraft_flight_info == 1) by (flight, hex, lat, lon)
```
