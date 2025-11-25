package database

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/olekukonko/tablewriter"
)

// query → returns rows (SELECT, SHOW, etc.)


type QueryResult struct {
	rowStreamer
}

func (q *QueryResult) RenderTable() string {
	defer q.Close()

	var sb strings.Builder
	table := tablewriter.NewWriter(&sb)

	cols := q.Columns()
	if len(cols) == 0 {
		return "(no columns)"
	}
	table.Header(cols)

	var rows [][]string
	for {
		row, err := q.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Sprintf("Error reading rows: %v", err)
		}

		strRow := make([]string, len(row))
		for i, v := range row {
			if v == nil {
				strRow[i] = "NULL"
			} else {
				strRow[i] = fmt.Sprintf("%v", v)
			}
		}
		rows = append(rows, strRow)
	}
	table.Bulk(rows)
	table.Render()

	return sb.String()
}

type rowStreamer struct {
	rows    pgx.Rows
	columns []string
	closed  bool
	duration time.Duration
}

func (r *rowStreamer) Columns() []string {
	return r.columns
}

// Next returns the next row as []interface{} or io.EOF when done.
func (r *rowStreamer) Next() ([]interface{}, error) {
	if r.closed {
		return nil, io.EOF
	}
	if r.rows.Next() {
		vals, err := r.rows.Values()
		if err != nil {
			r.rows.Close()
			r.closed = true
			return nil, err
		}
		// convert []interface{} as-is; nil for NULLs
		return vals, nil
	}
	if err := r.rows.Err(); err != nil {
		r.rows.Close()
		r.closed = true
		return nil, err
	}
	// no more rows
	r.rows.Close()
	r.closed = true
	return nil, io.EOF
}

func (r *rowStreamer) Close() error {
	if r.closed {
		return nil
	}
	r.rows.Close()
	r.closed = true
	return nil
}

func (r *rowStreamer) Duration() time.Duration {
	return r.duration
}

func (r *rowStreamer) GetType() string {
	return "QUERY"
}


