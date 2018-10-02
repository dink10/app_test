package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/dink10/app_test/cmd"
	"github.com/sirupsen/logrus"
)

var (
	version    = "1.0.0"
	fileName   = flag.String("f", "", "input csv file")
	allowEmpty = flag.Bool("e", false, "allow empty values in the input csv file")
	anyHeader  = flag.Bool("d", false, "allow any header type as is")
	cpus       = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
)

var usage = `Usage: converter [options...]

Options:
  -f      	Input csv file.
  -e    	Allow empty values in the input csv file (default: false)
  -d    	Allow any header type as is (default: false)
  -cpus     Number of used cpu cores (default for current machine is %d cores)
`

func main() {
	fmt.Printf("Converter version: %s\n\n", version)
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}
	flag.Parse()
	runtime.GOMAXPROCS(*cpus)
	if *fileName == "" {
		logrus.Fatal("Empty input file")
	}
	s := cmd.NewConverter(*fileName, *allowEmpty, *anyHeader)
	s.Run()
}
