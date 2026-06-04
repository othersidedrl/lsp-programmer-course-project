package models

// Konstanta peran (role) pengguna sesuai hak akses pada sistem.
const (
	RoleAdmin = "admin"
	RoleGuru  = "guru"
	RoleSiswa = "siswa"
)

// User merepresentasikan akun login beserta perannya.
// Atribut Password tidak pernah dikirim ke klien (tag json:"-").
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`      // hash bcrypt
	Role     string `json:"role"`   // admin | guru | siswa
	RefID    string `json:"ref_id"` // NIS untuk siswa, ID guru untuk guru
}
