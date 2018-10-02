package table

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	testHeaders = []string{ColumnInt, ColumnString, ColumnMoney}
	testRows    = [][]string{
		{"1", "aaaaaa bbb ccc", "1000.33"},
		{"7", "aaaaaa bbb", "0.01"},
		{"11", "aaaaaa", "10000"},
	}
)

func BenchmarkTable(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		table := New(len(testHeaders), false)
		for hIndex, header := range testHeaders {
			table.SetColumn(hIndex, header, AlignLeft)
		}
		buf.Reset()
		if err := table.Append(testRows...); err != nil {
			logrus.Fatal(err)
		}

		table.Write(&buf)
	}
}
