package integration_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/grafana-labs/surrealdb-datasource/pkg/client"
	"github.com/grafana-labs/surrealdb-datasource/pkg/plugin"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/slo"
)

// table is a list of test cases for the QueryData method.
var cases = []struct {
	input string
	name  string
}{
	{
		input: "SELECT * FROM person WHERE name = 'Aaron Brooks';",
		name:  "SELECT query with WHERE clause",
	},
	{
		input: "SELECT id, name FROM person;",
		name:  "SELECT specific columns",
	},
	{
		input: "SELECT count() FROM person GROUP all;",
		name:  "Aggregate function",
	},
	{
		input: "SELECT * FROM person ORDER BY name DESC;",
		name:  "SELECT query with ORDER BY clause",
	},
	{
		input: "SELECT * FROM person WHERE first_name ~ 'aaron';",
		name:  "SELECT with fuzzy equality",
	},
	{
		input: "SELECT * FROM person WHERE company_name == NONE;",
		name:  "SELECT with IS NONE",
	},
	{
		input: "SELECT * FROM person WHERE company_name != NONE;",
		name:  "SELECT with IS NOT NONE",
	},
	{
		input: "SELECT * FROM person LIMIT 10;",
		name:  "SELECT with LIMIT",
	},
	{
		input: "SELECT out.name FROM review;",
		name:  "SELECT with Record links",
	},
	{
		input: "SELECT * FROM person WHERE time.created_at = d'2023-11-27T22:30:23Z'",
		name:  "SELECT with datetime equality",
	},
	{
		input: "SELECT * FROM person WHERE time.created_at > d'2022-11-27T22:30:23Z'",
		name:  "SELECT with datetime greater than",
	},
	{
		input: "SELECT * FROM person WHERE time.created_at < d'2023-11-27T22:30:23Z'",
		name:  "SELECT with datetime less than",
	},
}

// getSurrealEndpoint returns the endpoint for the SurrealDB server.
func getSurrealEndpoint() string {
	env := os.Getenv("SURREAL_DB_URL")

	if env != "" {
		return env
	}

	return "ws://localhost:8000/rpc"
}

// createTestInstance creates a new SurrealDatasource instance for testing.
func createTestInstance() (instancemgmt.Instance, error) {
	config := client.SurrealConfig{
		Database:  "test",
		Endpoint:  getSurrealEndpoint(),
		Namespace: "test",
		Username:  "root",
	}

	msg, err := json.Marshal(config)

	if err != nil {
		return nil, err
	}

	dsiConfig := backend.DataSourceInstanceSettings{
		JSONData:                msg,
		DecryptedSecureJSONData: map[string]string{"password": "test"},
	}

	return plugin.NewDatasource(context.Background(), dsiConfig)
}

// createJsonRequest creates a new JSON request for testing.
func createJsonRequest(query string) json.RawMessage {
	return json.RawMessage(`{"rawSql":"` + query + `"}`)
}

func TestNewDatasource(t *testing.T) {
	instance, err := createTestInstance()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if instance == nil {
		t.Errorf("expected instance to be non-nil")
	}

	if _, ok := instance.(*plugin.SurrealDatasource); ok {
		t.Errorf("expected instance to be the correct type")
	}
}

func TestCheckHealth(t *testing.T) {
	instance, err := createTestInstance()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	res, err := (*instance.(*slo.MetricsWrapper)).CheckHealth(context.Background(), &backend.CheckHealthRequest{})

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if res == nil {
		t.Errorf("expected response to be non-nil")
	}
}

func TestQueryData_Success(t *testing.T) {
	instance, err := createTestInstance()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			dqs := []backend.DataQuery{
				{
					RefID: "A",
					JSON:  createJsonRequest(tt.input),
				},
			}

			req := backend.QueryDataRequest{
				PluginContext: backend.PluginContext{
					OrgID: 1,
				},
				Queries: dqs,
			}

			res, err := (*instance.(*slo.MetricsWrapper)).QueryData(context.Background(), &req)

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if res == nil {
				t.Errorf("expected response to be non-nil")
			} else {
				if len(res.Responses) != 1 {
					t.Errorf("expected 1 response, got %d", len(res.Responses))
				}

				if res.Responses["A"].Frames == nil {
					t.Errorf("expected frames to be non-nil")
				}

				if res.Responses["A"].Frames[0].Fields == nil {
					t.Errorf("expected fields to be non-nil")
				}
			}
		})
	}
}

func TestQueryData_BadTableName(t *testing.T) {
	instance, err := createTestInstance()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	dqs := []backend.DataQuery{
		{
			RefID: "A",
			JSON:  createJsonRequest("SELECT * FROM does_not_exist;"),
		},
	}

	req := backend.QueryDataRequest{
		PluginContext: backend.PluginContext{
			OrgID: 1,
		},
		Queries: dqs,
	}

	res, err := (*instance.(*slo.MetricsWrapper)).QueryData(context.Background(), &req)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if len(res.Responses) != 1 {
			t.Errorf("expected 1 response, got %d", len(res.Responses))
		}

		if res.Responses["A"].Frames != nil {
			t.Errorf("expected frames to be nil")
		}
	}
}

func TestQueryData_BadQuerySyntax(t *testing.T) {
	instance, err := createTestInstance()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	dqs := []backend.DataQuery{
		{
			RefID: "A",
			JSON:  createJsonRequest("SELECT * FROM person JOIN address ON person.id = address.person_id;"),
		},
	}

	req := backend.QueryDataRequest{
		PluginContext: backend.PluginContext{
			OrgID: 1,
		},
		Queries: dqs,
	}

	res, err := (*instance.(*slo.MetricsWrapper)).QueryData(context.Background(), &req)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if len(res.Responses) != 1 {
			t.Errorf("expected 1 response, got %d", len(res.Responses))
		}

		if res.Responses["A"].Frames != nil {
			t.Errorf("expected frames to be nil")
		}

		if res.Responses["A"].Error == nil {
			t.Errorf("expected error to be non-nil")
		}
	}
}
