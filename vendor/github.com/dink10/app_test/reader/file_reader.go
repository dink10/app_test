package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dink10/app_test/table"
	"github.com/sirupsen/logrus"
)

const CsvExtension = ".csv"

// FileReader describes reader struct.
type FileReader struct {
	allowEmpty bool
	anyHeader  bool
	fileName   string
}

// NewFileReader return new FileReader instance.
func NewFileReader(fileName string, allowEmpty bool, anyHeader bool) *FileReader {
	return &FileReader{fileName: fileName, allowEmpty: allowEmpty, anyHeader: anyHeader}
}

// Validate checks file extension.
func (r *FileReader) Validate() error {
	if !strings.HasSuffix(r.fileName, CsvExtension) {
		return fmt.Errorf("inccorect file type")
	}
	return nil
}

// Read reads file.
func (r *FileReader) Read() error {
	file, err := os.Open(r.fileName)
	if err != nil {
		return err
	}
	defer func() {
		if errCl := file.Close(); errCl != nil {
			logrus.Error(errCl)
		}
	}()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true // Allow lazy quotes.

	record, err := reader.Read()
	if err == io.EOF {
		return fmt.Errorf("empty file")
	}

	headerLength := len(record)
	if headerLength == 0 {
		return fmt.Errorf("empty header")
	}
	t := table.New(headerLength, r.allowEmpty)

	for k, v := range record {
		switch v {
		case table.ColumnInt:
			t.SetColumn(k, table.ColumnInt, table.AlignRight)
		case table.ColumnString:
			t.SetColumn(k, table.ColumnString, table.AlignLeft)
		case table.ColumnMoney:
			t.SetColumn(k, table.ColumnMoney, table.AlignRight)
		default:
			if !r.anyHeader {
				return fmt.Errorf("wrong header type %s in position %d", v, k)
			}
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err := t.Append(record); err != nil {
			return err
		}
	}

	t.Write(os.Stdout)

	return nil
}
