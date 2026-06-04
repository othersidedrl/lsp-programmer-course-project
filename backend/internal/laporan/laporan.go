// Package laporan berisi prosedur pengolahan laporan hasil belajar
// sebagai bagian dari PEMROGRAMAN TERSTRUKTUR.
package laporan

import "github.com/othersidedrl/lsp-programmer-course-project/backend/internal/models"

// Laporan adalah ringkasan hasil belajar dari sekumpulan nilai.
type Laporan struct {
	TotalNilai       int            `json:"total_nilai"`
	JumlahLulus      int            `json:"jumlah_lulus"`
	JumlahTidakLulus int            `json:"jumlah_tidak_lulus"`
	RataRataAkhir    float64        `json:"rata_rata_akhir"`
	NilaiTertinggi   float64        `json:"nilai_tertinggi"`
	NilaiTerendah    float64        `json:"nilai_terendah"`
	Detail           []models.Nilai `json:"detail"`
}

// BuatLaporan mengolah daftar nilai menjadi ringkasan laporan.
// (PROSEDUR PEMROGRAMAN TERSTRUKTUR #4 — pengolahan laporan)
func BuatLaporan(daftar []models.Nilai) Laporan {
	lap := Laporan{Detail: daftar}
	if len(daftar) == 0 {
		return lap
	}

	var total float64
	lap.NilaiTertinggi = daftar[0].NilaiAkhir
	lap.NilaiTerendah = daftar[0].NilaiAkhir

	for _, n := range daftar {
		total += n.NilaiAkhir
		if n.Lulus() {
			lap.JumlahLulus++
		} else {
			lap.JumlahTidakLulus++
		}
		if n.NilaiAkhir > lap.NilaiTertinggi {
			lap.NilaiTertinggi = n.NilaiAkhir
		}
		if n.NilaiAkhir < lap.NilaiTerendah {
			lap.NilaiTerendah = n.NilaiAkhir
		}
	}

	lap.TotalNilai = len(daftar)
	lap.RataRataAkhir = float64(int((total/float64(len(daftar)))*100+0.5)) / 100
	return lap
}
