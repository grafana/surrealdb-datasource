package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/surrealdb/surrealdb.go"
)

type queryModel struct {
	Text string `json:"queryText"`
}

// query is the internal implementation of executing queries in the Surreal datasource.
func (d *SurrealDatasource) createQuery(ctx context.Context, _ backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponseWithSource(backend.StatusBadRequest, backend.ErrorSourcePlugin, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	r, err := d.queryWithContext(ctx, qm.Text, nil)
	if err != nil {
		return backend.ErrDataResponseWithSource(backend.StatusBadRequest, backend.ErrorSourceDownstream, fmt.Sprintf("query: %v", err.Error()))
	}

	var response backend.DataResponse

	// unmarshal the response into a slice of maps.
	// each map represents a row in the response from the database
	// as a map of column name to value. The value is a `json.RawMessage`.
	var res []map[string]json.RawMessage
	ok, err := surrealdb.UnmarshalRaw(r, &res)

	if err != nil {
		return backend.ErrDataResponseWithSource(backend.StatusBadRequest, backend.ErrorSourcePlugin, fmt.Sprintf("unmarshal: %v", err.Error()))
	}

	if !ok {
		return response
	}

	// convert the response to a data frame.
	frame := d.toDataFrame(res)

	// add the fields to the frame.
	response.Frames = append(response.Frames, frame)

	return response
}

// toDataFrame converts the response from the database into a data frame.
func (d *SurrealDatasource) toDataFrame(resp []map[string]json.RawMessage) *data.Frame {
	// @adamyeats: TODO: what should the name here be?
	frame := data.NewFrame("response")

	buckets := map[string][]json.RawMessage{}

	if len(resp) == 0 {
		return frame
	}

	for _, entity := range resp {
		for k, v := range entity {
			buckets[k] = append(buckets[k], v)
		}
	}

	for key, vals := range buckets {
		field := data.NewField(key, nil, vals)
		frame.Fields = append(frame.Fields, field)
	}

	return frame
}

// QueryWithContext wraps the Query method to handle context for cancellation/timeout
func (d *SurrealDatasource) queryWithContext(ctx context.Context, query string, args interface{}) (interface{}, error) {
	rc := make(chan interface{})
	ec := make(chan error)

	go func() {
		r, err := d.db.Query(query, args)
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
