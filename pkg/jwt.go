package pkg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
)

// DecodeFiles reads a list of filenames and writes the decoded tokens to
// the writer
func DecodeFiles(out io.Writer, files []string) error {
	decoder := TokenDecoder{}
	for _, filename := range files {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		decodedTokens, err := decoder.DecodeLines(content)
		if err != nil {
			return err
		}
		for _, token := range decodedTokens {
			for _, line := range token {
				fmt.Fprintln(out, string(line))
			}
		}
	}
	return nil
}

// DecodeLines splits a file on newlines and decodes each line
func (d *TokenDecoder) DecodeLines(input []byte) ([][][]byte, error) {
	tokenLines := [][][]byte{}
	for _, line := range bytes.Split(input, []byte("\n")) {
		if len(line) < 2 {
			continue
		}
		header, body, err := d.Decode(line)
		if err != nil {
			return nil, err
		}
		headerContent, err := marshalIndent(header)
		if err != nil {
			return nil, err
		}
		bodyContent, err := marshalIndent(body)
		if err != nil {
			return nil, err
		}
		tokenLines = append(tokenLines, [][]byte{headerContent, bodyContent})
	}
	return tokenLines, nil
}

func marshalIndent(input []byte) ([]byte, error) {
	data := &map[string]interface{}{}
	err := json.Unmarshal([]byte(input), data)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not unmarshal payload")
	}
	return json.MarshalIndent(data, "", "    ")
}

// Decode returns the header, body, and an error of the included token.
// It attempts to use both the raw and regular URL and standard Base64
// decoding
func (d *TokenDecoder) Decode(input []byte) (header []byte, body []byte, err error) {
	parts := bytes.Split(input, []byte("."))
	if len(parts) != 3 {
		return nil, nil, errors.Errorf("Expected 3 parts to the JWT, got %d", len(parts))
	}
	header, err = d.b64decode(parts[0])
	if err != nil {
		return nil, nil, err
	}
	body, err = d.b64decode(parts[1])
	if err != nil {
		return nil, nil, err
	}
	return header, body, nil
}

// TokenDecoder decodes tokens
type TokenDecoder struct{}

func (d *TokenDecoder) b64decode(data []byte) ([]byte, error) {
	var encodings = []*base64.Encoding{
		base64.URLEncoding,
		base64.StdEncoding,
		base64.RawURLEncoding,
		base64.RawStdEncoding,
	}
	for i, enc := range encodings {
		dst := make([]byte, enc.DecodedLen(len(data)))
		l, err := enc.Decode(dst, data)
		if err != nil && i == len(encodings)-1 {
			return nil, errors.Wrap(err, "No encodings passed")
		}
		if len(dst) > 0 {
			if dst[l-1] != byte('}') {
				continue
			}
			return dst[:l], nil
		}
	}
	return nil, errors.New("Could not decode input")
}
