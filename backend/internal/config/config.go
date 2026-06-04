// Package config memuat konfigurasi aplikasi dari environment variable
// dengan nilai default yang aman sehingga aplikasi dapat langsung dijalankan.
package config

import "os"

// Config menampung seluruh pengaturan yang dibutuhkan aplikasi.
type Config struct {
	Port       string
	DBPath     string
	JWTSecret  string
	CORSOrigin string
}

// Load membaca konfigurasi dari environment variable.
func Load() Config {
	return Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "nilai_siswa.db"),
		JWTSecret:  getEnv("JWT_SECRET", "rahasia-lsp-programmer-2026"),
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:5173"),
	}
}

// getEnv mengembalikan nilai environment variable atau nilai default bila kosong.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
