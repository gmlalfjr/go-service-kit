package db

import (
	"database/sql"
	"fmt"

	"github.com/gmlalfjr/go-service-kit/config"
	"github.com/gmlalfjr/go-service-kit/logger"
)

// PostgreSQLService represents a PostgreSQL service
type PostgreSQLService struct {
	db     *sql.DB
	logger logger.LoggerConfig
}

// NewPostgreSQLService creates a new PostgreSQLService instance
func NewPostgreSQLService(config *config.PostgreSQLConfig, logger logger.LoggerConfig) (*PostgreSQLService, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLService{
		db:     db,
		logger: logger,
	}, nil
}

// ...
