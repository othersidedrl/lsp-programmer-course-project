// Package handlers berisi handler HTTP untuk setiap endpoint API.
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/middleware"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/models"
	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/utils"
)

// AuthHandler menangani proses login dan info profil pengguna.
type AuthHandler struct {
	DB        *sql.DB
	JWTSecret string
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login memvalidasi kredensial dan mengembalikan token JWT beserta peran.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "format permintaan tidak valid")
		return
	}
	if req.Username == "" || req.Password == "" {
		utils.Error(w, http.StatusBadRequest, "username dan password wajib diisi")
		return
	}

	var u models.User
	err := h.DB.QueryRow(
		"SELECT id, username, password, role, COALESCE(ref_id,'') FROM users WHERE username = ?",
		req.Username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.RefID)
	if err == sql.ErrNoRows {
		utils.Error(w, http.StatusUnauthorized, "username atau password salah")
		return
	} else if err != nil {
		utils.Error(w, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		utils.Error(w, http.StatusUnauthorized, "username atau password salah")
		return
	}

	token, err := utils.GenerateToken(h.JWTSecret, u.ID, u.Username, u.Role, u.RefID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "gagal membuat token")
		return
	}

	utils.OK(w, "login berhasil", map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"username": u.Username,
			"role":     u.Role,
			"ref_id":   u.RefID,
		},
	})
}

// Me mengembalikan profil pengguna dari token aktif.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "tidak terautentikasi")
		return
	}
	utils.OK(w, "profil pengguna", map[string]string{
		"username": claims.Username,
		"role":     claims.Role,
		"ref_id":   claims.RefID,
	})
}
