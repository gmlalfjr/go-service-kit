package postgre

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgreService struct {
	DB *sql.DB
}

func NewPostgreService(connStr string) *PostgreService {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	return &PostgreService{
		DB: db,
	}
}

func (p *PostgreService) Start() error {
	if err := p.DB.Ping(); err != nil {
		return fmt.Errorf("Failed to ping PostgreSQL: %w", err)
	}
	log.Println("[PostgreSQL] Connection established")
	return nil
}

func (p *PostgreService) Stop() error {
	if err := p.DB.Close(); err != nil {
		return fmt.Errorf("Failed to close PostgreSQL connection: %w", err)
	}
	log.Println("[PostgreSQL] Connection closed")
	return nil
}
