package handlers

import (
	"database/sql"
	"net/http"

	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/laporan"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/utils"
)

// LaporanHandler menangani penyajian laporan hasil belajar.
type LaporanHandler struct{ DB *sql.DB }

// Laporan menyusun ringkasan hasil belajar dari seluruh nilai.
func (h *LaporanHandler) Laporan(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(nilaiSelect + " ORDER BY n.nilai_akhir DESC")
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal mengambil data untuk laporan")
		return
	}
	defer rows.Close()

	list, err := scanNilai(rows)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal membaca data laporan")
		return
	}

	// Pengolahan laporan dilakukan oleh prosedur terstruktur BuatLaporan.
	hasil := laporan.BuatLaporan(list)
	utils.OK(w, "laporan hasil belajar siswa", hasil)
}
