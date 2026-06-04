package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/middleware"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/models"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/utils"
)

// NilaiHandler menangani pengelolaan nilai siswa.
type NilaiHandler struct{ DB *sql.DB }

// nilaiSelect adalah query dasar dengan JOIN ke tabel siswa (dipakai ulang oleh laporan).
const nilaiSelect = `
SELECT n.id, n.nis, COALESCE(n.id_guru,''), n.mata_pelajaran,
       n.tugas, n.uts, n.uas, n.nilai_akhir, n.status, n.tervalidasi,
       s.nama, s.kelas
FROM nilai n
JOIN siswa s ON s.nis = n.nis`

// scanNilai membaca hasil query menjadi slice models.Nilai.
func scanNilai(rows *sql.Rows) ([]models.Nilai, error) {
	list := make([]models.Nilai, 0)
	for rows.Next() {
		var n models.Nilai
		var tervalidasi int
		if err := rows.Scan(&n.ID, &n.NIS, &n.IDGuru, &n.MataPelajaran,
			&n.Tugas, &n.UTS, &n.UAS, &n.NilaiAkhir, &n.Status, &tervalidasi,
			&n.NamaSiswa, &n.Kelas); err != nil {
			return nil, err
		}
		n.Tervalidasi = tervalidasi == 1
		list = append(list, n)
	}
	return list, nil
}

// List mengembalikan nilai sesuai peran:
//   - admin : seluruh nilai
//   - guru  : nilai yang diampunya (sesuai id_guru)
//   - siswa : hanya nilai miliknya (sesuai NIS)
func (h *NilaiHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var (
		rows *sql.Rows
		err  error
	)
	switch claims.Role {
	case models.RoleSiswa:
		rows, err = h.DB.Query(nilaiSelect+" WHERE n.nis = ? ORDER BY n.mata_pelajaran", claims.RefID)
	case models.RoleGuru:
		rows, err = h.DB.Query(nilaiSelect+" WHERE n.id_guru = ? ORDER BY s.nama", claims.RefID)
	default:
		rows, err = h.DB.Query(nilaiSelect + " ORDER BY s.nama, n.mata_pelajaran")
	}
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil data nilai")
		return
	}
	defer rows.Close()

	list, err := scanNilai(rows)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal membaca data nilai")
		return
	}
	utils.OK(w, "daftar nilai", list)
}

type nilaiRequest struct {
	NIS           string  `json:"nis"`
	IDGuru        string  `json:"id_guru"`
	MataPelajaran string  `json:"mata_pelajaran"`
	Tugas         float64 `json:"tugas"`
	UTS           float64 `json:"uts"`
	UAS           float64 `json:"uas"`
}

// Create menyimpan nilai baru. Validasi & perhitungan dilakukan lewat method OOP.
func (h *NilaiHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var req nilaiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "format data tidak valid")
		return
	}

	n := models.Nilai{
		NIS:           req.NIS,
		IDGuru:        req.IDGuru,
		MataPelajaran: req.MataPelajaran,
		Tugas:         req.Tugas,
		UTS:           req.UTS,
		UAS:           req.UAS,
	}

	// Guru hanya boleh menginput nilai untuk mata pelajaran yang diampunya.
	if claims.Role == models.RoleGuru {
		n.IDGuru = claims.RefID
		var mapel string
		if err := h.DB.QueryRow("SELECT mata_pelajaran FROM guru WHERE id_guru = ?", claims.RefID).
			Scan(&mapel); err == nil {
			n.MataPelajaran = mapel
		}
	}

	if n.NIS == "" || n.MataPelajaran == "" {
		utils.Error(w, http.StatusBadRequest, "NIS dan mata pelajaran wajib diisi")
		return
	}

	// Validasi rentang nilai (method OOP -> fungsi terstruktur).
	if err := n.Validasi(); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// Hitung nilai akhir + status kelulusan (method OOP -> fungsi terstruktur).
	n.HitungAkhir()

	res, err := h.DB.Exec(`INSERT INTO nilai
		(nis, id_guru, mata_pelajaran, tugas, uts, uas, nilai_akhir, status, tervalidasi)
		VALUES (?,?,?,?,?,?,?,?,0)`,
		n.NIS, nullable(n.IDGuru), n.MataPelajaran, n.Tugas, n.UTS, n.UAS, n.NilaiAkhir, n.Status)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "gagal menyimpan nilai, pastikan NIS sudah terdaftar")
		return
	}
	id, _ := res.LastInsertId()
	n.ID = int(id)
	utils.JSON(w, http.StatusCreated, utils.Response{Success: true, Message: "nilai berhasil disimpan", Data: n})
}

// Update memperbarui nilai. Perubahan mereset status validasi menjadi 0.
func (h *NilaiHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id nilai tidak valid")
		return
	}
	var req nilaiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "format data tidak valid")
		return
	}

	n := models.Nilai{Tugas: req.Tugas, UTS: req.UTS, UAS: req.UAS}
	if err := n.Validasi(); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	n.HitungAkhir()

	res, err := h.DB.Exec(
		"UPDATE nilai SET tugas=?, uts=?, uas=?, nilai_akhir=?, status=?, tervalidasi=0 WHERE id=?",
		n.Tugas, n.UTS, n.UAS, n.NilaiAkhir, n.Status, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal memperbarui nilai")
		return
	}
	if cnt, _ := res.RowsAffected(); cnt == 0 {
		utils.Error(w, http.StatusNotFound, "data nilai tidak ditemukan")
		return
	}
	n.ID = id
	utils.OK(w, "nilai diperbarui (status validasi direset, perlu validasi ulang)", n)
}

// Validate menandai nilai sebagai tervalidasi (oleh guru/admin).
func (h *NilaiHandler) Validate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id nilai tidak valid")
		return
	}
	res, err := h.DB.Exec("UPDATE nilai SET tervalidasi = 1 WHERE id = ?", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal memvalidasi nilai")
		return
	}
	if cnt, _ := res.RowsAffected(); cnt == 0 {
		utils.Error(w, http.StatusNotFound, "data nilai tidak ditemukan")
		return
	}
	utils.OK(w, "nilai berhasil divalidasi", nil)
}

// Delete menghapus nilai berdasarkan ID.
func (h *NilaiHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "id nilai tidak valid")
		return
	}
	res, err := h.DB.Exec("DELETE FROM nilai WHERE id = ?", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal menghapus nilai")
		return
	}
	if cnt, _ := res.RowsAffected(); cnt == 0 {
		utils.Error(w, http.StatusNotFound, "data nilai tidak ditemukan")
		return
	}
	utils.OK(w, "nilai berhasil dihapus", nil)
}

// nullable mengubah string kosong menjadi NULL untuk kolom opsional.
func nullable(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
