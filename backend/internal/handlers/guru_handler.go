package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/models"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/utils"
)

// GuruHandler menangani operasi CRUD data guru.
type GuruHandler struct{ DB *sql.DB }

// List mengembalikan seluruh data guru.
func (h *GuruHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id_guru, nama, mata_pelajaran FROM guru ORDER BY id_guru")
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil data guru")
		return
	}
	defer rows.Close()

	list := make([]models.Guru, 0)
	for rows.Next() {
		var g models.Guru
		if err := rows.Scan(&g.IDGuru, &g.Nama, &g.MataPelajaran); err != nil {
			utils.Error(w, http.StatusInternalServerError, "gagal membaca data guru")
			return
		}
		list = append(list, g)
	}
	utils.OK(w, "daftar guru", list)
}

// Create menambahkan data guru baru.
func (h *GuruHandler) Create(w http.ResponseWriter, r *http.Request) {
	var g models.Guru
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		utils.Error(w, http.StatusBadRequest, "format data tidak valid")
		return
	}
	if !g.Valid() {
		utils.Error(w, http.StatusBadRequest, "ID guru, nama, dan mata pelajaran wajib diisi")
		return
	}
	_, err := h.DB.Exec("INSERT INTO guru (id_guru, nama, mata_pelajaran) VALUES (?,?,?)",
		g.IDGuru, g.Nama, g.MataPelajaran)
	if err != nil {
		utils.Error(w, http.StatusConflict, "gagal menyimpan, ID guru mungkin sudah terdaftar")
		return
	}
	utils.JSON(w, http.StatusCreated, utils.Response{Success: true, Message: "guru berhasil ditambahkan", Data: g})
}

// Update memperbarui nama & mata pelajaran guru.
func (h *GuruHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var g models.Guru
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		utils.Error(w, http.StatusBadRequest, "format data tidak valid")
		return
	}
	res, err := h.DB.Exec("UPDATE guru SET nama = ?, mata_pelajaran = ? WHERE id_guru = ?",
		g.Nama, g.MataPelajaran, id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal memperbarui data guru")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		utils.Error(w, http.StatusNotFound, "guru tidak ditemukan")
		return
	}
	g.IDGuru = id
	utils.OK(w, "data guru berhasil diperbarui", g)
}

// Delete menghapus guru berdasarkan ID.
func (h *GuruHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := h.DB.Exec("DELETE FROM guru WHERE id_guru = ?", id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal menghapus data guru")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		utils.Error(w, http.StatusNotFound, "guru tidak ditemukan")
		return
	}
	utils.OK(w, "guru berhasil dihapus", nil)
}
