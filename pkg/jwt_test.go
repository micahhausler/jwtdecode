package pkg

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestDecode(t *testing.T) {

	cases := []struct {
		name                string
		tokenFiles          []string
		expectedError       error
		expectedContentFile string
	}{
		{
			"Valid single b64url token",
			[]string{"./test_fixtures/valid-b64url.jwt"},
			nil,
			"./test_fixtures/valid-b64url.jwt-output",
		},
		{
			"Valid single b64 token",
			[]string{"./test_fixtures/valid-b64.jwt"},
			nil,
			"./test_fixtures/valid-b64.jwt-output",
		},
		{
			"Valid multiple token",
			[]string{"./test_fixtures/valid-multi-token.jwt"},
			nil,
			"./test_fixtures/valid-multi-token.jwt-output",
		},
		{
			"Invalid single token",
			[]string{"./test_fixtures/invalid-single-token.jwt"},
			fmt.Errorf("Expected 3 parts to the JWT, got 2"),
			"./test_fixtures/invalid-single-token.jwt-output",
		},
		{
			"Invalid b64",
			[]string{"./test_fixtures/not-a-token.jwt"},
			fmt.Errorf("Could not decode input"),
			"./test_fixtures/not-a-token.jwt-output",
		},
		{
			"Invalid json",
			[]string{"./test_fixtures/invalid-json.jwt"},
			fmt.Errorf("Could not unmarshal payload: invalid character '}' looking for beginning of value"),
			"./test_fixtures/invalid-json.jwt-output",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			var output io.Writer
			if *update {
				f, err := os.OpenFile(c.expectedContentFile, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					t.Errorf("Error opening %s: %v", c.expectedContentFile, err)
					return
				}
				defer f.Close()
				output = f
			} else {
				output = buf
			}

			err := DecodeFiles(output, c.tokenFiles)
			if c.expectedError == nil && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			} else if c.expectedError != nil && err == nil {
				t.Errorf("Unexpected success, expected error: %v", c.expectedError)
				return
			}
			if c.expectedError != nil && err != nil {
				if c.expectedError.Error() != err.Error() {
					t.Errorf("Unexpected error. Got '%v', expected '%v'", err, c.expectedError)
					return
				}
				return
			}

			if *update {
				return
			}

			expectedContent, err := ioutil.ReadFile(c.expectedContentFile)
			if err != nil {
				t.Errorf("Unexpected error reading file: %v", err)
				return
			}

			if !bytes.Equal(buf.Bytes(), expectedContent) {
				t.Errorf("Unexpected content. Expected %s, got %s", string(buf.Bytes()), string(expectedContent))
				return
			}
		})
	}
}
