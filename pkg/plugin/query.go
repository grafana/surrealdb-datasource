package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/surrealdb/surrealdb.go"
)

// createDataResponse creates a data response from a data query.
func (d *SurrealDatasource) createDataResponse(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	str, err := sqlStringFromDataQuery(query)
	if err != nil {
		return backend.ErrDataResponseWithSource(backend.StatusBadRequest, backend.ErrorSourcePlugin, fmt.Sprintf("sql: %v", err.Error()))
	}

	result, err := d.client.QueryWithContext(ctx, str, nil)
	if err != nil {
		return backend.ErrDataResponseWithSource(backend.StatusBadRequest, backend.ErrorSourceDownstream, fmt.Sprintf("query: %v", err.Error()))
	}

	response, err := buildResponse(result)
	if err != nil {
		return backend.ErrDataResponseWithSource(backend.StatusBadRequest, backend.ErrorSourcePlugin, fmt.Sprintf("response: %v", err.Error()))
	}

	return response
}

// sqlStringFromDataQuery converts a data query into a SQL string, interpolating any macros.
func sqlStringFromDataQuery(query backend.DataQuery) (string, error) {
	sq, err := sqlutil.GetQuery(query)
	if err != nil {
		return "", err
	}

	// apply grafana macros to the query
	str, err := sqlutil.Interpolate(sq, sqlutil.DefaultMacros)
	if err != nil {
		return "", err
	}

	return str, nil
}

// buildResponse converts the response from the database into a data response.
func buildResponse(result interface{}) (backend.DataResponse, error) {
	var response backend.DataResponse

	// unmarshal the response into a slice of maps.
	// each map represents a row in the response from the database
	// as a map of column name to value. The value is a `json.RawMessage`.
	var res []map[string]json.RawMessage
	ok, err := surrealdb.UnmarshalRaw(result, &res)
	if err != nil {
		return response, err
	}

	if !ok {
		return response, nil
	}

	// convert the response to a data frame.
	frame := toDataFrame(res)

	// add the fields to the frame.
	response.Frames = append(response.Frames, frame)

	return response, nil
}

// toDataFrame converts the response from the database into a data frame.
func toDataFrame(resp []map[string]json.RawMessage) *data.Frame {
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
