package models

// Siswa adalah class (OOP) yang merepresentasikan data seorang siswa.
// Menerapkan konsep OOP berupa atribut (field) dan method (perilaku).
type Siswa struct {
	NIS   string `json:"nis"` // Nomor Induk Siswa (unik)
	Nama  string `json:"nama"`
	Kelas string `json:"kelas"`
}

// Label mengembalikan representasi ringkas siswa. (METHOD OOP)
func (s Siswa) Label() string {
	return s.NIS + " - " + s.Nama
}

// Valid memeriksa kelengkapan atribut wajib siswa. (METHOD OOP)
func (s Siswa) Valid() bool {
	return s.NIS != "" && s.Nama != "" && s.Kelas != ""
}
