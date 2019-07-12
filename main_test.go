package main

import (
	"bytes"
	"testing"
)

func TestMeasure(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)

	measure("testdata/test.zip", out, err)
	expected := "testdata/test.zip: \x1b[32m8.5K\x1b[0m\n"
	if out.String() != expected {
		t.Errorf("Expect out is %q, but %q", out.String(), expected)
	}

	out.Reset()
	measure("https://raw.githubusercontent.com/y-yagi/rails_api_sample/master/Gemfile", out, err)
	expected = "https://raw.githubusercontent.com/y-yagi/rails_api_sample/master/Gemfile: \x1b[32m958B\x1b[0m\n"
	if out.String() != expected {
		t.Errorf("Expect out is %q, but %q", out.String(), expected)
	}

	out.Reset()
	measure("testdata/", out, err)
	expected = "test.txt: \x1b[32m9B\x1b[0m\ntest.zip: \x1b[32m8.5K\x1b[0m\n"
	if out.String() != expected {
		t.Errorf("Expect out is %q, but %q", out.String(), expected)
	}
}
