package client

import "github.com/surrealdb/surrealdb.go"

// SurrealConfig defines the configuration for the SurrealDB database.
type SurrealConfig struct {
	Database  string `json:"database,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Password  string `json:"password,omitempty"`
	Username  string `json:"username,omitempty"`
}

// NewConnection creates a new connection to the SurrealDB database.
func NewConnection(config *SurrealConfig) (*surrealdb.DB, error) {
	db, err := surrealdb.New(config.Endpoint)

	if err != nil {
		return nil, err
	}

	credentials := map[string]interface{}{
		"user": config.Username,
		"pass": config.Password,
	}

	if _, err = db.Signin(credentials); err != nil {
		return nil, err
	}

	if _, err = db.Use(config.Namespace, config.Database); err != nil {
		return nil, err
	}

	return db, nil
}
