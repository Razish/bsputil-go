package main

import (
	"bsputil/parse"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		prog := os.Args[0]
		log.Fatal(fmt.Errorf("usage: %s <filename> <lump name>", prog))
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
		return fmt.Errorf("unable to read BSP: %v", err)
	}

	r := bytes.NewReader(data)
	header, err := parse.ReadHeader(r)
	if err != nil {
		return fmt.Errorf("unable to parse BSP header: %v", err)
	}

	// each lump should only emit JSON-L for external consumption
	switch lump {
	case "ents", "entities":
		entities, err := parse.ReadEntityLump(r, header)
		if err != nil {
			return fmt.Errorf("unable to parse entity lump: %v", err)
		}
		entities.Write(enc)

	case "shaders":
		shaders, err := parse.ReadShadersLump(r, header)
		if err != nil {
			return fmt.Errorf("unable to parse shaders lump: %v", err)
		}
		shaders.Write(enc)
	}

	return nil
}
