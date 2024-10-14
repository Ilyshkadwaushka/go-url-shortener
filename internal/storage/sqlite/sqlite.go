package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3" // init sqlite3 driver
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err) // Показываем, где именно произошла ошибка
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	  id INTEGER PRIMARY KEY,
	  alias TEXT NOT NULL UNIQUE,
	  url TEXT NOT NULL);
	
	CREATE INDEX IF NOT EXISTS idx_aliad ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// s - receiver которая ссылается на экземпляр типа Storage
// *Storage - указатель на тип Storage. Без него мы бы работали с копией

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// TODO: refactor this
		// Возвращаем общую ошибку для всех storage
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.getURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resUrl string

	err = stmt.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resUrl, nil
}

// fund (s *Storage) DeleteURL(alias string) error
