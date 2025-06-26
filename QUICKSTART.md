# Quick Start Guide - Golang AI Agent

## Langkah 1: Setup Environment

### Install Dependencies
```bash
# Install Go 1.21+
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Git
sudo apt update && sudo apt install -y git

# Install Docker (optional)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

## Langkah 2: Clone dan Setup Project

```bash
# Clone repository
git clone https://github.com/kevinpranata97/golang-ai-agent.git
cd golang-ai-agent

# Setup development environment
make setup

# Install dependencies
make deps
```

## Langkah 3: Konfigurasi

### Environment Variables
```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN"
export WEBHOOK_SECRET="your_webhook_secret"
export PORT="8080"
```

### Edit Configuration File
```bash
cp config.json.example config.json
# Edit config.json dengan settings yang sesuai
```

## Langkah 4: Build dan Run

```bash
# Build application
make build

# Run application
make run

# Atau run dalam development mode
make dev
```

## Langkah 5: Setup GitHub Webhook

1. Buka repository GitHub: https://github.com/kevinpranata97/golang-ai-agent
2. Pergi ke Settings > Webhooks
3. Klik "Add webhook"
4. Payload URL: `http://your-server:8080/webhook`
5. Content type: `application/json`
6. Secret: masukkan WEBHOOK_SECRET yang sama
7. Events: pilih "Push" dan "Pull requests"
8. Klik "Add webhook"

## Langkah 6: Test Installation

```bash
# Health check
curl http://localhost:8080/health

# Status check
curl http://localhost:8080/status

# Test dengan push ke repository
git add .
git commit -m "Test webhook"
git push origin main
```

## Docker Deployment

```bash
# Build Docker image
make docker-build

# Run dengan Docker
make docker-run

# Atau manual
docker run -d \
  -p 8080:8080 \
  -e GITHUB_TOKEN="YOUR_GITHUB_TOKEN" \
  -e WEBHOOK_SECRET="your_secret" \
  -v $(pwd)/data:/root/data \
  golang-ai-agent:latest
```

## Troubleshooting

### Common Issues

1. **Port sudah digunakan**
   ```bash
   export PORT="8081"
   make run
   ```

2. **Permission denied**
   ```bash
   chmod +x golang-ai-agent
   ```

3. **Git authentication error**
   - Pastikan GITHUB_TOKEN valid dan memiliki akses ke repository

### Debug Mode
```bash
export LOG_LEVEL=debug
make dev
```

## Next Steps

1. Customize workflow di `internal/workflow/engine.go`
2. Tambah custom testing di `internal/testing/testing.go`
3. Setup monitoring dan alerting
4. Deploy ke production environment

## Support

- Documentation: [README.md](README.md)
- Issues: https://github.com/kevinpranata97/golang-ai-agent/issues
- Architecture: [architecture_design.md](architecture_design.md)

