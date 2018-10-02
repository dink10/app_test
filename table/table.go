package table

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/sirupsen/logrus"
)

// Alignment represents the supported cell content alignment modes.
type Alignment uint8

const (
	// Format constants.
	Center  byte = '+'
	Column  byte = '|'
	Row     byte = '-'
	Space        = " "
	NewLine      = "\n"
	// Column types.
	ColumnInt    = "int"
	ColumnString = "string"
	ColumnMoney  = "money"
	// Cell alignments.
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// Table can be rendered in a terminal.
type Table struct {
	allowEmpty   bool
	columnsTypes []string
	alignments   []Alignment
	rows         [][]string
	once         sync.Once
}

// Create a new empty table with the specified number of columns.
func New(columns int, allowEmpty bool) *Table {
	return &Table{
		allowEmpty:   allowEmpty,
		columnsTypes: make([]string, columns),
		alignments:   make([]Alignment, columns),
		rows:         make([][]string, 0, columns),
	}
}

// SetColumn set header title and column alignment settings. Column indices are 0-based.
func (t *Table) SetColumn(col int, columnType string, alignment Alignment) error {
	if col < 0 || col > len(t.columnsTypes)-1 {
		return fmt.Errorf("index out of range while attempting to set table header for column %d", col)
	}

	t.columnsTypes[col] = columnType
	t.alignments[col] = alignment

	return nil
}

// Append adds one or more rows to the table.
func (t *Table) Append(rows ...[]string) error {
	colCount := len(t.columnsTypes)

	for rowIndex, row := range rows {
		lenRaw := len(row)
		if lenRaw != colCount {
			return fmt.Errorf(
				"inconsistent number of colums for row %d; expected %d but got %d", rowIndex, colCount, lenRaw,
			)
		}
		// to prevent rewriting source slice
		cpRaw := make([]string, lenRaw)
		for valueIndex, val := range row {
			parsed, err := t.parse(t.columnsTypes[valueIndex], val)
			if err != nil {
				return fmt.Errorf(
					"row #%d: value #%d: expected type : %s got %s",
					rowIndex+1, valueIndex+1, t.columnsTypes[valueIndex], row[valueIndex],
				)
			}
			cpRaw[valueIndex] = parsed
		}
		t.rows = append(t.rows, cpRaw)
	}

	return nil
}

// Write renders table to any io.Writer.
func (t *Table) Write(to io.Writer) {
	var w bytes.Buffer

	// Calculate col widths,raw heights
	colWidths, rawHeights := t.cellSizes()

	hLine := t.hLine(colWidths)
	w.Write(hLine)

	// Render rows
	for rowIndex, row := range t.rows {
		for i := 0; i < rawHeights[rowIndex]; i++ {
			w.WriteByte(Column)
			for colIndex, cell := range row {
				cellVal := Space
				switch t.columnsTypes[colIndex] {
				case ColumnString:
					if cellLen := strings.Split(cell, Space); len(cellLen) > i {
						cellVal = cellLen[i]
					}
				default:
					if i == 0 {
						cellVal = cell
					}
				}
				w.WriteString(t.align(cellVal, colIndex, colWidths[colIndex]))
				w.WriteByte(Column)
			}
			w.Write([]byte(NewLine))
		}
		w.Write(hLine)
	}

	to.Write(w.Bytes())
}

// hLine generates horizontal line.
func (t *Table) hLine(colWidths []int) []byte {
	var w bytes.Buffer

	w.WriteByte(Center)
	for _, colWidth := range colWidths {
		w.Write(bytes.Repeat([]byte{Row}, colWidth))
		w.WriteByte(Center)
	}
	w.WriteString(NewLine)
	return w.Bytes()
}

// align makes input string alignment.
func (t *Table) align(val string, colIndex, maxWidth int) string {
	vLen := utf8.RuneCountInString(val)
	switch t.alignments[colIndex] {
	case AlignLeft:
		return val + string(bytes.Repeat([]byte(Space), maxWidth-vLen))
	case AlignRight:
		return string(bytes.Repeat([]byte(Space), maxWidth-vLen)) + val
	default:
		buffer := &bytes.Buffer{}
		lPad := (maxWidth - vLen) / 2
		buffer.Write(bytes.Repeat([]byte(Space), lPad))
		buffer.WriteString(val)
		buffer.Write(bytes.Repeat([]byte(Space), maxWidth-lPad-vLen))
		return buffer.String()
	}
}

// cellSizes calculates max width/height for each cell.
func (t *Table) cellSizes() ([]int, []int) {
	colWidths := make([]int, len(t.columnsTypes))
	rawHeights := make([]int, len(t.rows))
	for rowIndex, row := range t.rows {
		maxHeight := 1
		for columnIndex, v := range row {
			rawHeights[rowIndex] = maxHeight
			if t.columnsTypes[columnIndex] == ColumnString {
				cLen := strings.Split(v, Space)
				if len(cLen) > maxHeight {
					maxHeight = len(cLen)
				}
				for _, rw := range cLen {
					cWidth := utf8.RuneCountInString(rw)
					if cWidth > colWidths[columnIndex] {
						colWidths[columnIndex] = cWidth
					}
				}
				continue
			}
			cellWidth := utf8.RuneCountInString(v)
			if cellWidth > colWidths[columnIndex] {
				colWidths[columnIndex] = cellWidth
			}
		}
	}

	return colWidths, rawHeights
}

// parse maps cell value to type.
func (t *Table) parse(cellType string, cellValue string) (string, error) {
	if t.allowEmpty && cellValue == "" {
		return "", nil
	}
	switch cellType {
	case ColumnInt:
		_, err := strconv.Atoi(cellValue)
		if err != nil {
			return "", err
		}
		return cellValue, nil
	case ColumnString:
		return cellValue, nil
	case ColumnMoney:
		m, err := strconv.ParseFloat(cellValue, 64)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		return string(CurrencyFormat(m, 2)), nil
	default:
		return cellValue, nil
	}
}
