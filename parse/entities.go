package parse

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"

	"bsputil/util"
)

type EntitiesString string
type Entity map[string]string
type Entities []Entity

func (entities Entities) Write(enc *json.Encoder) error {
	for _, entity := range entities {
		if err := enc.Encode(entity); err != nil {
			return fmt.Errorf("couldn't encode entity: %w", err)
		}
	}

	return nil
}

type EntityScanner struct {
	scanner.Scanner
}

var ErrEntityParse = errors.New("entity parse error")

func EntityParseError(detail string) error {
	return fmt.Errorf("%w: %s", ErrEntityParse, detail)
}

// token is the entire `"quoted string"`.
func (es *EntityScanner) ReadQuotedString(token rune) (string, error) {
	if token == scanner.EOF {
		return "", EntityParseError("unexpected EOF parsing quoted string")
	}

	// return the inside portion of the `"quoted string"`
	text := es.TokenText()
	if !strings.HasSuffix(text, "\"") {
		return "", fmt.Errorf("%w: quoted string did not end in a quote: %s", ErrEntityParse, text)
	}

	return text[1 : len(text)-1], nil
}

func (entitiesString EntitiesString) Parse() (*Entities, error) {
	entities := make(Entities, 0, 16)

	var es EntityScanner
	es.Init(strings.NewReader(string(entitiesString)))

	// note: we don't have to support escaped quotes 💡
	// format:
	//   {
	//     "key1" "value1"
	//     "key2" "value two"
	//   }
	//   {
	//     "key1" "value"
	//     "key2" "value"
	//   }
	for token := es.Scan(); token != scanner.EOF; token = es.Scan() {
		// we're parsing an array of entity defs between braces: "{" ... "}" until EOF
		text := es.TokenText()
		if text != "{" {
			return nil, fmt.Errorf("%w: unexpected token \"%s\" reading entity, expected \"{\"", ErrEntityParse, text)
		}

		entity := Entity{}

		for token = es.Scan(); token != scanner.EOF; token = es.Scan() {
			// we're parsing k,v pairs of entity properties until the next "}"
			text = es.TokenText()

			if text == "}" {
				entities = append(entities, entity)

				break
			}

			key, err := es.ReadQuotedString(token)
			if err != nil {
				return nil, fmt.Errorf("couldn't parse entity key: %w", err)
			}

			token = es.Scan()
			value, err := es.ReadQuotedString(token)
			if err != nil {
				return nil, fmt.Errorf("couldn't parse entity value: %w", err)
			}

			entity[key] = value
		}
	}

	return &entities, nil
}

func loadEntitiesString(r *bytes.Reader, lump *Lump) (*Entities, error) {
	if _, err := r.Seek(int64(lump.FileOffset), io.SeekStart); err != nil {
		return nil, fmt.Errorf("unable to locate entities lump: %w", err)
	}
	data := make([]byte, lump.FileLength) // nozero: we are reading directly into the prealloced slice
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		return nil, fmt.Errorf("unable to read data from entities lump: %w", err)
	}

	entitiesString := EntitiesString(util.CToGoString(data[:]))

	return entitiesString.Parse()
}

func ReadEntityLump(r *bytes.Reader, header *Header) (*Entities, error) {
	return loadEntitiesString(r, &header.Lumps[LumpEntities])
}
