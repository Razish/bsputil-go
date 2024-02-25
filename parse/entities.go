package parse

import (
	"bsputil/util"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/scanner"
)

type EntitiesString string
type Entity map[string]string
type Entities []Entity

func (entities Entities) Write(enc *json.Encoder) error {
	for _, entity := range entities {
		enc.Encode(entity)
	}
	return nil
}

type EntityScanner struct {
	scanner.Scanner
}

// token is the entire `"quoted string"`
func (es *EntityScanner) ReadQuotedString(token rune) (string, error) {
	if token == scanner.EOF {
		return "", fmt.Errorf("unexpected EOF parsing quoted string")
	}

	// return the inside portion of the `"quoted string"`
	text := es.TokenText()
	if text[len(text)-1:] != "\"" {
		return "", fmt.Errorf("quoted string did not end in a quote")
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
		text := es.TokenText()
		if text != "{" {
			return nil, fmt.Errorf("unexpected token \"%s\" reading entity, expected \"{\"", text)
		}

		entity := Entity{}

		// we are parsing an entity
		for token = es.Scan(); token != scanner.EOF; token = es.Scan() {
			text = es.TokenText()

			if text == "}" {
				entities = append(entities, entity)
				break
			}

			// we're parsing k,v pairs of entity properties until the next "}"
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
	r.Seek(int64(lump.FileOffset), io.SeekStart)
	data := make([]byte, lump.FileLength)
	binary.Read(r, binary.LittleEndian, data)

	entitiesString := EntitiesString(util.CToGoString(data[:]))
	return entitiesString.Parse()
}

func ReadEntityLump(r *bytes.Reader, header *Header) (*Entities, error) {
	return loadEntitiesString(r, &header.Lumps[LUMP_ENTITIES])
}
