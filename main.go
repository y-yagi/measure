package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/fatih/color"
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
	var retryCount int

	exitCode = 0

	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.SetOutput(errStream)
	flags.BoolVar(&showVersion, "v", false, "show version")
	flags.IntVar(&retryCount, "r", 3, "retry count")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintf(errStream, "%v\n", err)
		exitCode = 1
		return
	}

	if showVersion {
		fmt.Fprintf(outStream, "%s version: %s\n", cmd, version)
		return
	}

	if len(flags.Args()) != 1 {
		exitCode = 2
		usage(errStream)
		return
	}

	exitCode = measure(flags.Args()[0], retryCount, outStream, errStream)

	return
}

func usage(errStream io.Writer) {
	fmt.Fprintf(errStream, "Usage: %s [OPTIONS] LOCATION\n", cmd)
}

func measure(location string, retryCount int, outStream, errStream io.Writer) int {
	debugLogger = debuglog.New(outStream)
	if strings.HasPrefix(location, "http") {
		return measureURL(location, retryCount, outStream, errStream)
	}

	return measureFileOrDir(location, outStream, errStream)
}

func measureFileOrDir(location string, outStream, errStream io.Writer) int {
	fileInfo, err := os.Stat(location)
	if err != nil {
		fmt.Fprintf(errStream, "%v\n", err)
		return 1
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(fileInfo.Name())
		if err != nil {
			fmt.Fprintf(errStream, "%v\n", err)
			return 1
		}

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				fmt.Fprintf(errStream, "%v\n", err)
			} else {
				fmt.Fprintf(outStream, "%s: %s\n", file.Name(), decoratedSize(uint64(info.Size())))
			}
		}
	} else {
		fmt.Fprintf(outStream, "%s: %s\n", location, decoratedSize(uint64(fileInfo.Size())))
	}
	return 0
}

func measureURL(location string, retryCount int, outStream, errStream io.Writer) int {
	var resp *http.Response
	var err error
	lastLocation := location

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		debugLogger.Printf("Redirectd to %s\n", req.URL)
		lastLocation = req.URL.String()
		return nil
	}

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

	if isSuccess(resp) {
		if resp.ContentLength >= 0 {
			fmt.Fprintf(outStream, "%s: %s\n", location, decoratedSize(uint64(resp.ContentLength)))
		} else {
			fmt.Fprintf(outStream, "Can not get Content-Length from %s\n", location)
		}
	} else {
		fmt.Fprintf(outStream, "Error: %s\n", resp.Status)
		return 1
	}

	return 0
}

func isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func decoratedSize(size uint64) string {
	return color.GreenString(bytefmt.ByteSize(size))
}
