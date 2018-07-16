package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"code.cloudfoundry.org/bytefmt"
)

const (
	version = "0.1.0"
	cmd     = "measure"
)

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
	res, err := http.Head(location)
	if err != nil {
		fmt.Fprintf(errStream, "%v\n", err)
		return 1
	}

	fmt.Fprintf(outStream, "%s: %s\n", location, bytefmt.ByteSize(uint64(res.ContentLength)))
	return 0
}
