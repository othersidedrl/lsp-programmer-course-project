// Package database menangani koneksi ke SQLite, migrasi skema tabel,
// dan pengisian data awal (seeding).
package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite" // driver SQLite murni-Go (tanpa CGO)
)

// Connect membuka koneksi ke database SQLite dan memverifikasinya.
func Connect(path string) (*sql.DB, error) {
	// _pragma foreign_keys(1) -> aktifkan relasi antar tabel
	// busy_timeout(5000)      -> tunggu 5 detik bila database terkunci
	dsn := path + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// SQLite menulis lebih aman bila koneksi diserialkan.
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate membuat tabel-tabel bila belum ada.
func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role     TEXT NOT NULL CHECK (role IN ('admin','guru','siswa')),
		ref_id   TEXT
	);

	CREATE TABLE IF NOT EXISTS siswa (
		nis   TEXT PRIMARY KEY,
		nama  TEXT NOT NULL,
		kelas TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS guru (
		id_guru        TEXT PRIMARY KEY,
		nama           TEXT NOT NULL,
		mata_pelajaran TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS nilai (
		id             INTEGER PRIMARY KEY AUTOINCREMENT,
		nis            TEXT NOT NULL,
		id_guru        TEXT,
		mata_pelajaran TEXT NOT NULL,
		tugas          REAL NOT NULL DEFAULT 0,
		uts            REAL NOT NULL DEFAULT 0,
		uas            REAL NOT NULL DEFAULT 0,
		nilai_akhir    REAL NOT NULL DEFAULT 0,
		status         TEXT NOT NULL DEFAULT 'TIDAK LULUS',
		tervalidasi    INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (nis)     REFERENCES siswa(nis)    ON DELETE CASCADE,
		FOREIGN KEY (id_guru) REFERENCES guru(id_guru) ON DELETE SET NULL
	);`
	_, err := db.Exec(schema)
	return err
}

// Seed mengisi data awal hanya bila tabel users masih kosong.
func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // data sudah ada, lewati
	}
	log.Println("database kosong, mengisi data awal (seeding)...")

	if _, err := db.Exec(`INSERT INTO guru (id_guru, nama, mata_pelajaran) VALUES
		('G001','Budi Santoso','Matematika'),
		('G002','Siti Aminah','Bahasa Indonesia')`); err != nil {
		return err
	}

	if _, err := db.Exec(`INSERT INTO siswa (nis, nama, kelas) VALUES
		('2024001','Ahmad Fauzi','XII RPL 1'),
		('2024002','Dewi Lestari','XII RPL 1'),
		('2024003','Rizki Pratama','XII RPL 2')`); err != nil {
		return err
	}

	users := []struct {
		username, password, role, refID string
	}{
		{"admin", "admin123", "admin", ""},
		{"guru", "guru123", "guru", "G001"},
		{"siswa", "siswa123", "siswa", "2024001"},
	}
	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if _, err := db.Exec(
			"INSERT INTO users (username, password, role, ref_id) VALUES (?,?,?,?)",
			u.username, string(hash), u.role, nullableSeed(u.refID),
		); err != nil {
			return err
		}
	}

	// Contoh satu data nilai: (0.3*80)+(0.3*75)+(0.4*90) = 82.5 -> LULUS.
	if _, err := db.Exec(`INSERT INTO nilai
		(nis, id_guru, mata_pelajaran, tugas, uts, uas, nilai_akhir, status, tervalidasi)
		VALUES ('2024001','G001','Matematika',80,75,90,82.5,'LULUS',1)`); err != nil {
		return err
	}

	log.Println("seeding selesai (akun: admin/admin123, guru/guru123, siswa/siswa123)")
	return nil
}

func nullableSeed(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
