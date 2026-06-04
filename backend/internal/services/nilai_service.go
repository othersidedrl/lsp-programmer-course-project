// Package services berisi fungsi/prosedur logika program sebagai
// implementasi PEMROGRAMAN TERSTRUKTUR: validasi nilai, perhitungan
// nilai akhir, dan penentuan status kelulusan. Fungsi-fungsi di sini
// bekerja pada tipe data primitif sehingga mudah diuji dan dipakai ulang.
package services

import "fmt"

// Bobot & ketentuan penilaian sesuai spesifikasi sistem.
const (
	BobotTugas = 0.30  // 30% nilai tugas
	BobotUTS   = 0.30  // 30% nilai UTS
	BobotUAS   = 0.40  // 40% nilai UAS

	NilaiMin   = 0.0   // batas bawah nilai valid
	NilaiMaks  = 100.0 // batas atas nilai valid
	BatasLulus = 70.0  // nilai akhir minimal untuk lulus
)

// ValidasiNilai memastikan sebuah nilai berada pada rentang 0–100.
// (FUNGSI PEMROGRAMAN TERSTRUKTUR #1 — validasi nilai)
func ValidasiNilai(label string, nilai float64) error {
	if nilai < NilaiMin || nilai > NilaiMaks {
		return fmt.Errorf("nilai %s harus berada pada rentang %.0f–%.0f (diterima: %.2f)",
			label, NilaiMin, NilaiMaks, nilai)
	}
	return nil
}

// ValidasiSemuaNilai memvalidasi nilai tugas, UTS, dan UAS sekaligus.
func ValidasiSemuaNilai(tugas, uts, uas float64) error {
	if err := ValidasiNilai("tugas", tugas); err != nil {
		return err
	}
	if err := ValidasiNilai("UTS", uts); err != nil {
		return err
	}
	if err := ValidasiNilai("UAS", uas); err != nil {
		return err
	}
	return nil
}

// HitungNilaiAkhir menghitung nilai akhir dengan rumus:
// (30% × Tugas) + (30% × UTS) + (40% × UAS), dibulatkan 2 desimal.
// (FUNGSI PEMROGRAMAN TERSTRUKTUR #2 — perhitungan nilai akhir)
func HitungNilaiAkhir(tugas, uts, uas float64) float64 {
	akhir := (BobotTugas * tugas) + (BobotUTS * uts) + (BobotUAS * uas)
	return bulatkan(akhir)
}

// TentukanStatusKelulusan mengembalikan "LULUS" bila nilai akhir >= 70,
// selain itu "TIDAK LULUS".
// (FUNGSI PEMROGRAMAN TERSTRUKTUR #3 — penentuan status kelulusan)
func TentukanStatusKelulusan(nilaiAkhir float64) string {
	if nilaiAkhir >= BatasLulus {
		return "LULUS"
	}
	return "TIDAK LULUS"
}

// bulatkan membulatkan angka ke 2 angka di belakang koma.
func bulatkan(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
