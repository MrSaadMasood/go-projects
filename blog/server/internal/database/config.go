package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("postgres", "postgres://admin:admin@localhost:5432/blog")
	if err != nil {
		panic(fmt.Errorf("database connection failed %w", err))
	}
	err = DB.Ping()
	if err != nil {
		panic(fmt.Errorf("failed to ping the database %w", err))
	}
}

func Setup(db *sql.DB) error {
	createUserTable := `CREATE TABLE users (
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
	  email VARCHAR(255) NOT NULL UNIQUE,
	  created_at DATE NOT NULL
	) IF NOT EXISTS`

	createBlogsTable := `CREATE TABLE blogs (
	  author_id SERIAL REFERENCES users,
    title VARCHAR(500) NOT NULL,
	  body TEXT NOT NULL,
	  created_at DATE NOT NULL,
	  updated_at DATE NOT NULL
	) IF NOT EXISTS`

	createCommentsTable := `CREATE TABLE comments (
	  author_id SERIAL REFERENCES users,
	  body VARCHAR(500) NOT NULL,
	  created_at DATE NOT NULL, 
	  updated_at DATE NOT NULL
	) IF NOT EXISTS`

	_, err := db.Query(createUserTable)
	if err != nil {
		return fmt.Errorf("failed To Create User Table %w", err)
	}

	_, err = db.Query(createBlogsTable)
	if err != nil {
		return fmt.Errorf("failed to create posts table %w", err)
	}

	_, err = db.Query(createCommentsTable)
	if err != nil {
		return fmt.Errorf("failed to create comments table %w", err)
	}

	return nil
}
