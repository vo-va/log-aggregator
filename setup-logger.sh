#!/bin/bash
set -e


if [ $# -lt 3 ]; then
	echo "Usage: $ <firehose_stream> <product> <component> <units/service to ignore>"
	echo "example: $0 log-stream \"product a\" \"component b\" \"redis.service,postgresql\""
	echo "WARNING logger will not be installed, safely exiting with 0 to not abort deployment"

	exit 0
fi

eleven_product="$2"
eleven_component="$3"

ignore_systemd_units="${4:-}"

# Install logger to /usr/bin
cp ./log-aggregator /usr/local/bin/log-aggregator
chmod +x /usr/local/bin/log-aggregator

endpoint="169.254.169.254"

cat <<EOF >/usr/local/bin/start-logger
#!/bin/bash
set -e
export EC2_METADATA_INSTANCE_ID=$(curl http://$endpoint/latest/meta-data/instance-id)
export EC2_METADATA_LOCAL_IPV4=$(curl http://$endpoint/latest/meta-data/local-ipv4)
export EC2_METADATA_LOCAL_HOSTNAME=$(curl http://$endpoint/latest/meta-data/local-hostname)
/usr/local/bin/log-aggregator
EOF

chmod +x /usr/local/bin/start-logger

# Create service file to run log-aggregator
cat <<EOF >/etc/systemd/system/log-aggregator.service
[Unit]
Description=log-aggregator
After=network-online.target
Requires=network-online.target
[Service]
Environment="FAIR_LOG_CURSOR_PATH=/var/log/log-aggregator.cursor"
Environment="FAIR_LOG_FIREHOSE_STREAM=$1"
Environment="FAIR_LOG_FIREHOSE_CREDENTIALS_ENDPOINT=$endpoint"
Environment="ELEVEN_PRODUCT=$eleven_product"
Environment="ELEVEN_COMPONENT=$eleven_component"
Environment="ENV=production"
Environment="SYSTEMD_UNITS_IGNORE=$ignore_systemd_units"
ExecStart=/usr/local/bin/start-logger
Restart=always
RestartSec=5
[Install]
WantedBy=multi-user.target
EOF

# Enable service
systemctl enable log-aggregator.service

# Start service
systemctl start log-aggregator.service
