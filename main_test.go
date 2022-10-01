package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result, err := parseContent(input, inputFile, "")
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
}

func TestRun(t *testing.T) {
	var mockStdout bytes.Buffer

	if err := run(inputFile, "", &mockStdout, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdout.String())
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
