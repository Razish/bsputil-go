package parse

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"bsputil/util"
)

//TODO: support more BSP versions?
// lump indexes to be determined by BSP ident+version

const (
	BSPIdent   = (('P' << 24) + ('S' << 16) + ('B' << 8) + 'R') // little endian {'R', 'B', 'S', 'P'}
	BSPVersion = 1
)

type LumpIndex uint32

const (
	LumpEntities LumpIndex = iota
	LumpShaders
	LumpPlanes
	LumpNodes
	LumpLeafs
	LumpLeafSurfaces
	LumpLeafBrushes
	LumpModels
	LumpBrushes
	LumpBrushSides
	LumpDrawVerts
	LumpDrawIndexes
	LumpFogs
	LumpSurfaces
	LumpLightmaps
	LumpLightGrid
	LumpVisibility
	LumpLightArray
	NumLumps
)

type Header struct {
	Ident   int32
	Version int32

	Lumps [NumLumps]Lump
}

type Lump struct {
	FileOffset uint32
	FileLength uint32
}

var ErrHeaderParse = errors.New("header parse error")

func ReadHeader(r io.Reader) (*Header, error) {
	var header Header

	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, fmt.Errorf("unable to read header: %w", err)
	}

	if header.Ident != BSPIdent {
		return nil, fmt.Errorf("%w: invalid ident (expected \"%s\", got \"%s\"", ErrHeaderParse, util.Int32AsString(BSPIdent), util.Int32AsString(header.Ident))
	}
	if header.Version != BSPVersion {
		return nil, fmt.Errorf("%w: invalid version (expected \"%d\", got \"%d\"", ErrHeaderParse, BSPVersion, header.Version)
	}

	return &header, nil
}
