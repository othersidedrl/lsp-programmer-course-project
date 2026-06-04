// Package middleware berisi middleware HTTP untuk autentikasi (JWT) dan
// otorisasi berbasis peran (hak akses pengguna: admin, guru, siswa).
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/othersidedrl/lsp-programmer-course-project/backend/internal/utils"
)

type ctxKey string

const claimsKey ctxKey = "claims"

// Authenticate memverifikasi token JWT pada header "Authorization: Bearer <token>".
func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				utils.Error(w, http.StatusUnauthorized, "token tidak ditemukan, silakan login")
				return
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := utils.ParseToken(secret, tokenStr)
			if err != nil {
				utils.Error(w, http.StatusUnauthorized, "sesi tidak valid, silakan login kembali")
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole membatasi akses endpoint hanya untuk peran tertentu.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r)
			if claims == nil {
				utils.Error(w, http.StatusUnauthorized, "tidak terautentikasi")
				return
			}
			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.Error(w, http.StatusForbidden, "anda tidak memiliki hak akses untuk aksi ini")
		})
	}
}

// GetClaims mengambil data klaim JWT dari context request.
func GetClaims(r *http.Request) *utils.Claims {
	claims, ok := r.Context().Value(claimsKey).(*utils.Claims)
	if !ok {
		return nil
	}
	return claims
}
