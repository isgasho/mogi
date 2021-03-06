package mogi

import (
	"database/sql/driver"
	"encoding/csv"
	"io"
	"strings"
	"time"
)

type rows struct {
	cols []string
	data [][]driver.Value

	cursor int
	closed bool
}

func newRows(cols []string, data [][]driver.Value) *rows {
	return &rows{
		cols: cols,
		data: data,
	}
}

func (r *rows) Columns() []string {
	return r.cols
}

// Close closes the rows iterator.
func (r *rows) Close() error {
	r.closed = true
	return nil
}

func (r *rows) Err() error {
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	r.cursor++
	if r.cursor > len(r.data) {
		r.closed = true
		return io.EOF
	}

	for i, col := range r.data[r.cursor-1] {
		dest[i] = col
	}

	return nil
}

// cribbed from DATA-DOG/go-sqlmock
// TODO rewrite
func csvToValues(cols []string, s string) [][]driver.Value {
	var data [][]driver.Value
	if s == "" {
		return nil
	}

	res := strings.NewReader(strings.TrimSpace(s))
	csvReader := csv.NewReader(res)

	for {
		res, err := csvReader.Read()
		if err != nil || res == nil {
			break
		}

		row := []driver.Value{}
		for _, v := range res {
			if timeLayout != "" {
				if t, err := time.Parse(timeLayout, v); err == nil {
					row = append(row, t)
					continue
				}
			}
			row = append(row, []byte(strings.TrimSpace(v)))
		}
		data = append(data, row)
	}
	return data
}
