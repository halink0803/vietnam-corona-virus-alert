package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = `CREATE TABLE IF NOT EXISTS news (
		content TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
)

type Storage struct {
	db *sql.DB
}

func NewSqliteStorage() (*Storage, error) {
	db, err := sql.Open("sqlite3", "./news.db")
	if err != nil {
		return nil, err
	}
	if _, err = db.Exec(schema); err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

// StoreNews store latest news and returns error if any
func (s *Storage) StoreNews(news string) error {
	var (
		query = `INSERT INTO news (content) values ($1);`
	)
	_, err := s.db.Exec(query, news)
	return err
}

//GetLatestNews return latest news store in db
func (s *Storage) GetLatestNews() (string, error) {
	var (
		query  = `SELECT content FROM news ORDER BY timestamp DESC LIMIT 1;`
		result string
	)
	rows := s.db.QueryRow(query)
	if err := rows.Scan(&result); err != nil {
		if err != sql.ErrNoRows {
			return result, err
		}
	}
	return result, nil
}
