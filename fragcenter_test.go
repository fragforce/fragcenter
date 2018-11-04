package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"
)

var parseTests = []string{
	"resources/test_stats_payloads/multi_stream.xml",
	"resources/test_stats_payloads/single_stream.xml",
}

// TestParsing ensures that all of the example payloads provided parse into the struct we defined
func TestParsing(t *testing.T) {
	for _, fileName := range parseTests {
		f, err := os.Open(fileName)
		if err != nil {
			t.Errorf("Couldn't open the file '%s'.", fileName)
		}

		xmlBytes, err := ioutil.ReadAll(f)
		if err != nil {
			t.Errorf("Couldn't read the bytes in the file '%s'.", fileName)
		}

		var streams LiveStreams
		err = xml.Unmarshal(xmlBytes, &streams)
		if err != nil {
			t.Errorf("Couldn't unmarshal the XML in file '%s'.", fileName)
		}
	}
}
