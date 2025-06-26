## Golang AI Agent

Aplikasi AI Agent yang dikembangkan menggunakan Go (Golang) ini adalah solusi komprehensif untuk **pembuatan dan pengetesan aplikasi secara otomatis**. Agen ini dirancang untuk mengubah deskripsi bahasa alami menjadi aplikasi yang berfungsi penuh, lengkap dengan pengujian, analisis kualitas, dan integrasi CI/CD.

### Fitur Utama:

-   **Generasi Aplikasi Berbasis AI**: Mengubah deskripsi bahasa alami menjadi kode aplikasi Go yang berfungsi penuh (API, Web, CLI).
-   **Pengujian Komprehensif**: Melakukan unit test, integration test, static analysis, security scan, dan performance benchmark secara otomatis.
-   **Analisis Cerdas**: Memberikan wawasan mendalam tentang kualitas kode, keamanan, dan performa aplikasi yang dihasilkan.
-   **Fine-tuning Iteratif**: Secara otomatis mengidentifikasi dan menerapkan perbaikan untuk meningkatkan kualitas dan performa aplikasi.
-   **Penyimpanan & Debugging**: Menyimpan data proyek, riwayat analisis, dan menyediakan alat debugging.
-   **Workflow Otomatis**: Integrasi penuh dengan GitHub Actions untuk CI/CD otomatis.

### Struktur Proyek:

-   `main.go`: Titik masuk utama aplikasi.
-   `internal/agent/`: Logika inti agen AI.
-   `internal/requirements/`: Modul untuk analisis kebutuhan dari deskripsi bahasa alami.
-   `internal/codegen/`: Modul untuk menghasilkan kode aplikasi.
-   `internal/apptesting/`: Modul untuk melakukan pengujian komprehensif pada aplikasi yang dihasilkan.
-   `internal/analysis/`: Modul untuk analisis kualitas kode dan performa.
-   `internal/finetuning/`: Modul untuk fine-tuning dan perbaikan otomatis.
-   `internal/storage/`: Modul untuk persistensi data proyek.
-   `internal/debugging/`: Modul untuk debugging dan logging.
-   `Dockerfile`: Konfigurasi untuk membangun image Docker aplikasi.
-   `.github/workflows/ci-cd.yml`: Definisi pipeline CI/CD GitHub Actions.
-   `deploy.sh`: Script untuk mempermudah deployment lokal, Docker, atau Docker Compose.
-   `docker-compose.yml`: Konfigurasi Docker Compose untuk deployment multi-layanan.
-   `nginx.conf`: Contoh konfigurasi Nginx untuk reverse proxy.
-   `FINAL_DOCUMENTATION.pdf`: Dokumentasi teknis lengkap proyek.

### Cara Menggunakan (Quick Start):

1.  **Clone Repositori:**
    ```bash
    git clone https://github.com/kevinpranata97/golang-ai-agent.git
    cd golang-ai-agent
    ```

2.  **Deployment Lokal:**
    ```bash
    ./deploy.sh local
    ```
    Aplikasi akan berjalan di `http://localhost:8080`.

3.  **Deployment dengan Docker Compose (Disarankan untuk Produksi):**
    Pastikan Docker dan Docker Compose terinstal.
    ```bash
    ./deploy.sh compose
    ```
    Aplikasi akan berjalan di `http://localhost:80` (melalui Nginx proxy).

4.  **Menggunakan API (Contoh):**

    **Generate Aplikasi:**
    ```bash
    curl -X POST http://localhost:8080/generate-app \
      -H "Content-Type: application/json" \
      -d '{"description": "Create a simple blog API with posts and comments"}'
    ```

    **Generate dan Test Aplikasi Sekaligus:**
    ```bash
    curl -X POST http://localhost:8080/generate-and-test \
      -H "Content-Type: application/json" \
      -d '{"description": "Create a simple task management API"}'
    ```

### GitHub Actions CI/CD:

Setiap push ke branch `main` atau `develop` akan memicu pipeline CI/CD otomatis yang mencakup:

-   **Test dan Quality Checks**: Unit test, static analysis, code formatting.
-   **Build Aplikasi**: Kompilasi aplikasi Go.
-   **Security Scan**: Pemindaian kerentanan dengan Gosec dan Trivy.
-   **Build dan Push Docker Image**: Membuat dan mendorong image Docker ke Docker Hub (hanya untuk `main`).
-   **Integration Tests**: Menguji endpoint API agen.
-   **Performance Tests**: Mengukur performa endpoint.
-   **Deployment**: Notifikasi deployment ke lingkungan produksi (jika dikonfigurasi).

Anda dapat memantau status pipeline di halaman GitHub Actions repositori Anda.

---

**Catatan:** Untuk penggunaan fitur yang memerlukan akses ke Google Gemini API atau GitHub API, pastikan Anda telah mengkonfigurasi `GEMINI_API_KEY` dan `GITHUB_TOKEN` di lingkungan Anda atau di file `.env` jika menggunakan Docker Compose.gkan menggunakan Golang untuk mengotomatisasi pembuatan dan pengetesan aplikasi secara umum, termasuk aplikasi berbasis web, analisis, fine-tuning, penyimpanan sederhana, debugging, dan pembuatan workflow otomatis untuk repositori GitHub.

## Fitur Utama

### 🤖 Core AI Agent
- **Webhook Handler**: Menerima dan memproses webhook dari GitHub
- **Workflow Orchestration**: Mengorkestrasi alur kerja CI/CD otomatis
- **Status Monitoring**: Memantau status dan kesehatan agen

### 🔧 Web Testing & Analysis
- **Unit Testing**: Menjalankan tes unit untuk berbagai bahasa pemrograman (Go, JavaScript, Python)
- **Integration Testing**: Melakukan pengetesan integrasi
- **Code Analysis**: Analisis statis kode untuk menemukan masalah dan kerentanan
- **Performance Testing**: Pengujian kinerja dan load testing
- **Security Scanning**: Pemindaian keamanan untuk menemukan kerentanan

### 🐛 Debugging & Monitoring
- **Code Issue Detection**: Deteksi masalah umum dalam kode
- **Log Analysis**: Analisis log untuk menemukan error dan warning
- **Performance Profiling**: Profiling kinerja aplikasi
- **Memory Leak Detection**: Deteksi kebocoran memori
- **Suggestion Engine**: Memberikan saran perbaikan berdasarkan analisis

### 📊 Storage & Analytics
- **File-based Storage**: Penyimpanan sederhana berbasis file JSON
- **Data Persistence**: Menyimpan hasil analisis dan laporan
- **Storage Statistics**: Statistik penggunaan penyimpanan
- **Data Cleanup**: Pembersihan data lama secara otomatis

### 🔄 Workflow Automation
- **CI/CD Pipeline**: Pipeline otomatis untuk build, test, dan deploy
- **Multi-language Support**: Dukungan untuk Go, JavaScript, Python, dan lainnya
- **Parallel Execution**: Eksekusi tugas secara paralel
- **Retry Mechanism**: Mekanisme retry untuk tugas yang gagal

## Instalasi

### Prasyarat
- Go 1.21 atau lebih baru
- Git
- Docker (opsional)

### Clone Repository
```bash
git clone https://github.com/kevinpranata97/golang-ai-agent.git
cd golang-ai-agent
```

### Install Dependencies
```bash
go mod download
```

### Build Application
```bash
go build -o golang-ai-agent .
```

## Konfigurasi

### Environment Variables
```bash
export GITHUB_TOKEN="your_github_token"
export WEBHOOK_SECRET="your_webhook_secret"
export PORT="8080"
```

### Configuration File (config.json)
```json
{
  "server": {
    "port": "8080",
    "host": "0.0.0.0",
    "read_timeout": 30,
    "write_timeout": 30
  },
  "github": {
    "token": "your_github_token",
    "webhook_secret": "your_webhook_secret",
    "base_url": "https://api.github.com"
  },
  "storage": {
    "type": "file",
    "path": "./data"
  },
  "testing": {
    "timeout": 300,
    "parallel": true,
    "coverage": true,
    "security_scan": true
  },
  "debugging": {
    "log_level": "info",
    "profile_mode": false,
    "max_sessions": 5
  },
  "workflow": {
    "max_concurrent": 3,
    "retry_attempts": 3,
    "cleanup_after": 24
  }
}
```

## Penggunaan

### Menjalankan Agen
```bash
./golang-ai-agent
```

### Setup GitHub Webhook
1. Buka repository GitHub Anda
2. Pergi ke Settings > Webhooks
3. Klik "Add webhook"
4. Masukkan URL: `http://your-server:8080/webhook`
5. Pilih "application/json" sebagai Content type
6. Masukkan secret yang sama dengan WEBHOOK_SECRET
7. Pilih events: "Push" dan "Pull requests"

### API Endpoints

#### Health Check
```bash
GET /health
```

#### Status Agent
```bash
GET /status
```

#### Generate Application
```bash
POST /generate-app
```
**Description:** Generates a new application based on a natural language description.
**Request Body (JSON):**
```json
{
  "description": "Create a simple blog API with posts and comments"
}
```

#### Test Application
```bash
POST /test-app
```
**Description:** Tests an existing application at a given path.
**Request Body (JSON):**
```json
{
  "app_path": "/path/to/your/generated_app"
}
```

#### Generate and Test Application
```bash
POST /generate-and-test
```
**Description:** Generates an application and immediately runs tests on it.
**Request Body (JSON):**
```json
{
  "description": "Create a simple task management API"
}
```

#### Webhook Handler
```bash
POST /webhook
```

## Arsitektur

### Komponen Utama

1. **Agent Core** (`internal/agent/`)
   - Mengelola webhook dan orkestrasi workflow
   - Mengintegrasikan semua modul lainnya

2. **GitHub Integration** (`internal/github/`)
   - Klien untuk berinteraksi dengan GitHub API
   - Kloning repository dan analisis struktur

3. **Testing Module** (`internal/testing/`)
   - Menjalankan berbagai jenis tes
   - Analisis kode dan keamanan

4. **Debugging Module** (`internal/debugging/`)
   - Deteksi masalah dan debugging
   - Profiling dan analisis performa

5. **Storage Module** (`internal/storage/`)
   - Penyimpanan data dan konfigurasi
   - Manajemen file dan cleanup

6. **Workflow Engine** (`internal/workflow/`)
   - Eksekusi workflow CI/CD
   - Manajemen tugas dan retry

### Alur Kerja

1. **GitHub Event** → Webhook diterima
2. **Workflow Trigger** → Workflow CI/CD dimulai
3. **Repository Clone** → Repository di-clone ke temporary directory
4. **Analysis & Testing** → Kode dianalisis dan ditest
5. **Debugging** → Masalah dideteksi dan saran diberikan
6. **Results Storage** → Hasil disimpan ke storage
7. **GitHub Feedback** → Status dikirim kembali ke GitHub

## Docker

### Build Image
```bash
docker build -t golang-ai-agent .
```

### Run Container
```bash
docker run -d \
  -p 8080:8080 \
  -e GITHUB_TOKEN="your_token" \
  -e WEBHOOK_SECRET="your_secret" \
  -v $(pwd)/data:/root/data \
  golang-ai-agent
```

## CI/CD

Proyek ini menggunakan GitHub Actions untuk CI/CD pipeline yang mencakup:

- **Testing**: Unit tests, integration tests, dan security scanning
- **Linting**: Code quality checks dengan golangci-lint
- **Building**: Multi-platform builds (Linux, macOS, Windows)
- **Docker**: Automated Docker image building dan publishing
- **Deployment**: Automated deployment ke production

## Monitoring & Logging

### Health Monitoring
- Health check endpoint tersedia di `/health`
- Status monitoring di `/status`
- Metrics collection untuk performance monitoring

### Logging
- Structured logging dengan level konfigurasi
- Log rotation dan cleanup otomatis
- Error tracking dan alerting

## Security

### Best Practices
- Webhook signature verification
- Token-based authentication untuk GitHub API
- Input validation dan sanitization
- Secure secret management
- Container security scanning

### Security Features
- Automated security scanning
- Vulnerability detection
- Dependency checking
- Code analysis untuk security issues

## Troubleshooting

### Common Issues

1. **Webhook tidak diterima**
   - Periksa URL webhook di GitHub settings
   - Pastikan port 8080 dapat diakses dari internet
   - Verifikasi webhook secret

2. **Build gagal**
   - Pastikan Go version 1.21+
   - Jalankan `go mod tidy` untuk update dependencies
   - Periksa log error untuk detail

3. **Tests gagal**
   - Periksa dependencies yang diperlukan
   - Pastikan Git tersedia untuk cloning
   - Verifikasi permissions untuk temporary directories

### Debug Mode
```bash
export LOG_LEVEL=debug
./golang-ai-agent
```

## Contributing

1. Fork repository
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Kevin Pranata - [@kevinpranata97](https://github.com/kevinpranata97)

Project Link: [https://github.com/kevinpranata97/golang-ai-agent](https://github.com/kevinpranata97/golang-ai-agent)

## Acknowledgments

- [Go](https://golang.org/) - Programming language
- [GitHub API](https://docs.github.com/en/rest) - GitHub integration
- [Docker](https://www.docker.com/) - Containerization
- [GitHub Actions](https://github.com/features/actions) - CI/CD pipeline

