package config

// PostgreSQLConfig struct holds PostgreSQL configuration parameters
type PostgreSQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewPostgreSQLConfig creates a new PostgreSQLConfig instance with default values
func NewPostgreSQLConfig() *PostgreSQLConfig {
	return &PostgreSQLConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Database: "dbname",
	}
}
