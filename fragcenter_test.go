package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/fragforce/fragcenter/internal/stream_checker"
)

// TestParsing ensures that all of the example payloads provided parse into the struct we defined
func TestParsing(t *testing.T) {
	files, err := filepath.Glob("assets/resources/test_stats_payloads/*xml")
	if err != nil {
		t.Errorf("No XML files found in assets/resources/test_stats_payloads.")
		t.FailNow()
	}

	for _, fileName := range files {
		f, err := os.Open(fileName)
		if err != nil {
			t.Errorf("Couldn't open the file '%s'.", fileName)
			continue
		}

		xmlBytes, err := ioutil.ReadAll(f)
		if err != nil {
			t.Errorf("Couldn't read the bytes in the file '%s'.", fileName)
			continue
		}

		var streams stream_checker.LiveStreams
		err = xml.Unmarshal(xmlBytes, &streams)
		if err != nil {
			t.Errorf("Couldn't unmarshal the XML in file '%s'.", fileName)
			continue
		}
	}
}
