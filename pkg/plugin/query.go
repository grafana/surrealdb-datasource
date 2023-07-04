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
func (d *SurrealDatasource) createQuery(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	r, err := d.db.Query(qm.Text, nil)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("query: %v", err.Error()))
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/
	var response backend.DataResponse

	// unmarshal the response into a slice of maps.
	// each map represents a row in the response from the database
	// as a map of column name to value. The value is a `json.RawMessage`.
	var res []map[string]json.RawMessage
	ok, err := surrealdb.UnmarshalRaw(r, &res)

	if !ok {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unmarshal: %v", err.Error()))
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
