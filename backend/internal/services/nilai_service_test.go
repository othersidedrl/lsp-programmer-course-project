package services

import "testing"

// TestValidasiNilai menguji validasi rentang nilai 0–100.
func TestValidasiNilai(t *testing.T) {
	tests := []struct {
		nama    string
		nilai   float64
		wantErr bool
	}{
		{"batas bawah valid (0)", 0, false},
		{"batas atas valid (100)", 100, false},
		{"nilai tengah valid", 75.5, false},
		{"kurang dari 0", -1, true},
		{"lebih dari 100", 101, true},
	}
	for _, tc := range tests {
		t.Run(tc.nama, func(t *testing.T) {
			err := ValidasiNilai("uji", tc.nilai)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidasiNilai(%.2f) error = %v, wantErr = %v", tc.nilai, err, tc.wantErr)
			}
		})
	}
}

// TestHitungNilaiAkhir menguji rumus (30%*Tugas)+(30%*UTS)+(40%*UAS).
func TestHitungNilaiAkhir(t *testing.T) {
	tests := []struct {
		nama            string
		tugas, uts, uas float64
		want            float64
	}{
		{"semua 80", 80, 80, 80, 80},
		{"contoh data seed", 80, 75, 90, 82.5},
		{"nilai rendah", 50, 50, 50, 50},
		{"nilai sempurna", 100, 100, 100, 100},
	}
	for _, tc := range tests {
		t.Run(tc.nama, func(t *testing.T) {
			got := HitungNilaiAkhir(tc.tugas, tc.uts, tc.uas)
			if got != tc.want {
				t.Errorf("HitungNilaiAkhir(%.0f,%.0f,%.0f) = %.2f, want %.2f",
					tc.tugas, tc.uts, tc.uas, got, tc.want)
			}
		})
	}
}

// TestTentukanStatusKelulusan menguji ambang kelulusan (>= 70 LULUS).
func TestTentukanStatusKelulusan(t *testing.T) {
	tests := []struct {
		nilai float64
		want  string
	}{
		{70, "LULUS"},
		{69.99, "TIDAK LULUS"},
		{100, "LULUS"},
		{0, "TIDAK LULUS"},
		{82.5, "LULUS"},
	}
	for _, tc := range tests {
		got := TentukanStatusKelulusan(tc.nilai)
		if got != tc.want {
			t.Errorf("TentukanStatusKelulusan(%.2f) = %q, want %q", tc.nilai, got, tc.want)
		}
	}
}
