package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/grafana-labs/surrealdb-datasource/pkg/client"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/slo"
	"github.com/surrealdb/surrealdb.go"
)

var (
	_ backend.QueryDataHandler      = (*SurrealDatasource)(nil)
	_ backend.CheckHealthHandler    = (*SurrealDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*SurrealDatasource)(nil)
)

// SurrealDatasource defines how to connect to the datasource and describes the query model.
type SurrealDatasource struct {
	client *client.Client
	config *client.SurrealConfig
}

// NewDatasourceInstance creates a new SurrealDatasource instance.
func NewDatasourceInstance(client *client.Client, config *client.SurrealConfig) *SurrealDatasource {
	return &SurrealDatasource{
		client: client,
		config: config,
	}
}

// NewDatasource creates a new datasource instance.
func NewDatasource(ctx context.Context, dsiConfig backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var config client.SurrealConfig

	err := json.Unmarshal(dsiConfig.JSONData, &config)
	if err != nil {
		return nil, fmt.Errorf("unable to get settings from JSON config: %w", err)
	}

	config.Password = dsiConfig.DecryptedSecureJSONData["password"]

	db, err := surrealdb.New(config.Endpoint)
	if err != nil {
		return nil, errorsource.DownstreamError(err, false)
	}

	client := client.Use(db)

	_, err = client.Connect(&config)
	if err != nil {
		return nil, errorsource.DownstreamError(fmt.Errorf("unable to connect to database: %w", err), false)
	}

	return slo.NewMetricsWrapper(NewDatasourceInstance(client, &config), dsiConfig), nil
}

// Dispose cleans up the datasource instance resources.
func (d *SurrealDatasource) Dispose() {
	// Clean up datasource instance resources.
	// d.db.Close()
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *SurrealDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, query := range req.Queries {
		wg.Add(1)

		go func(ctx context.Context, pluginCtx backend.PluginContext, q backend.DataQuery) {
			defer wg.Done()

			mutex.Lock()
			response.Responses[q.RefID] = d.CreateDataResponse(ctx, q)
			mutex.Unlock()
		}(ctx, req.PluginContext, query)
	}

	wg.Wait()

	return response, nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *SurrealDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	status := backend.HealthStatusOk
	message := "Data source is working"

	_, err := d.client.QueryWithContext(ctx, "BEGIN TRANSACTION; CANCEL TRANSACTION;", nil)

	if err != nil {
		status = backend.HealthStatusError
		message = fmt.Sprintf("error while checking database health: %v", err)
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
