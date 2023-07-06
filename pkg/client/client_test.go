package client_test

import (
	"errors"
	"testing"

	"github.com/grafana-labs/surrealdb-datasource/internal/mocks"
	"github.com/grafana-labs/surrealdb-datasource/pkg/client"
)

func TestConnect_Success(t *testing.T) {
	mockDB := mocks.MockSurrealDBClient{
		SigninFunc: func(vars interface{}) (interface{}, error) {
			return nil, nil
		},
		UseFunc: func(namespace string, database string) (interface{}, error) {
			return nil, nil
		},
	}

	c := client.Use(&mockDB)

	config := client.SurrealConfig{
		Database:  "test_db",
		Endpoint:  "ws://localhost:8000",
		Namespace: "test-namespace",
		Password:  "password",
		Username:  "username",
	}

	ok, err := c.Connect(&config)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if ok != true {
		t.Error("expected result to be successful, got false")
	}
}

func TestConnect_SigninError(t *testing.T) {
	mockDB := mocks.MockSurrealDBClient{
		SigninFunc: func(vars interface{}) (interface{}, error) {
			return nil, errors.New("signin error")
		},
	}

	c := client.Use(&mockDB)

	config := client.SurrealConfig{
		Database:  "test_db",
		Endpoint:  "ws://localhost:8000",
		Namespace: "test-namespace",
		Password:  "password",
		Username:  "username",
	}

	_, err := c.Connect(&config)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestConnect_UseError(t *testing.T) {
	mockDB := mocks.MockSurrealDBClient{
		SigninFunc: func(vars interface{}) (interface{}, error) {
			return nil, nil
		},
		UseFunc: func(namespace string, database string) (interface{}, error) {
			return nil, errors.New("use error")
		},
	}

	c := client.Use(&mockDB)

	config := client.SurrealConfig{
		Database:  "test-db",
		Endpoint:  "ws://localhost:8000",
		Namespace: "test-namespace",
		Password:  "password",
		Username:  "username",
	}

	_, err := c.Connect(&config)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
