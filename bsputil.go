package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"bsputil/parse"
)

var ErrUsage = errors.New("usage error")

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		prog := os.Args[0]
		log.Fatal(fmt.Errorf("%w: %s <filename> <lump name>", ErrUsage, prog))
	}

	enc := json.NewEncoder(os.Stdout)

	filename := args[0]
	lump := args[1]
	if err := readBSP(filename, lump, enc); err != nil {
		log.Fatal(err)
	}
}

func readBSP(filename string, lump string, enc *json.Encoder) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("unable to read BSP: %w", err)
	}

	r := bytes.NewReader(data)
	header, err := parse.ReadHeader(r)
	if err != nil {
		return fmt.Errorf("unable to parse BSP header: %w", err)
	}

	// each lump should only emit JSON-L for external consumption
	switch lump {
	case "ents", "entities":
		entities, err := parse.ReadEntityLump(r, header)
		if err != nil {
			return fmt.Errorf("unable to parse entity lump: %w", err)
		}
		if err := entities.Write(enc); err != nil {
			return fmt.Errorf("unable to encode entities: %w", err)
		}

	case "shaders":
		shaders, err := parse.ReadShadersLump(r, header)
		if err != nil {
			return fmt.Errorf("unable to parse shaders lump: %w", err)
		}
		if err := shaders.Write(enc); err != nil {
			return fmt.Errorf("unable to encode shaders: %w", err)
		}

	default:
		return fmt.Errorf("%w: unknown lump %s", ErrUsage, lump)
	}

	return nil
}
