package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/models"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/utils"
)

// SiswaHandler menangani operasi CRUD data siswa.
type SiswaHandler struct{ DB *sql.DB }

// List mengembalikan seluruh data siswa.
func (h *SiswaHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT nis, nama, kelas FROM siswa ORDER BY nis")
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil data siswa")
		return
	}
	defer rows.Close()

	list := make([]models.Siswa, 0)
	for rows.Next() {
		var s models.Siswa
		if err := rows.Scan(&s.NIS, &s.Nama, &s.Kelas); err != nil {
			utils.Error(w, http.StatusInternalServerError, "gagal membaca data siswa")
			return
		}
		list = append(list, s)
	}
	utils.OK(w, "daftar siswa", list)
}

// Get mengembalikan satu siswa berdasarkan NIS.
func (h *SiswaHandler) Get(w http.ResponseWriter, r *http.Request) {
	nis := r.PathValue("nis")
	var s models.Siswa
	err := h.DB.QueryRow("SELECT nis, nama, kelas FROM siswa WHERE nis = ?", nis).
		Scan(&s.NIS, &s.Nama, &s.Kelas)
	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusNotFound, "siswa tidak ditemukan")
		return
	} else if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil data siswa")
		return
	}
	utils.OK(w, "detail siswa", s)
}

// Create menambahkan data siswa baru.
func (h *SiswaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s models.Siswa
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		utils.Error(w, http.StatusBadRequest, "format data tidak valid")
		return
	}
	if !s.Valid() {
		utils.Error(w, http.StatusBadRequest, "NIS, nama, dan kelas wajib diisi")
		return
	}
	_, err := h.DB.Exec("INSERT INTO siswa (nis, nama, kelas) VALUES (?,?,?)", s.NIS, s.Nama, s.Kelas)
	if err != nil {
		utils.Error(w, http.StatusConflict, "gagal menyimpan, NIS mungkin sudah terdaftar")
		return
	}
	utils.JSON(w, http.StatusCreated, utils.Response{Success: true, Message: "siswa berhasil ditambahkan", Data: s})
}

// Update memperbarui nama & kelas siswa berdasarkan NIS.
func (h *SiswaHandler) Update(w http.ResponseWriter, r *http.Request) {
	nis := r.PathValue("nis")
	var s models.Siswa
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		utils.Error(w, http.StatusBadRequest, "format data tidak valid")
		return
	}
	res, err := h.DB.Exec("UPDATE siswa SET nama = ?, kelas = ? WHERE nis = ?", s.Nama, s.Kelas, nis)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal memperbarui data siswa")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		utils.Error(w, http.StatusNotFound, "siswa tidak ditemukan")
		return
	}
	s.NIS = nis
	utils.OK(w, "data siswa berhasil diperbarui", s)
}

// Delete menghapus siswa berdasarkan NIS (beserta nilainya via ON DELETE CASCADE).
func (h *SiswaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	nis := r.PathValue("nis")
	res, err := h.DB.Exec("DELETE FROM siswa WHERE nis = ?", nis)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal menghapus data siswa")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		utils.Error(w, http.StatusNotFound, "siswa tidak ditemukan")
		return
	}
	utils.OK(w, "siswa berhasil dihapus", nil)
}
