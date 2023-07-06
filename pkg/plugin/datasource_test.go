package plugin_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/grafana-labs/surrealdb-datasource/internal/mocks"
	"github.com/grafana-labs/surrealdb-datasource/pkg/client"
	"github.com/grafana-labs/surrealdb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// IMPORTANT!
//
// As it stands, this test suite requires a running instance of SurrealDB. You
// can start one by bringing up the Docker containers with the following command:
//
// `docker-compose up -d`
//
// What can be tested in isolation at the moment, is being tested in isolation.

var mock = mocks.MockSurrealDBClient{}

var config = client.SurrealConfig{
	Database:  "grafana_ds_tests",
	Endpoint:  "ws://localhost:8000/rpc",
	Namespace: "grafana",
	Username:  "grafana",
}

func TestNewDatasource_Success(t *testing.T) {
	msg, err := json.Marshal(config)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	config := backend.DataSourceInstanceSettings{
		JSONData:                msg,
		DecryptedSecureJSONData: map[string]string{"password": "password"},
	}

	instance, err := plugin.NewDatasource(config)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if instance == nil {
		t.Error("expected instance to be non-nil")
	}
}

func TestNewDatasource_InvalidJSON(t *testing.T) {
	config := backend.DataSourceInstanceSettings{
		JSONData:                json.RawMessage(`invalid json`),
		DecryptedSecureJSONData: map[string]string{"password": "password"},
	}

	_, err := plugin.NewDatasource(config)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestQueryData(t *testing.T) {
	// Create a mock SurrealDatasource with a mock DB client
	datasource := plugin.NewDatasourceInstance(&mock, &config)

	mockRequest := backend.QueryDataRequest{
		Queries: []backend.DataQuery{
			{RefID: "query1"},
			{RefID: "query2"},
		},
	}

	// this returns a DataResponse no matter what, and that includes errors
	response, err := datasource.QueryData(context.Background(), &mockRequest)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if response == nil {
		t.Error("expected response to be non-nil")
	}
}

func TestCheckHealth_Success(t *testing.T) {
	successMock := mocks.MockSurrealDBClient{
		QueryFunc: func(sql string, vars interface{}) (interface{}, error) {
			return nil, nil
		},
	}

	datasource := plugin.NewDatasourceInstance(&successMock, &config)
	result, err := datasource.CheckHealth(context.Background(), &backend.CheckHealthRequest{})

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if result == nil {
		t.Error("expected result to be non-nil")
	}
}

func TestCheckHealth_Error(t *testing.T) {
	errorMock := mocks.MockSurrealDBClient{
		QueryFunc: func(sql string, vars interface{}) (interface{}, error) {
			return nil, errors.New("query error")
		},
	}

	datasource := plugin.NewDatasourceInstance(&errorMock, &config)
	result, err := datasource.CheckHealth(context.Background(), &backend.CheckHealthRequest{})

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if result == nil {
		t.Error("expected result to be non-nil")
	}
}
