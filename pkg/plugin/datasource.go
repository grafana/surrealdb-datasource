package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/grafana-labs/surrealdb-datasource/pkg/client"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/surrealdb/surrealdb.go"
)

var (
	_ backend.QueryDataHandler      = (*SurrealDatasource)(nil)
	_ backend.CheckHealthHandler    = (*SurrealDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*SurrealDatasource)(nil)
)

// SurrealDatasource defines how to connect to the datasource and describes the query model.
type SurrealDatasource struct {
	db     client.SurrealDBClient
	config *client.SurrealConfig
}

// NewDatasourceInstance creates a new SurrealDatasource instance.
func NewDatasourceInstance(db client.SurrealDBClient, config *client.SurrealConfig) *SurrealDatasource {
	return &SurrealDatasource{
		db:     db,
		config: config,
	}
}

// NewDatasource creates a new datasource instance.
func NewDatasource(dsiConfig backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var config client.SurrealConfig
	err := json.Unmarshal(dsiConfig.JSONData, &config)

	config.Password = dsiConfig.DecryptedSecureJSONData["password"]

	if err != nil {
		return nil, fmt.Errorf("unable to get settings from JSON config: %w", err)
	}

	db, err := surrealdb.New(config.Endpoint)

	if err != nil {
		return nil, err
	}

	_, err = client.Use(db).Connect(&config)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return NewDatasourceInstance(db, &config), nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
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

			res := d.createQuery(ctx, pluginCtx, q)

			mutex.Lock()
			response.Responses[q.RefID] = res
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
func (d *SurrealDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	status := backend.HealthStatusOk
	message := "Data source is working"

	result := make(chan error)

	go func() {
		_, err := d.db.Query("BEGIN TRANSACTION; CANCEL TRANSACTION;", nil)
		result <- err
	}()

	err := <-result

	if err != nil {
		status = backend.HealthStatusError
		message = fmt.Sprintf("error while checking database health: %v", err)
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
