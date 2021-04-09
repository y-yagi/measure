package main

import (
	"bytes"
	"testing"
)

func TestMeasure(t *testing.T) {
	retryCount := 1
	out, err := new(bytes.Buffer), new(bytes.Buffer)

	measure("testdata/test.zip", retryCount, out, err)
	expected := "testdata/test.zip: \x1b[32m8.5K\x1b[0m\n"
	if out.String() != expected {
		t.Fatalf("Expected '%q', but got '%q'", expected, out.String())
	}

	out.Reset()
	measure("https://raw.githubusercontent.com/y-yagi/rails_api_sample/master/Gemfile", retryCount, out, err)
	expected = "https://raw.githubusercontent.com/y-yagi/rails_api_sample/master/Gemfile: \x1b[32m958B\x1b[0m\n"
	if out.String() != expected {
		t.Fatalf("Expected '%q', but got '%q'", expected, out.String())
	}

	out.Reset()
	measure("testdata/", retryCount, out, err)
	expected = "test.txt: \x1b[32m9B\x1b[0m\ntest.zip: \x1b[32m8.5K\x1b[0m\n"
	if out.String() != expected {
		t.Fatalf("Expected '%q', but got '%q'", expected, out.String())
	}
}
