# TUGAS 3 — Pengujian dan Dokumentasi Program

Aplikasi Pengolahan Nilai Siswa | LSP Skema Programmer

---

## 1. Dokumentasi Tahapan Pengujian

Pengujian dilakukan dalam tiga lapisan:

1. **Unit test** — menguji fungsi inti pengolahan nilai secara terisolasi
   (`go test ./...`).
2. **Pengujian API (integrasi)** — menjalankan server lalu memanggil endpoint
   memakai `curl` (login, CRUD, validasi nilai, perhitungan, laporan).
3. **Pengujian integrasi frontend–backend** — menjalankan Vite + backend, lalu
   memverifikasi proxy `/api` dan proses login dari sisi aplikasi web.

## 2. Skenario dan Test Case

| ID    | Skenario                          | Input                               | Harapan                                  |
| ----- | --------------------------------- | ----------------------------------- | ---------------------------------------- |
| TC-01 | Login berhasil                    | admin / admin123                    | HTTP 200 + token JWT                      |
| TC-02 | Login gagal (password salah)      | admin / salah                       | HTTP 401, pesan "username atau password salah" |
| TC-03 | Validasi nilai valid              | tugas=80, uts=75, uas=90            | Diterima, nilai akhir = 82.5             |
| TC-04 | Validasi nilai di luar rentang    | tugas=150                           | HTTP 400, pesan validasi                 |
| TC-05 | Hitung nilai akhir                | 90, 85, 95                          | 90.5 (LULUS)                             |
| TC-06 | Penentuan status (batas)          | nilai akhir = 70                    | LULUS                                    |
| TC-07 | Penentuan status (di bawah batas) | nilai akhir = 69.99                 | TIDAK LULUS                              |
| TC-08 | Siswa lihat nilai sendiri         | login siswa → GET /api/nilai        | Hanya nilai milik NIS tsb                |
| TC-09 | Otorisasi peran                   | siswa → GET /api/siswa              | HTTP 403 Forbidden                       |
| TC-10 | Laporan rekap                     | GET /api/laporan                    | total, lulus, rata-rata, dll. benar      |
| TC-11 | Koneksi database                  | jalankan server                     | Log "koneksi database berhasil"          |
| TC-12 | Integrasi proxy frontend          | GET http://localhost:5173/api/health| Diteruskan ke backend, HTTP 200          |

## 3. Hasil Pengujian Aplikasi

### 3.1 Unit Test (`go test ./...`)

```
ok  github.com/othersidedrl/lsp-programmer-course-project/backend/internal/services  2.526s
```

Mencakup `TestValidasiNilai`, `TestHitungNilaiAkhir`, `TestTentukanStatusKelulusan`
(seluruhnya **PASS**).

### 3.2 Pengujian API (curl)

| ID    | Hasil Aktual                                                          | Status |
| ----- | --------------------------------------------------------------------- | :----: |
| TC-01 | `{"success":true,"message":"login berhasil","data":{"token":"eyJ..."}}` | ✅ |
| TC-02 | `HTTP 401`                                                            | ✅ |
| TC-03 | `nilai_akhir: 82.5, status: LULUS`                                    | ✅ |
| TC-04 | `{"success":false,"message":"nilai tugas harus berada pada rentang 0–100 (diterima: 150.00)"}` `[HTTP 400]` | ✅ |
| TC-05 | `nilai_akhir: 90.5, status: LULUS`                                    | ✅ |
| TC-08 | Hanya 1 data (NIS 2024001) muncul untuk akun siswa                    | ✅ |
| TC-09 | `HTTP 403`                                                            | ✅ |
| TC-10 | `total_nilai:2, jumlah_lulus:2, rata_rata_akhir:86.5, tertinggi:90.5, terendah:82.5` | ✅ |
| TC-11 | `koneksi database berhasil -> nilai_siswa.db`                         | ✅ |
| TC-12 | `{"success":true,"message":"server berjalan"}` (via port 5173)        | ✅ |

## 4. Bukti Pengujian Berhasil / Gagal

**Bukti berhasil (TC-05, perhitungan):**

```
POST /api/nilai  body: {"nis":"2024002","tugas":90,"uts":85,"uas":95}
-> {"success":true,"message":"nilai berhasil disimpan",
    "data":{ ... "nilai_akhir":90.5,"status":"LULUS" }}
```

**Bukti gagal yang ditangani dengan benar (TC-04, validasi):**

```
POST /api/nilai  body: {"tugas":150, ...}
-> {"success":false,"message":"nilai tugas harus berada pada rentang 0–100 (diterima: 150.00)"}
   [HTTP 400]
```

**Bukti otorisasi (TC-09):**

```
GET /api/siswa  (token siswa)  ->  HTTP 403
```

## 5. Dokumentasi Debugging

Lihat juga Tugas 2 §6–§7. Ringkasan siklus debugging:

1. **Float presisi (E2)** — gejala: nilai akhir bernilai panjang. Diagnosis:
   aritmetika `float64`. Solusi: fungsi `bulatkan()` (pembulatan 2 desimal).
   Verifikasi ulang lewat `TestHitungNilaiAkhir`.
2. **Foreign key SQLite (E3)** — gejala: hapus siswa tidak menghapus nilainya.
   Diagnosis: PRAGMA FK default OFF. Solusi: DSN `?_pragma=foreign_keys(1)`.
3. **CORS (E4)** — gejala: request frontend ditolak browser. Solusi:
   `corsMiddleware` + proxy Vite. Verifikasi: TC-12 lolos.
4. **Otorisasi (E5)** — gejala: siswa bisa lihat data siswa. Solusi:
   `RequireRole`. Verifikasi: TC-09 → 403.

## 6. Dokumentasi Kode Program

```
backend/
  main.go ............... entry point, daftar route, CORS
  internal/
    config/config.go .... muat konfigurasi dari env
    database/database.go  Connect / Migrate / Seed (SQLite)
    models/ ............. OOP: User, Siswa, Guru, Nilai (+ method)
    services/ ........... TERSTRUKTUR: ValidasiNilai, HitungNilaiAkhir,
                          TentukanStatusKelulusan (+ unit test)
    laporan/ ............ TERSTRUKTUR: BuatLaporan
    middleware/auth.go .. Authenticate (JWT) + RequireRole
    handlers/ ........... auth, siswa, guru, nilai, laporan
    utils/ .............. Response JSON, GenerateToken/ParseToken
frontend/
  src/api/client.js ..... axios + interceptor token
  src/context/ .......... AuthContext (login/logout)
  src/components/ ....... Layout, ProtectedRoute
  src/pages/ ............ Login, Dashboard, Siswa, Guru, Nilai, Laporan, NilaiSaya
```

## 7. Penjelasan Fungsi, Modul, dan Class

### Modul/Fungsi (Pemrograman Terstruktur)

- **`ValidasiNilai`** — memastikan nilai 0–100, mengembalikan `error` bila tidak.
- **`HitungNilaiAkhir`** — menerapkan rumus bobot 30/30/40.
- **`TentukanStatusKelulusan`** — memetakan nilai akhir ke LULUS/TIDAK LULUS.
- **`BuatLaporan`** — meringkas daftar nilai (total, lulus, rata-rata, dsb.).

### Class & Method (OOP)

- **`Siswa`** — atribut NIS, Nama, Kelas; method `Label()`, `Valid()`.
- **`Guru`** — atribut IDGuru, Nama, MataPelajaran; method `Mengampu()`, `Valid()`.
- **`Nilai`** — atribut nilai & status; method `HitungAkhir()`, `Validasi()`,
  `Lulus()` yang memanggil fungsi terstruktur (integrasi OOP + terstruktur).

## 8. Evaluasi Hasil Pengujian Aplikasi

- Seluruh test case (TC-01 s.d. TC-12) **berhasil** sesuai harapan.
- Fungsi inti (validasi, perhitungan, kelulusan) terbukti benar lewat unit test
  otomatis sehingga mudah diuji ulang (regresi).
- Keamanan dasar terpenuhi: hashing password, JWT, otorisasi peran, query
  berparameter.
- Koneksi database & operasi CRUD berjalan baik.

**Kesimpulan:** aplikasi memenuhi seluruh kebutuhan sistem (menyimpan data siswa,
mengelola nilai, menghitung nilai akhir otomatis, menentukan kelulusan, dan
menampilkan laporan) serta berhasil menunjukkan integrasi **pemrograman
terstruktur** dan **OOP** dalam satu aplikasi.

**Saran pengembangan lanjutan:** ekspor laporan ke PDF/Excel, paginasi data,
serta penambahan satuan pendidikan/semester.
