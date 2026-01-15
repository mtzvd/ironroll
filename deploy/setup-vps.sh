#!/bin/bash
# VPS Setup Script for IronRoll
# Run this once on your VPS to prepare for deployments

set -e

echo "=== Creating ironroll user ==="
sudo useradd -r -s /bin/false ironroll || echo "User already exists"

echo "=== Creating directories ==="
sudo mkdir -p /opt/ironroll
sudo chown ironroll:ironroll /opt/ironroll

echo "=== Installing systemd service ==="
sudo cp ironroll.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable ironroll

echo "=== Creating .env template ==="
if [ ! -f /opt/ironroll/.env ]; then
    sudo tee /opt/ironroll/.env > /dev/null << 'EOF'
TELEGRAM_BOT_TOKEN=your_telegram_token_here
DISCORD_BOT_TOKEN=your_discord_token_here
PORT=8080
EOF
    sudo chown ironroll:ironroll /opt/ironroll/.env
    sudo chmod 600 /opt/ironroll/.env
    echo ">>> Edit /opt/ironroll/.env with your actual tokens!"
fi

echo "=== Setup complete ==="
echo "Next steps:"
echo "1. Edit /opt/ironroll/.env with your bot tokens"
echo "2. Configure GitHub secrets (SSH_PRIVATE_KEY, VPS_HOST, VPS_USER)"
echo "3. Push to master branch to trigger deployment"
