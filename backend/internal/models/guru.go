package models

// Guru adalah class (OOP) untuk data guru pengampu mata pelajaran.
type Guru struct {
	IDGuru        string `json:"id_guru"`
	Nama          string `json:"nama"`
	MataPelajaran string `json:"mata_pelajaran"`
}

// Mengampu mengembalikan true bila guru mengampu mata pelajaran tertentu. (METHOD OOP)
func (g Guru) Mengampu(mapel string) bool {
	return g.MataPelajaran == mapel
}

// Valid memeriksa kelengkapan atribut wajib guru. (METHOD OOP)
func (g Guru) Valid() bool {
	return g.IDGuru != "" && g.Nama != "" && g.MataPelajaran != ""
}
