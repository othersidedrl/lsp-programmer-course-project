package models

import "github.com/othersidedrl/lsp-programmer-course-project/backend/internal/services"

// Nilai adalah class (OOP) yang merepresentasikan nilai seorang siswa
// pada satu mata pelajaran. Method-methodnya mendelegasikan logika
// perhitungan ke package services (pemrograman terstruktur), sehingga
// menunjukkan integrasi OOP + pemrograman terstruktur.
type Nilai struct {
	ID            int     `json:"id"`
	NIS           string  `json:"nis"`
	IDGuru        string  `json:"id_guru"`
	MataPelajaran string  `json:"mata_pelajaran"`
	Tugas         float64 `json:"tugas"`
	UTS           float64 `json:"uts"`
	UAS           float64 `json:"uas"`
	NilaiAkhir    float64 `json:"nilai_akhir"`
	Status        string  `json:"status"`
	Tervalidasi   bool    `json:"tervalidasi"`

	// Field tampilan hasil JOIN (tidak disimpan langsung di tabel nilai).
	NamaSiswa string `json:"nama_siswa,omitempty"`
	Kelas     string `json:"kelas,omitempty"`
}

// HitungAkhir menghitung lalu menetapkan NilaiAkhir dan Status.
// Method ini memanggil fungsi terstruktur HitungNilaiAkhir dan
// TentukanStatusKelulusan. (METHOD OOP)
func (n *Nilai) HitungAkhir() {
	n.NilaiAkhir = services.HitungNilaiAkhir(n.Tugas, n.UTS, n.UAS)
	n.Status = services.TentukanStatusKelulusan(n.NilaiAkhir)
}

// Validasi memvalidasi rentang nilai melalui fungsi terstruktur. (METHOD OOP)
func (n Nilai) Validasi() error {
	return services.ValidasiSemuaNilai(n.Tugas, n.UTS, n.UAS)
}

// Lulus mengembalikan true bila siswa dinyatakan lulus. (METHOD OOP)
func (n Nilai) Lulus() bool {
	return n.NilaiAkhir >= services.BatasLulus
}
