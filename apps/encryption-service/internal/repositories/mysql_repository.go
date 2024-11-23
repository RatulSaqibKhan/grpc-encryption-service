package repository

import (
	"database/sql"
)

type MySQLRepository struct {
	DB *sql.DB
}

func (r *MySQLRepository) GetEncrypted(plaintext string) (string, error) {
	var encrypted string
	err := r.DB.QueryRow("SELECT encrypted FROM demodb.encryption WHERE plaintext = ?", plaintext).Scan(&encrypted)
	return encrypted, err
}

func (r *MySQLRepository) GetDecrypted(encrypted string) (string, error) {
	var plaintext string
	err := r.DB.QueryRow("SELECT plaintext FROM demodb.encryption WHERE encrypted = ?", encrypted).Scan(&plaintext)
	return plaintext, err
}

func (r *MySQLRepository) SaveEncrypted(plaintext, encrypted string) error {
	_, err := r.DB.Exec("INSERT INTO demodb.encryption (plaintext, encrypted) VALUES (?, ?)", plaintext, encrypted)
	return err
}
