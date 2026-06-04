# Aplikasi Pengolahan Nilai Siswa

Aplikasi pengolahan nilai siswa berbasis web untuk institusi pendidikan, dibuat
sebagai project pembekalan **Uji Kompetensi LSP Skema Programmer**.

Aplikasi memadukan dua pendekatan pemrograman dalam satu sistem:

- **Pemrograman Terstruktur** — fungsi/prosedur untuk validasi nilai, perhitungan
  nilai akhir, penentuan status kelulusan, dan pengolahan laporan
  (`backend/internal/services`, `backend/internal/laporan`).
- **Pemrograman Berorientasi Objek (OOP)** — pemodelan `Siswa`, `Guru`, `Nilai`
  sebagai class beserta atribut & method (`backend/internal/models`).

## Teknologi

| Bagian    | Teknologi                                             |
| --------- | ----------------------------------------------------- |
| Frontend  | React 18 + Vite, React Router, Axios                  |
| Backend   | Go (Golang) — `net/http`                              |
| Database  | SQLite (driver murni-Go `modernc.org/sqlite`)         |
| Keamanan  | JWT (`golang-jwt`) + hashing password `bcrypt`        |

## Hak Akses Pengguna

| Peran     | Kemampuan                                                                 |
| --------- | ------------------------------------------------------------------------- |
| **Admin** | CRUD penuh data siswa, guru, nilai + mengelola laporan                    |
| **Guru**  | Input nilai, melihat rekap nilai, memvalidasi nilai (sesuai mapel)        |
| **Siswa** | Melihat nilai pribadi & status kelulusan                                  |

## Rumus & Ketentuan

- **Nilai Akhir** = `(30% × Tugas) + (30% × UTS) + (40% × UAS)`
- Rentang nilai valid: **0 – 100**
- Status: **LULUS** bila Nilai Akhir **≥ 70**, selain itu **TIDAK LULUS**

## Struktur Proyek

```
lsp/
├── backend/                 # REST API (Go)
│   ├── main.go              # entry point + routing + CORS
│   └── internal/
│       ├── config/          # konfigurasi (env)
│       ├── database/        # koneksi SQLite, migrasi, seeding
│       ├── models/          # OOP: Siswa, Guru, Nilai, User (class + method)
│       ├── services/        # TERSTRUKTUR: validasi, hitung, status kelulusan
│       ├── laporan/         # TERSTRUKTUR: pengolahan laporan
│       ├── middleware/      # JWT auth + otorisasi peran
│       ├── handlers/        # handler HTTP tiap endpoint
│       └── utils/           # respons JSON & helper JWT
├── frontend/                # SPA (React + Vite)
│   └── src/
│       ├── api/             # axios client
│       ├── context/         # AuthContext (status login)
│       ├── components/      # Layout, ProtectedRoute
│       └── pages/           # Login, Dashboard, Siswa, Guru, Nilai, Laporan, NilaiSaya
└── docs/                    # Dokumen Tugas 1, 2, 3
```

## Cara Menjalankan

### 1. Backend (Go)

> Butuh Go 1.21+ terpasang.

```bash
cd backend
go mod tidy          # unduh dependensi (sekali saja)
go run .             # server berjalan di http://localhost:8080
```

Saat pertama dijalankan, database `nilai_siswa.db` dibuat otomatis (migrasi +
seeding data awal). Di terminal akan muncul:

```
koneksi database berhasil -> nilai_siswa.db
seeding selesai (akun: admin/admin123, guru/guru123, siswa/siswa123)
server berjalan di http://localhost:8080
```

### 2. Frontend (React + Vite)

> Butuh Node.js 18+.

```bash
cd frontend
npm install
npm run dev          # aplikasi di http://localhost:5173
```

Buka **http://localhost:5173** di browser. Request `/api/*` otomatis di-proxy ke
backend (lihat `frontend/vite.config.js`).

### Akun Demo

| Peran | Username | Password   |
| ----- | -------- | ---------- |
| Admin | `admin`  | `admin123` |
| Guru  | `guru`   | `guru123`  |
| Siswa | `siswa`  | `siswa123` |

## Daftar Endpoint API

| Method | Endpoint                     | Hak Akses          | Keterangan                          |
| ------ | ---------------------------- | ------------------ | ----------------------------------- |
| POST   | `/api/login`                 | publik             | Login, mengembalikan token JWT      |
| GET    | `/api/me`                    | semua (login)      | Profil pengguna                     |
| GET    | `/api/siswa`                 | admin, guru        | Daftar siswa                        |
| POST   | `/api/siswa`                 | admin              | Tambah siswa                        |
| PUT    | `/api/siswa/{nis}`           | admin              | Ubah siswa                          |
| DELETE | `/api/siswa/{nis}`           | admin              | Hapus siswa                         |
| GET    | `/api/guru`                  | admin, guru        | Daftar guru                         |
| POST   | `/api/guru`                  | admin              | Tambah guru                         |
| PUT    | `/api/guru/{id}`             | admin              | Ubah guru                           |
| DELETE | `/api/guru/{id}`             | admin              | Hapus guru                          |
| GET    | `/api/nilai`                 | semua (login)      | Daftar nilai (difilter per peran)   |
| POST   | `/api/nilai`                 | admin, guru        | Input nilai (+ hitung otomatis)     |
| PUT    | `/api/nilai/{id}`            | admin, guru        | Ubah nilai                          |
| PATCH  | `/api/nilai/{id}/validasi`   | admin, guru        | Validasi nilai                      |
| DELETE | `/api/nilai/{id}`            | admin              | Hapus nilai                         |
| GET    | `/api/laporan`               | admin, guru        | Laporan rekap hasil belajar         |

## Pengujian

```bash
cd backend
go test ./...        # unit test fungsi nilai (validasi, hitung, kelulusan)
```

## Dokumentasi Lengkap

- [Tugas 1 — Analisis & Perancangan Sistem](docs/TUGAS-1-Analisis-dan-Perancangan.md)
- [Tugas 2 — Implementasi Program](docs/TUGAS-2-Implementasi.md)
- [Tugas 3 — Pengujian & Dokumentasi Program](docs/TUGAS-3-Pengujian-dan-Dokumentasi.md)
