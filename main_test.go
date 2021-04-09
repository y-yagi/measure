package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMeasure(t *testing.T) {
	retryCount := 1
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	expected := ""

	measure("testdata/test.zip", retryCount, out, err)
	if len(os.Getenv("CI")) != 0 {
		expected = "testdata/test.zip: 8.5K\n"
	} else {
		expected = "testdata/test.zip: \x1b[32m8.5K\x1b[0m\n"
	}
	if out.String() != expected {
		t.Fatalf("Expected '%q', but got '%q'", expected, out.String())
	}

	out.Reset()
	measure("https://raw.githubusercontent.com/y-yagi/rails_api_sample/master/Gemfile", retryCount, out, err)
	if len(os.Getenv("CI")) != 0 {
		expected = "https://raw.githubusercontent.com/y-yagi/rails_api_sample/master/Gemfile: 958B\n"
	} else {
		expected = "https://raw.githubusercontent.com/y-yagi/rails_api_sample/master/Gemfile: \x1b[32m958B\x1b[0m\n"
	}
	if out.String() != expected {
		t.Fatalf("Expected '%q', but got '%q'", expected, out.String())
	}

	out.Reset()
	measure("testdata/", retryCount, out, err)
	if len(os.Getenv("CI")) != 0 {
		expected = "test.txt: 9B\ntest.zip: 8.5K\n"
	} else {
		expected = "test.txt: \x1b[32m9B\x1b[0m\ntest.zip: \x1b[32m8.5K\x1b[0m\n"
	}
	if out.String() != expected {
		t.Fatalf("Expected '%q', but got '%q'", expected, out.String())
	}
}
