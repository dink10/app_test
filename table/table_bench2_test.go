package table

import (
	"bytes"
	"testing"

	"github.com/olekukonko/tablewriter"
)

var (
	headers = []string{"Name", "Sign", "Rating"}
	data    = [][]string{
		{"A", "The Good", "500"},
		{"B", "The Very very Bad Man", "288"},
		{"C", "The Ugly", "120"},
	}
)

func BenchmarkThirdPartyTable(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		buf.Reset()
		table := tablewriter.NewWriter(&buf)
		table.SetHeader(headers)

		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	}
}
