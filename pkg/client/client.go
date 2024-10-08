package client

import "context"

// SurrealConfig defines the configuration for the SurrealDB database.
type SurrealConfig struct {
	Database  string `json:"database,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Password  string `json:"password,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Username  string `json:"username,omitempty"`
}

// SurrealDBClient defines the interface for the SurrealDB database.
type SurrealDBClient interface {
	Close()
	Create(thing string, data interface{}) (interface{}, error)
	Query(sql string, vars interface{}) (interface{}, error)
	Signin(vars interface{}) (interface{}, error)
	Use(namespace string, database string) (interface{}, error)
}

// Client defines the client for the SurrealDB database.
type Client struct {
	db SurrealDBClient
}

// Use returns a new client for the SurrealDB database.
func Use(db SurrealDBClient) *Client {
	return &Client{db}
}

// Connect connects to the SurrealDB database.
func (c *Client) Connect(config *SurrealConfig) (bool, error) {
	credentials := map[string]interface{}{
		"user": config.Username,
		"pass": config.Password,
	}

	if config.Scope != "" {
		credentials["scope"] = config.Scope
	}

	if _, err := c.db.Signin(credentials); err != nil {
		return false, err
	}

	if _, err := c.db.Use(config.Namespace, config.Database); err != nil {
		return false, err
	}

	return true, nil
}

// QueryWithContext wraps the Query method to handle context for cancellation/timeout
func (c *Client) QueryWithContext(ctx context.Context, query string, args interface{}) (interface{}, error) {
	rc := make(chan interface{})
	ec := make(chan error)

	go func() {
		r, err := c.db.Query(query, args)
		if err != nil {
			ec <- err
			return
		}
		rc <- r
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-ec:
		return nil, err
	case result := <-rc:
		return result, nil
	}
}
