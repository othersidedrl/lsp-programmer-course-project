# TUGAS 2 — Implementasi Program

Aplikasi Pengolahan Nilai Siswa | LSP Skema Programmer

---

## 1. Halaman Login Berdasarkan Role

File: `frontend/src/pages/Login.jsx` + `frontend/src/context/AuthContext.jsx`.

Pengguna login dengan username & password. Backend memverifikasi kredensial,
membuat token JWT yang memuat **peran** (role), lalu frontend mengarahkan ke
halaman sesuai peran. Menu navigasi (`frontend/src/components/Layout.jsx`)
ditampilkan berbeda untuk admin, guru, dan siswa.

Pembatasan akses sisi server ada di `backend/internal/middleware/auth.go`
(`Authenticate` + `RequireRole`). Contoh respons login berhasil:

```json
{
  "success": true,
  "message": "login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { "username": "admin", "role": "admin", "ref_id": "" }
  }
}
```

## 2. Form Input Data Siswa dan Nilai

- **Data siswa**: `frontend/src/pages/SiswaPage.jsx` — form NIS, Nama, Kelas
  (CRUD penuh untuk admin).
- **Data nilai**: `frontend/src/pages/NilaiPage.jsx` — form pilih Siswa, Mata
  Pelajaran, Guru, lalu Tugas/UTS/UAS. Untuk guru, mata pelajaran terisi
  otomatis sesuai mapel yang diampu.

Validasi sisi klien menggunakan `<input type="number" min="0" max="100">`,
sedangkan validasi otoritatif tetap di sisi server.

## 3. Proses Perhitungan Nilai Akhir

Perhitungan dilakukan oleh fungsi terstruktur `HitungNilaiAkhir` lalu dipanggil
melalui method OOP `Nilai.HitungAkhir()`.

```go
// backend/internal/services/nilai_service.go
const (
    BobotTugas = 0.30
    BobotUTS   = 0.30
    BobotUAS   = 0.40
    BatasLulus = 70.0
)

func HitungNilaiAkhir(tugas, uts, uas float64) float64 {
    akhir := (BobotTugas * tugas) + (BobotUTS * uts) + (BobotUAS * uas)
    return bulatkan(akhir) // dibulatkan 2 desimal
}
```

Contoh: Tugas=80, UTS=75, UAS=90 →
`(0.3×80)+(0.3×75)+(0.4×90) = 24 + 22.5 + 36 = 82.5` → **LULUS**.

## 4. Laporan Hasil Nilai Siswa

File: `frontend/src/pages/LaporanPage.jsx` (tampilan) +
`backend/internal/laporan/laporan.go` (pengolahan) +
`backend/internal/handlers/laporan_handler.go` (endpoint `GET /api/laporan`).

Laporan menampilkan ringkasan: total nilai, jumlah lulus / tidak lulus,
rata-rata, nilai tertinggi & terendah, beserta tabel detail. Tersedia tombol
**Cetak** (memanfaatkan `window.print()` + CSS `@media print`).

Contoh respons:

```json
{
  "total_nilai": 2, "jumlah_lulus": 2, "jumlah_tidak_lulus": 0,
  "rata_rata_akhir": 86.5, "nilai_tertinggi": 90.5, "nilai_terendah": 82.5,
  "detail": [ /* ... */ ]
}
```

## 5. Bukti Pengujian Database

Pengujian dilakukan dengan menjalankan server lalu memanggil API. Log startup
membuktikan koneksi & seeding:

```
koneksi database berhasil -> nilai_siswa.db
seeding selesai (akun: admin/admin123, guru/guru123, siswa/siswa123)
```

Operasi tulis & baca terbukti dari `POST /api/nilai` (menyimpan) dan
`GET /api/nilai` (membaca dengan JOIN ke tabel siswa). Lihat hasil lengkap pada
[Tugas 3](TUGAS-3-Pengujian-dan-Dokumentasi.md).

## 6. Catatan Error / Debugging

| No | Masalah ditemukan                                              | Penyebab                                                       |
| -- | -------------------------------------------------------------- | ------------------------------------------------------------- |
| E1 | Module path masih placeholder `nilaisiswa`                     | `go mod init` awal belum memakai path repo GitHub             |
| E2 | Hasil nilai akhir kadang muncul desimal panjang (mis. 82.4999)| Operasi `float64` tidak presisi (0.3 tidak eksak biner)       |
| E3 | Relasi antar tabel (FK) tidak aktif di SQLite                  | Secara default `PRAGMA foreign_keys` SQLite = OFF             |
| E4 | Browser menolak request: error CORS                            | Frontend (port 5173) beda origin dengan backend (port 8080)   |
| E5 | Siswa bisa mengakses endpoint data siswa                       | Endpoint belum dibatasi peran                                 |

## 7. Perbaikan Error dan Hasil Setelah Diperbaiki

| No | Perbaikan                                                                                  | Hasil                                            |
| -- | ------------------------------------------------------------------------------------------ | ------------------------------------------------ |
| E1 | Ganti module → `github.com/othersidedrl/lsp-programmer-course-project/backend`, rewrite import | `go build ./...` & `go vet ./...` lolos          |
| E2 | Tambah fungsi `bulatkan()` (pembulatan 2 desimal) di `services`                            | Nilai akhir rapi, mis. `82.5`, `90.5`            |
| E3 | Aktifkan FK lewat DSN `?_pragma=foreign_keys(1)`                                           | Hapus siswa otomatis menghapus nilainya (CASCADE)|
| E4 | Tambah `corsMiddleware` di backend + `proxy` di `vite.config.js`                           | Request frontend→backend berhasil (HTTP 200)     |
| E5 | Bungkus endpoint dengan `RequireRole(...)`                                                 | Siswa akses `/api/siswa` → **403 Forbidden**     |

## 8. Potongan Kode Fungsi/Procedure (Pemrograman Terstruktur)

```go
// backend/internal/services/nilai_service.go

// (1) Validasi nilai 0–100
func ValidasiNilai(label string, nilai float64) error {
    if nilai < NilaiMin || nilai > NilaiMaks {
        return fmt.Errorf("nilai %s harus berada pada rentang %.0f–%.0f (diterima: %.2f)",
            label, NilaiMin, NilaiMaks, nilai)
    }
    return nil
}

// (2) Hitung nilai akhir berbobot
func HitungNilaiAkhir(tugas, uts, uas float64) float64 {
    akhir := (BobotTugas * tugas) + (BobotUTS * uts) + (BobotUAS * uas)
    return bulatkan(akhir)
}

// (3) Tentukan status kelulusan
func TentukanStatusKelulusan(nilaiAkhir float64) string {
    if nilaiAkhir >= BatasLulus {
        return "LULUS"
    }
    return "TIDAK LULUS"
}
```

```go
// backend/internal/laporan/laporan.go
// (4) Prosedur pengolahan laporan
func BuatLaporan(daftar []models.Nilai) Laporan {
    lap := Laporan{Detail: daftar}
    if len(daftar) == 0 { return lap }
    var total float64
    lap.NilaiTertinggi, lap.NilaiTerendah = daftar[0].NilaiAkhir, daftar[0].NilaiAkhir
    for _, n := range daftar {
        total += n.NilaiAkhir
        if n.Lulus() { lap.JumlahLulus++ } else { lap.JumlahTidakLulus++ }
        if n.NilaiAkhir > lap.NilaiTertinggi { lap.NilaiTertinggi = n.NilaiAkhir }
        if n.NilaiAkhir < lap.NilaiTerendah { lap.NilaiTerendah = n.NilaiAkhir }
    }
    lap.TotalNilai = len(daftar)
    lap.RataRataAkhir = bulatkanRata(total, len(daftar))
    return lap
}
```

## 9. Potongan Kode Class dan Method (OOP)

```go
// backend/internal/models/siswa.go
type Siswa struct {
    NIS   string `json:"nis"`
    Nama  string `json:"nama"`
    Kelas string `json:"kelas"`
}
func (s Siswa) Label() string { return s.NIS + " - " + s.Nama }
func (s Siswa) Valid() bool   { return s.NIS != "" && s.Nama != "" && s.Kelas != "" }
```

```go
// backend/internal/models/nilai.go
type Nilai struct {
    NIS        string  `json:"nis"`
    Tugas      float64 `json:"tugas"`
    UTS        float64 `json:"uts"`
    UAS        float64 `json:"uas"`
    NilaiAkhir float64 `json:"nilai_akhir"`
    Status     string  `json:"status"`
    // ...
}

// Method OOP yang memanggil fungsi terstruktur (integrasi 2 pendekatan)
func (n *Nilai) HitungAkhir() {
    n.NilaiAkhir = services.HitungNilaiAkhir(n.Tugas, n.UTS, n.UAS)
    n.Status     = services.TentukanStatusKelulusan(n.NilaiAkhir)
}
func (n Nilai) Validasi() error { return services.ValidasiSemuaNilai(n.Tugas, n.UTS, n.UAS) }
func (n Nilai) Lulus() bool     { return n.NilaiAkhir >= services.BatasLulus }
```

## 10. Penjelasan Library / Komponen yang Digunakan

### Backend (Go)

| Library                       | Fungsi                                                         |
| ----------------------------- | -------------------------------------------------------------- |
| `net/http` (standar)          | Server HTTP & routing (`ServeMux` dengan pola method+path)     |
| `modernc.org/sqlite`          | Driver SQLite murni-Go (tanpa CGO) untuk `database/sql`        |
| `github.com/golang-jwt/jwt/v5`| Membuat & memverifikasi token JWT                              |
| `golang.org/x/crypto/bcrypt`  | Hashing & verifikasi password                                  |
| `database/sql` (standar)      | Abstraksi akses database (query, exec, prepared statement)     |

### Frontend (React)

| Library            | Fungsi                                            |
| ------------------ | ------------------------------------------------- |
| `react`, `react-dom` | Library UI berbasis komponen                    |
| `vite`             | Build tool & dev server cepat                     |
| `react-router-dom` | Routing antar halaman + proteksi rute             |
| `axios`            | HTTP client (interceptor untuk menyisipkan token) |

## 11. Penjelasan Coding Guidelines & Best Practices

1. **Pemisahan lapisan (separation of concerns)** — `models` (data/OOP),
   `services`/`laporan` (logika terstruktur), `handlers` (HTTP), `middleware`
   (auth), `database` (persistensi). Memudahkan pemeliharaan & pengujian.
2. **Penamaan jelas & konsisten** — nama fungsi/variabel deskriptif
   (`HitungNilaiAkhir`, `TentukanStatusKelulusan`).
3. **Validasi terpusat** — aturan nilai 0–100 & bobot ditaruh di satu tempat
   (konstanta `services`) sehingga tidak ada "magic number" tersebar.
4. **Keamanan** — password di-hash (bcrypt), endpoint dilindungi JWT + peran,
   query memakai parameter (`?`) untuk mencegah SQL injection.
5. **Penanganan error eksplisit** — setiap error dicek & dibalas pesan yang
   ramah pengguna dalam Bahasa Indonesia.
6. **Format respons konsisten** — semua API memakai struktur
   `{ success, message, data }`.
7. **Komponen frontend reusable** — `ProtectedRoute`, `Layout`, state form
   memakai hooks (`useState`, `useEffect`).
8. **Verifikasi otomatis** — `go vet`, `go build`, dan `go test` dijalankan
   untuk memastikan kode bersih sebelum dipakai.
