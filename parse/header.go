package parse

import (
	"bsputil/util"
	"bytes"
	"encoding/binary"
	"fmt"
)

//TODO: support more BSP versions?
// lump indexes to be determined by BSP ident+version

const (
	BSP_IDENT   = (('P' << 24) + ('S' << 16) + ('B' << 8) + 'R') // big endian {'R', 'B', 'S', 'P'}
	BSP_VERSION = 1
)

const (
	LUMP_ENTITIES     = 0
	LUMP_SHADERS      = 1
	LUMP_PLANES       = 2
	LUMP_NODES        = 3
	LUMP_LEAFS        = 4
	LUMP_LEAFSURFACES = 5
	LUMP_LEAFBRUSHES  = 6
	LUMP_MODELS       = 7
	LUMP_BRUSHES      = 8
	LUMP_BRUSHSIDES   = 9
	LUMP_DRAWVERTS    = 10
	LUMP_DRAWINDEXES  = 11
	LUMP_FOGS         = 12
	LUMP_SURFACES     = 13
	LUMP_LIGHTMAPS    = 14
	LUMP_LIGHTGRID    = 15
	LUMP_VISIBILITY   = 16
	LUMP_LIGHTARRAY   = 17
	HEADER_LUMPS      = 18
)

type Header struct {
	Ident   int32
	Version int32

	Lumps [HEADER_LUMPS]Lump
}

type Lump struct {
	FileOffset uint32
	FileLength uint32
}

func ReadHeader(r *bytes.Reader) (*Header, error) {
	var header Header

	binary.Read(r, binary.LittleEndian, &header)

	if header.Ident != BSP_IDENT {
		return nil, fmt.Errorf("invalid header ident (expected \"%s\", got \"%s\"", util.Int32AsString(BSP_IDENT), util.Int32AsString(header.Ident))
	}
	if header.Version != BSP_VERSION {
		return nil, fmt.Errorf("invalid header version (expected \"%d\", got \"%d\"", BSP_VERSION, header.Version)
	}

	return &header, nil
}
