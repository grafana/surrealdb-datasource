package plugin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/grafana-labs/surrealdb-datasource/internal/mocks"
	"github.com/grafana-labs/surrealdb-datasource/pkg/client"
	"github.com/grafana-labs/surrealdb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestCreateDataResponse_QueryError(t *testing.T) {
	query := backend.DataQuery{
		RefID: "A",
		JSON:  []byte(`{"rawSql": "SELECT * FROM test"}`),
	}

	errorMock := mocks.MockSurrealDBClient{
		QueryFunc: func(sql string, vars interface{}) (interface{}, error) {
			return nil, errors.New("query error")
		},
	}

	ds := plugin.NewDatasourceInstance(client.Use(&errorMock), &config)

	ctx := context.TODO()
	response := ds.CreateDataResponse(ctx, query)

	if response.Error == nil {
		t.Error("expected error, got nil")
	}
	if response.Status != backend.StatusBadRequest {
		t.Errorf("expected status bad request, got %v", response.Status)
	}
	if response.Error.Error() != "query: query error" {
		t.Errorf("expected error message 'query: query error', got %v", response.Error.Error())
	}
}

func TestCreateDataResponse_BuildResponseError(t *testing.T) {
	query := backend.DataQuery{
		RefID: "A",
		JSON:  []byte(`{"rawSql": "SELECT * FROM test"}`),
	}

	errorMock := mocks.MockSurrealDBClient{
		QueryFunc: func(sql string, vars interface{}) (interface{}, error) {
			return "invalid response", nil
		},
	}

	ds := plugin.NewDatasourceInstance(client.Use(&errorMock), &config)

	ctx := context.TODO()
	response := ds.CreateDataResponse(ctx, query)

	if response.Error == nil {
		t.Error("expected error, got nil")
	}
	if response.Status != backend.StatusBadRequest {
		t.Errorf("expected status bad request, got %v", response.Status)
	}
	if response.Error.Error() != "response: failed raw unmarshaling to interface slice: invalid SurrealDB response" {
		t.Errorf("expected error message 'response: failed raw unmarshaling to interface slice: invalid SurrealDB response', got %v", response.Error.Error())
	}
}

func TestCreateDataResponse_Success(t *testing.T) {
	query := backend.DataQuery{
		RefID: "A",
		JSON:  []byte(`{"rawSql": "SELECT * FROM test"}`),
	}

	successMock := mocks.MockSurrealDBClient{
		QueryFunc: func(sql string, vars interface{}) (interface{}, error) {
			return []interface{}{
				map[string]interface{}{
					"status": "OK",
					"result": []interface{}{
						map[string]interface{}{
							"column1": "value1",
							"column2": "value2",
						},
					},
					"time": "10ms",
				},
			}, nil
		},
	}

	ds := plugin.NewDatasourceInstance(client.Use(&successMock), &config)

	ctx := context.TODO()
	response := ds.CreateDataResponse(ctx, query)

	if response.Error != nil {
		t.Errorf("unexpected error: %s", response.Error)
	}
	if response.Frames == nil {
		t.Error("expected frames to be non-nil")
	}
	if len(response.Frames) != 1 {
		t.Errorf("expected 1 frame, got %d", len(response.Frames))
	}
	if response.Frames[0].Fields[0].Name != "column1" {
		t.Errorf("expected first field name to be 'column1', got %v", response.Frames[0].Fields[0].Name)
	}
	if response.Frames[0].Fields[1].Name != "column2" {
		t.Errorf("expected second field name to be 'column2', got %v", response.Frames[0].Fields[1].Name)
	}
}
