## Desain Arsitektur Modul Fine-tuning dengan Database Lokal

### Pendahuluan

Modul fine-tuning akan menjadi komponen krusial yang memungkinkan agen AI untuk belajar dan beradaptasi dari setiap interaksi. Desain ini akan merinci struktur database lokal (SQLite), skema data untuk menyimpan interaksi, dan bagaimana modul fine-tuning akan memproses data ini untuk meningkatkan performa agen.

### 1. Skema Database SQLite

Kita akan menggunakan SQLite sebagai database lokal karena sifatnya yang ringan dan kemudahan integrasinya dengan Go. Database akan berisi setidaknya satu tabel utama untuk menyimpan log interaksi yang akan digunakan untuk fine-tuning.

**Tabel: `interactions_log`**

| Kolom                 | Tipe Data       | Deskripsi                                                                 |
| :-------------------- | :-------------- | :------------------------------------------------------------------------ |
| `id`                  | `TEXT PRIMARY KEY` | ID unik untuk setiap interaksi (UUID)                                     |
| `timestamp`           | `TEXT`          | Waktu interaksi (ISO 8601 format)                                         |
| `endpoint`            | `TEXT`          | Endpoint API yang dipanggil (e.g., `/generate-app`, `/test-app`)          |
| `request_payload`     | `TEXT`          | Payload JSON dari request pengguna                                        |
| `response_payload`    | `TEXT`          | Payload JSON dari response agen                                           |
| `app_name`            | `TEXT`          | Nama aplikasi yang dihasilkan/diuji (jika ada)                            |
| `app_path`            | `TEXT`          | Path ke direktori aplikasi yang dihasilkan (jika ada)                     |
| `test_results_json`   | `TEXT`          | Hasil pengujian dalam format JSON (jika ada)                              |
| `analysis_results_json` | `TEXT`          | Hasil analisis kode dalam format JSON (jika ada)                          |
| `feedback_json`       | `TEXT`          | Feedback pengguna eksplisit (jika ada, dalam JSON)                        |
| `status`              | `TEXT`          | Status interaksi (e.g., `success`, `failure`, `error`)                    |
| `processed_for_finetuning` | `INTEGER`    | Flag (0/1) menunjukkan apakah data ini sudah diproses untuk fine-tuning |

**Indeks:**
*   `CREATE INDEX idx_timestamp ON interactions_log (timestamp);`
*   `CREATE INDEX idx_endpoint ON interactions_log (endpoint);`
*   `CREATE INDEX idx_processed ON interactions_log (processed_for_finetuning);`

### 2. Modul Database (`internal/database`) Baru

Untuk mengelola interaksi dengan SQLite, kita akan membuat modul baru `internal/database`. Modul ini akan bertanggung jawab untuk:

*   **Inisialisasi Database**: Membuka koneksi ke file SQLite, membuat tabel jika belum ada.
*   **Penyimpanan Log Interaksi**: Fungsi untuk menyimpan data interaksi baru ke tabel `interactions_log`.
*   **Pengambilan Data untuk Fine-tuning**: Fungsi untuk mengambil data interaksi yang belum diproses untuk fine-tuning.
*   **Pembaruan Status**: Fungsi untuk memperbarui status `processed_for_finetuning` setelah data digunakan.

```go
// internal/database/database.go
package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"log"
	"os"
	"path/filepath"
	"time"
)

const dbFileName = "finetuning.db"

type InteractionLog struct {
	ID                     string
	Timestamp              time.Time
	Endpoint               string
	RequestPayload         string
	ResponsePayload        string
	AppName                string
	AppPath                string
	TestResultsJSON        string
	AnalysisResultsJSON    string
	FeedbackJSON           string
	Status                 string
	ProcessedForFinetuning bool
}

type DB struct {
	*sql.DB
}

func NewDB(dataDir string) (*DB, error) {
	dbPath := filepath.Join(dataDir, dbFileName)
	// Ensure the directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Printf("Database initialized at %s", dbPath)
	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS interactions_log (
		id TEXT PRIMARY KEY,
		timestamp TEXT NOT NULL,
		endpoint TEXT NOT NULL,
		request_payload TEXT,
		response_payload TEXT,
		app_name TEXT,
		app_path TEXT,
		test_results_json TEXT,
		analysis_results_json TEXT,
		feedback_json TEXT,
		status TEXT NOT NULL,
		processed_for_finetuning INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON interactions_log (timestamp);
	CREATE INDEX IF NOT EXISTS idx_endpoint ON interactions_log (endpoint);
	CREATE INDEX IF NOT EXISTS idx_processed ON interactions_log (processed_for_finetuning);
	`
	_, err := db.Exec(sqlStmt)
	return err
}

func (d *DB) InsertInteractionLog(logEntry InteractionLog) error {
	stmt, err := d.Prepare(`
	INSERT INTO interactions_log (
		id, timestamp, endpoint, request_payload, response_payload, app_name, app_path,
		test_results_json, analysis_results_json, feedback_json, status, processed_for_finetuning
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		logEntry.ID,
		logEntry.Timestamp.Format(time.RFC3339),
		logEntry.Endpoint,
		logEntry.RequestPayload,
		logEntry.ResponsePayload,
		logEntry.AppName,
		logEntry.AppPath,
		logEntry.TestResultsJSON,
		logEntry.AnalysisResultsJSON,
		logEntry.FeedbackJSON,
		logEntry.Status,
		logEntry.ProcessedForFinetuning,
	)
	return err
}

func (d *DB) GetUnprocessedLogs() ([]InteractionLog, error) {
	rows, err := d.Query(`
	SELECT id, timestamp, endpoint, request_payload, response_payload, app_name, app_path,
		test_results_json, analysis_results_json, feedback_json, status, processed_for_finetuning
	FROM interactions_log
	WHERE processed_for_finetuning = 0
	ORDER BY timestamp ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query unprocessed logs: %w", err)
	}
	defer rows.Close()

	var logs []InteractionLog
	for rows.Next() {
		var logEntry InteractionLog
		var timestampStr string
		var processedInt int
		if err := rows.Scan(
			&logEntry.ID, &timestampStr, &logEntry.Endpoint, &logEntry.RequestPayload,
			&logEntry.ResponsePayload, &logEntry.AppName, &logEntry.AppPath,
			&logEntry.TestResultsJSON, &logEntry.AnalysisResultsJSON, &logEntry.FeedbackJSON,
			&logEntry.Status, &processedInt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		logEntry.Timestamp, err = time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		logEntry.ProcessedForFinetuning = (processedInt == 1)
		logs = append(logs, logEntry)
	}

	return logs, nil
}

func (d *DB) MarkLogsAsProcessed(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	// Create a string of '?' placeholders for the IN clause
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(`
	UPDATE interactions_log
	SET processed_for_finetuning = 1
	WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	stmt, err := d.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err = stmt.Exec(args...)
	return err
}

```

### 3. Pembaruan Modul Fine-tuning (`internal/finetuning`)

Modul `finetuning` akan diperbarui untuk berinteraksi dengan database lokal. Ini akan mencakup:

*   **`Finetuner` struct**: Akan memegang referensi ke instance `database.DB`.
*   **`ProcessLogs()` method**: Akan mengambil log yang belum diproses dari database, menganalisisnya, dan menerapkan logika fine-tuning. Logika fine-tuning di sini akan berfokus pada peningkatan prompt engineering dan penyesuaian rule-based berdasarkan hasil interaksi.
*   **`Train()` method (placeholder)**: Untuk masa depan, ini bisa menjadi tempat di mana model AI yang lebih kecil dilatih ulang secara lokal jika diperlukan.

```go
// internal/finetuning/tuner.go (updated)
package finetuning

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/database"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
	"github.com/kevinpranata97/golang-ai-agent/internal/codegen"
	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
)

type Finetuner struct {
	db *database.DB
	// Tambahkan referensi ke komponen lain yang mungkin perlu di-fine-tune
	// Misalnya, requirements.Analyzer, codegen.Generator, dll.
}

func NewFinetuner(db *database.DB) *Finetuner {
	return &Finetuner{db: db}
}

// ProcessLogs mengambil log interaksi yang belum diproses dan menerapkan logika fine-tuning.
func (f *Finetuner) ProcessLogs() error {
	logs, err := f.db.GetUnprocessedLogs()
	if err != nil {
		return fmt.Errorf("failed to get unprocessed logs: %w", err)
	}

	if len(logs) == 0 {
		log.Println("No new interaction logs to process for fine-tuning.")
		return nil
	}

	log.Printf("Processing %d interaction logs for fine-tuning...", len(logs))
	var processedIDs []string

	for _, entry := range logs {
		// Contoh logika fine-tuning sederhana:
		// Jika generate-app berhasil dan test results menunjukkan kegagalan, analisis mengapa.
		if entry.Endpoint == "/generate-app" && entry.Status == "success" && entry.TestResultsJSON != "" {
			var testResults apptesting.TestResults // Asumsikan struktur TestResults
			if err := json.Unmarshal([]byte(entry.TestResultsJSON), &testResults); err == nil {
				if testResults.OverallStatus == "failure" {
					log.Printf("Fine-tuning opportunity: App '%s' generated successfully but failed tests. Analyzing...", entry.AppName)
					// Di sini, kita akan menambahkan logika untuk menganalisis request_payload,
					// response_payload, dan test_results_json untuk mengidentifikasi pola kegagalan.
					// Misalnya, jika sering gagal karena masalah database, mungkin perlu menyesuaikan prompt
					// atau template kode untuk database.
					// Ini adalah bagian yang paling kompleks dan akan memerlukan AI/heuristik yang lebih canggih.
				}
			}
		}
		// Tambahkan logika fine-tuning lainnya berdasarkan endpoint dan status

		processedIDs = append(processedIDs, entry.ID)
	}

	if len(processedIDs) > 0 {
		if err := f.db.MarkLogsAsProcessed(processedIDs); err != nil {
			return fmt.Errorf("failed to mark logs as processed: %w", err)
		}
		log.Printf("Successfully processed %d logs for fine-tuning.", len(processedIDs))
	}

	return nil
}

// Train method is a placeholder for future, more advanced model training.
func (f *Finetuner) Train() error {
	log.Println("Starting advanced fine-tuning model training (placeholder).")
	// Implementasi pelatihan model AI yang sebenarnya akan ada di sini.
	// Ini bisa melibatkan loading model, melatihnya dengan data dari database,
	// dan menyimpan model yang telah di-fine-tune.
	return nil
}
```

### 4. Integrasi ke `main.go`

*   **Inisialisasi Database**: Instance `database.DB` akan diinisialisasi di `main.go`.
*   **Inisialisasi Finetuner**: Instance `finetuning.Finetuner` akan dibuat dengan meneruskan instance `database.DB`.
*   **Logging Interaksi**: Setiap handler endpoint (`/generate-app`, `/test-app`, `/generate-and-test`) akan dimodifikasi untuk mencatat detail interaksi ke database menggunakan `db.InsertInteractionLog`.
*   **Scheduler Fine-tuning**: Sebuah goroutine atau scheduler sederhana akan ditambahkan di `main.go` untuk memanggil `finetuner.ProcessLogs()` secara berkala (misalnya, setiap 5 menit).

```go
// main.go (modifikasi)
package main

import (
	// ... import yang sudah ada ...
	"github.com/kevinpranata97/golang-ai-agent/internal/database"
	"github.com/kevinpranata97/golang-ai-agent/internal/finetuning"
	"github.com/google/uuid"
	"time"
)

func main() {
	// ... inisialisasi yang sudah ada ...

	// Inisialisasi Database Lokal untuk Fine-tuning
	dataDir := "./data"
	db, err := database.NewDB(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Inisialisasi Finetuner
	finetuner := finetuning.NewFinetuner(db)

	// Jadwalkan proses fine-tuning secara berkala
	go func() {
		for {
			log.Println("Running scheduled fine-tuning process...")
			if err := finetuner.ProcessLogs(); err != nil {
				log.Printf("Error during scheduled fine-tuning: %v", err)
			}
			time.Sleep(5 * time.Minute) // Proses setiap 5 menit
		}
	}()

	// ... HTTP Handlers yang sudah ada ...

	// Modifikasi handler untuk mencatat interaksi
	http.HandleFunc("/generate-app", func(w http.ResponseWriter, r *http.Request) {
		// ... logika handler yang sudah ada ...

		logEntry := database.InteractionLog{
			ID:            uuid.New().String(),
			Timestamp:     time.Now(),
			Endpoint:      "/generate-app",
			RequestPayload:  string(reqBody),
			ResponsePayload: string(jsonResponse),
			AppName:       appReq.Name,
			AppPath:       filepath.Join(outputDir, strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-"))),
			Status:        "success", // Atau "failure" jika ada error
		}
		if err := db.InsertInteractionLog(logEntry); err != nil {
			log.Printf("Failed to log interaction: %v", err)
		}
	})

	// Lakukan hal yang sama untuk /test-app dan /generate-and-test

	// ... Start server ...
}
```

### 5. Pembaruan Modul `storage` (Opsional)

Modul `storage` yang ada saat ini menggunakan file JSON. Jika data fine-tuning akan disimpan di SQLite, maka modul `storage` yang ada mungkin perlu direfaktor atau diganti dengan `internal/database` untuk konsistensi, atau `internal/database` akan menjadi bagian dari `internal/storage` jika kita ingin mengelola semua persistensi data di satu tempat. Untuk saat ini, kita akan membuat `internal/database` sebagai modul terpisah untuk fine-tuning.

### 6. Alur Kerja Fine-tuning

1.  **Interaksi Pengguna**: Pengguna memanggil endpoint seperti `/generate-app`.
2.  **Log Interaksi**: Detail request, response, dan metadata relevan dicatat ke tabel `interactions_log` di SQLite.
3.  **Proses Fine-tuning Terjadwal**: Setiap 5 menit (atau interval yang dikonfigurasi), `finetuner.ProcessLogs()` dipanggil.
4.  **Analisis Log**: `ProcessLogs` mengambil log yang belum diproses (`processed_for_finetuning = 0`).
5.  **Pembelajaran & Penyesuaian**: Logika fine-tuning menganalisis data (misalnya, jika aplikasi yang dihasilkan gagal tes, atau jika feedback pengguna menunjukkan masalah). Berdasarkan analisis ini, agen dapat:
    *   Menyesuaikan prompt yang dikirim ke Google Gemini API.
    *   Memodifikasi aturan dalam modul `requirements` atau `codegen`.
    *   Mencatat pola untuk pelatihan model AI yang lebih canggih di masa depan.
6.  **Penandaan Log**: Log yang telah diproses ditandai sebagai `processed_for_finetuning = 1`.

### Tantangan Implementasi

*   **Logika Fine-tuning yang Cerdas**: Bagian tersulit adalah mendesain algoritma yang dapat secara otomatis belajar dari data interaksi dan menerjemahkannya ke dalam perbaikan yang berarti pada agen. Ini mungkin memerlukan heuristik, analisis statistik, atau bahkan model AI kecil yang dilatih secara lokal.
*   **Manajemen Konkurensi**: Memastikan penulisan ke database log tidak mengganggu performa endpoint API.
*   **Skalabilitas Data**: Meskipun SQLite baik untuk lokal, perlu dipertimbangkan bagaimana data akan dikelola jika volume interaksi sangat besar.

Desain ini menyediakan kerangka kerja yang solid untuk menambahkan kemampuan fine-tuning adaptif ke agen AI. Langkah selanjutnya adalah mengimplementasikan modul `internal/database` dan mengintegrasikannya ke `main.go` dan handler endpoint yang relevan.

