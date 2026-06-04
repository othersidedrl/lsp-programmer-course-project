// Command nilaisiswa adalah REST API untuk Aplikasi Pengolahan Nilai Siswa.
// Mengintegrasikan pemrograman terstruktur (package services & laporan)
// dengan OOP (package models) dalam satu aplikasi.
package main

import (
	"log"
	"net/http"

	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/config"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/database"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/handlers"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/middleware"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/models"
)

func main() {
	cfg := config.Load()

	// 1) Koneksi database
	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("KONEKSI DATABASE GAGAL: %v", err)
	}
	defer db.Close()
	log.Printf("koneksi database berhasil -> %s", cfg.DBPath)

	// 2) Migrasi & seeding
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrasi database gagal: %v", err)
	}
	if err := database.Seed(db); err != nil {
		log.Fatalf("seeding database gagal: %v", err)
	}

	// 3) Inisialisasi handler
	auth := &handlers.AuthHandler{DB: db, JWTSecret: cfg.JWTSecret}
	siswa := &handlers.SiswaHandler{DB: db}
	guru := &handlers.GuruHandler{DB: db}
	nilai := &handlers.NilaiHandler{DB: db}
	lapor := &handlers.LaporanHandler{DB: db}

	authMW := middleware.Authenticate(cfg.JWTSecret)
	mux := http.NewServeMux()

	// Endpoint publik
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"message":"server berjalan"}`))
	})
	mux.HandleFunc("POST /api/login", auth.Login)

	// protect membungkus handler dengan autentikasi + pembatasan peran.
	protect := func(roles ...string) func(http.HandlerFunc) http.Handler {
		return func(h http.HandlerFunc) http.Handler {
			return authMW(middleware.RequireRole(roles...)(h))
		}
	}
	semua := protect(models.RoleAdmin, models.RoleGuru, models.RoleSiswa)
	adminGuru := protect(models.RoleAdmin, models.RoleGuru)
	adminOnly := protect(models.RoleAdmin)

	mux.Handle("GET /api/me", semua(auth.Me))

	// Data siswa (CRUD: admin penuh; guru hanya melihat daftar)
	mux.Handle("GET /api/siswa", adminGuru(siswa.List))
	mux.Handle("GET /api/siswa/{nis}", adminOnly(siswa.Get))
	mux.Handle("POST /api/siswa", adminOnly(siswa.Create))
	mux.Handle("PUT /api/siswa/{nis}", adminOnly(siswa.Update))
	mux.Handle("DELETE /api/siswa/{nis}", adminOnly(siswa.Delete))

	// Data guru (CRUD: admin penuh; guru hanya melihat daftar)
	mux.Handle("GET /api/guru", adminGuru(guru.List))
	mux.Handle("POST /api/guru", adminOnly(guru.Create))
	mux.Handle("PUT /api/guru/{id}", adminOnly(guru.Update))
	mux.Handle("DELETE /api/guru/{id}", adminOnly(guru.Delete))

	// Data nilai
	mux.Handle("GET /api/nilai", semua(nilai.List))                         // semua peran (difilter sesuai peran)
	mux.Handle("POST /api/nilai", adminGuru(nilai.Create))                  // admin & guru input nilai
	mux.Handle("PUT /api/nilai/{id}", adminGuru(nilai.Update))              // admin & guru ubah nilai
	mux.Handle("PATCH /api/nilai/{id}/validasi", adminGuru(nilai.Validate)) // validasi nilai
	mux.Handle("DELETE /api/nilai/{id}", adminOnly(nilai.Delete))           // hapus nilai (admin)

	// Laporan
	mux.Handle("GET /api/laporan", adminGuru(lapor.Laporan))

	handler := corsMiddleware(cfg.CORSOrigin)(mux)

	log.Printf("server berjalan di http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}

// corsMiddleware mengizinkan akses dari frontend (Vite dev server).
func corsMiddleware(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
