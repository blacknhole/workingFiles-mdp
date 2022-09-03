package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "test1.md.html"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result := parseContent(input)

	exp, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, exp) {
		t.Logf("Result:\n%s\n", result)
		t.Logf("Expected:\n%s\n", exp)
		t.Error("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	if err := run(inputFile); err != nil {
		t.Fatal(err)
	}

	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	exp, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, exp) {
		t.Logf("Result:\n%s\n", result)
		t.Logf("Expected:\n%s\n", exp)
		t.Error("Result content does not match golden file")
	}

	os.Remove(resultFile)
}
