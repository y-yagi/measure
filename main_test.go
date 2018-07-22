package main

import (
	"bytes"
	"testing"
)

func TestMeasure(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)

	measure("testdata/test.zip", out, err)
	expected := "testdata/test.zip: 8.5K\n"
	if out.String() != expected {
		t.Errorf("Expect out is %q, but %q", out.String(), expected)
	}

	out.Reset()
	measure("https://github.com/y-yagi/rails_api_sample/archive/master.zip", out, err)
	expected = "https://github.com/y-yagi/rails_api_sample/archive/master.zip: 28.7K\n"
	if out.String() != expected {
		t.Errorf("Expect out is %q, but %q", out.String(), expected)
	}
}
