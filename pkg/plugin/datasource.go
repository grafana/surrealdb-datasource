package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana-labs/surrealdb-datasource/pkg/client"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/surrealdb/surrealdb.go"
)

var (
	_ backend.QueryDataHandler      = (*SurrealDatasource)(nil)
	_ backend.CheckHealthHandler    = (*SurrealDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*SurrealDatasource)(nil)
)

// SurrealDatasource defines how to connect to the datasource and describes the query model.
type SurrealDatasource struct {
	db     *surrealdb.DB
	config *client.SurrealConfig
}

// NewDatasource creates a new datasource instance.
func NewDatasource(dsiConfig backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var config client.SurrealConfig
	err := json.Unmarshal(dsiConfig.JSONData, &config)

	config.Password = dsiConfig.DecryptedSecureJSONData["password"]

	if err != nil {
		return nil, fmt.Errorf("unable to get settings from JSON config: %w", err)
	}

	db, err := client.NewConnection(config)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &SurrealDatasource{
		db:     db,
		config: &config,
	}, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *SurrealDatasource) Dispose() {
	// Clean up datasource instance resources.
	d.db.Close()
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *SurrealDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct{}

// query is the internal implementation of executing queries in the Surreal datasource.
func (d *SurrealDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/
	frame := data.NewFrame("response")

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("values", nil, []int64{10, 20}),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
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
