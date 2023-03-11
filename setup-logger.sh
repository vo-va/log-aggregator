#!/bin/bash

if [[ -z $FAIR_LOG_FIREHOSE_STREAM || -z $ELEVEN_PRODUCT || -z $ELEVEN_COMPONENT ]]; then
	echo "Usage: FAIR_LOG_FIREHOSE_STREAM, ELEVEN_PRODUCT and ELEVEN_COMPONENT env variables need to be set "
	echo "WARNING logger will not be installed, safely exiting with 0 to not abort deployment"
	exit 0
fi

# Install logger to /usr/bin
curl -SL https://github.com/eleven-software/log-aggregator/releases/download/1.2/log-aggregator_1.2 -o /usr/local/bin/log-aggregator

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
Environment="FAIR_LOG_FIREHOSE_STREAM=$FAIR_LOG_FIREHOSE_STREAM"
Environment="FAIR_LOG_FIREHOSE_CREDENTIALS_ENDPOINT=$endpoint"
Environment="ELEVEN_PRODUCT=$ELEVEN_PRODUCT"
Environment="ELEVEN_COMPONENT=$ELEVEN_COMPONENT"
Environment="ENV=production"
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
