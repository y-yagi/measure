package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/y-yagi/debuglog"
)

const (
	version = "0.1.0"
	cmd     = "measure"
)

var debugLogger *debuglog.Logger

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func run(args []string, outStream, errStream io.Writer) (exitCode int) {

	var showVersion bool

	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.SetOutput(errStream)
	flags.BoolVar(&showVersion, "v", false, "show version")
	flags.Parse(args[1:])

	exitCode = 0

	if showVersion {
		fmt.Fprintf(outStream, "%s version: %s\n", cmd, version)
		return
	}

	if len(flags.Args()) != 1 {
		exitCode = 2
		usage(errStream)
		return
	}

	exitCode = measure(flags.Args()[0], outStream, errStream)

	return
}

func usage(errStream io.Writer) {
	fmt.Fprintf(errStream, "Usage: %s [OPTIONS] LOCATION\n", cmd)
}

func measure(location string, outStream, errStream io.Writer) int {
	debugLogger = debuglog.New(outStream)
	if strings.HasPrefix(location, "http") {
		return measureURL(location, outStream, errStream)
	}

	return measureFile(location, outStream, errStream)
}

func measureFile(location string, outStream, errStream io.Writer) int {
	fileInfo, err := os.Stat(location)
	if err != nil {
		fmt.Fprintf(errStream, "%v\n", err)
		return 1
	}

	fmt.Fprintf(outStream, "%s: %s\n", location, bytefmt.ByteSize(uint64(fileInfo.Size())))
	return 0
}

func measureURL(location string, outStream, errStream io.Writer) int {
	var resp *http.Response
	var err error
	lastLocation := location

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		debugLogger.Printf("Redirectd to %s\n", req.URL)
		lastLocation = req.URL.String()
		return nil
	}

	retryCount := 3 // TODO: Can specify a value from arguments.

	for i := 0; i < retryCount; i++ {
		resp, err = client.Head(lastLocation)
		if err != nil {
			fmt.Fprintf(errStream, "%v\n", err)
			return 1
		}

		if resp.ContentLength >= 0 {
			break
		}
		debugLogger.Printf("Retry(count: %d, URL: %s)\n", i+1, lastLocation)
	}

	if resp.ContentLength >= 0 {
		fmt.Fprintf(outStream, "%s: %s\n", location, bytefmt.ByteSize(uint64(resp.ContentLength)))
	} else {
		fmt.Fprintf(outStream, "Can not get Content-Length from %s\n", location)
	}

	return 0
}
