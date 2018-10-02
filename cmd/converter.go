package cmd

import (
	"github.com/dink10/app_test/reader"
	"github.com/sirupsen/logrus"
)

// Converter describes Converter struct.
type Converter struct {
	allowEmpty bool
	anyHeader  bool
	fileName   string
}

// NewConverter return new Converter instance.
func NewConverter(fileName string, allowEmpty bool, anyHeader bool) *Converter {
	return &Converter{fileName: fileName, allowEmpty: allowEmpty, anyHeader: anyHeader}
}

// Run starts converter job.
func (s *Converter) Run() {
	r := reader.NewFileReader(s.fileName, s.allowEmpty, s.anyHeader)
	if err := r.Validate(); err != nil {
		logrus.Fatal(err)
	}

	if err := r.Read(); err != nil {
		logrus.Fatalf("failed to open %v. %v", s.fileName, err)
	}
}
