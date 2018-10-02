package table

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestTableErrors(t *testing.T) {
	var err error
	table := New(3, true)

	err = table.SetColumn(-1, ColumnString, AlignCenter)
	if err == nil {
		t.Error("expected SetColumn with col arg -1 to fail")
	}
	err = table.SetColumn(3, ColumnString, AlignCenter)
	if err == nil {
		t.Error("expected SetColumn with col arg 3 to fail")
	}
	err = table.SetColumn(2, ColumnString, AlignCenter)
	if err != nil {
		t.Errorf("expected SetColumn with col arg 2 to succeed; got %v", err)
	}

	specs := [][]string{
		{""},
		{"", "", "", ""},
	}
	for specIndex, spec := range specs {
		err = table.Append(spec)
		if err == nil {
			t.Errorf("[spec %d] expected call to Append with %d column(s) (table has 3 columns) to fail", specIndex, len(spec))
		}
	}
}

func TestTableAppendErrors(t *testing.T) {
	var err error
	table := New(3, false)
	err = table.SetColumn(0, ColumnInt, AlignLeft)
	err = table.SetColumn(1, ColumnInt, AlignCenter)
	err = table.SetColumn(2, ColumnInt, AlignCenter)
	specs := [][]string{
		{"", "", ""},
	}

	for _, spec := range specs {
		err = table.Append(spec)
		if err == nil {
			t.Errorf("expected call to Append to INT column(s) string value to fail")
		}
	}
}

func TestTable(t *testing.T) {
	testHeaders := []string{ColumnInt, ColumnString, ColumnMoney, ColumnString}
	testRows := [][]string{
		{"1", "aaaaaa bbb ccc", "1000.33", "aaaaaa bbb ccc tttttttt"},
		{"7", "aaaaaa bbb", "0.01", "aaaaaa ccc tttttttt"},
		{"11", "aaaaaa", "10000", "tttttttt"},
	}

	specs := []struct {
		alignments []Alignment
		rows       [][]string
		expOutput  string
	}{
		{
			alignments: []Alignment{AlignLeft, AlignLeft, AlignRight, AlignLeft},
			rows:       testRows,
			expOutput: `+--+------+---------+--------+
|1 |aaaaaa| 1 000,33|aaaaaa  |
|  |bbb   |         |bbb     |
|  |ccc   |         |ccc     |
+--+------+---------+--------+
|7 |aaaaaa|     0,01|aaaaaa  |
|  |bbb   |         |ccc     |
+--+------+---------+--------+
|11|aaaaaa|10 000,00|tttttttt|
+--+------+---------+--------+
`,
		},
		{
			alignments: []Alignment{AlignLeft, AlignCenter, AlignRight, AlignCenter},
			rows:       testRows,
			expOutput: `+--+------+---------+--------+
|1 |aaaaaa| 1 000,33| aaaaaa |
|  | bbb  |         |  bbb   |
|  | ccc  |         |  ccc   |
+--+------+---------+--------+
|7 |aaaaaa|     0,01| aaaaaa |
|  | bbb  |         |  ccc   |
+--+------+---------+--------+
|11|aaaaaa|10 000,00|tttttttt|
+--+------+---------+--------+
`,
		},
	}

	var buf bytes.Buffer
	for specIndex, spec := range specs {
		table := New(len(testHeaders), false)
		for hIndex, header := range testHeaders {
			table.SetColumn(hIndex, header, spec.alignments[hIndex])
		}
		if err := table.Append(spec.rows...); err != nil {
			t.Error(err)
		}

		buf.Reset()
		table.Write(&buf)
		tableOutput := buf.String()

		if tableOutput != spec.expOutput {
			t.Errorf(
				"[spec %d]\n expected output to be: \n%s\n got:\n%s\n",
				specIndex,
				indent(spec.expOutput, 2),
				indent(tableOutput, 2),
			)
		}
	}
}

func TestTableWithAllowingSpaces(t *testing.T) {
	testHeaders := []string{ColumnInt, ColumnString, ColumnMoney, ColumnString}
	testRows := [][]string{
		{"1", "aaaaaa bbb ccc", "1000.33", "aaaaaa bbb ccc tttttttt"},
		{"7", "aaaaaa bbb", "", ""},
		{"", "aaaaaa", "10000", "tttttttt"},
	}

	specs := []struct {
		alignments []Alignment
		rows       [][]string
		expOutput  string
	}{
		{
			alignments: []Alignment{AlignLeft, AlignLeft, AlignRight, AlignLeft},
			rows:       testRows,
			expOutput: `+-+------+---------+--------+
|1|aaaaaa| 1 000,33|aaaaaa  |
| |bbb   |         |bbb     |
| |ccc   |         |ccc     |
+-+------+---------+--------+
|7|aaaaaa|         |        |
| |bbb   |         |        |
+-+------+---------+--------+
| |aaaaaa|10 000,00|tttttttt|
+-+------+---------+--------+
`,
		},
		{
			alignments: []Alignment{AlignLeft, AlignCenter, AlignRight, AlignCenter},
			rows:       testRows,
			expOutput: `+-+------+---------+--------+
|1|aaaaaa| 1 000,33| aaaaaa |
| | bbb  |         |  bbb   |
| | ccc  |         |  ccc   |
+-+------+---------+--------+
|7|aaaaaa|         |        |
| | bbb  |         |        |
+-+------+---------+--------+
| |aaaaaa|10 000,00|tttttttt|
+-+------+---------+--------+
`,
		},
	}

	var buf bytes.Buffer
	for specIndex, spec := range specs {
		table := New(len(testHeaders), true)
		for hIndex, header := range testHeaders {
			table.SetColumn(hIndex, header, spec.alignments[hIndex])
		}
		if err := table.Append(spec.rows...); err != nil {
			t.Error(err)
		}

		buf.Reset()
		table.Write(&buf)
		tableOutput := buf.String()

		if tableOutput != spec.expOutput {
			t.Errorf(
				"[spec %d]\n expected output to be: \n%s\n got:\n%s\n",
				specIndex,
				indent(spec.expOutput, 2),
				indent(tableOutput, 2),
			)
		}
	}
}

func indent(s string, pad int) string {
	return regexp.MustCompile("(?m)^").ReplaceAllString(s, strings.Repeat(" ", pad))
}
