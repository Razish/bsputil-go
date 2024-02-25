package parse

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"unsafe"

	"bsputil/util"
)

const MaxQPath = 64

type Flag uint32

type Shader struct {
	Shader       string `json:"shader"`
	SurfaceFlags Flag   `json:"surfaceFlags"`
	ContentFlags Flag   `json:"contentFlags"`
}

type Shaders []Shader

func (shaders Shaders) Write(enc *json.Encoder) error {
	for _, shader := range shaders {
		if err := enc.Encode(shader); err != nil {
			return fmt.Errorf("couldn't encode shader: %w", err)
		}
	}

	return nil
}

func loadShadersString(r *bytes.Reader, lump *Lump) (*Shaders, error) {
	if _, err := r.Seek(int64(lump.FileOffset), io.SeekStart); err != nil {
		return nil, fmt.Errorf("unable to locate shaders lump: %w", err)
	}
	data := make([]byte, lump.FileLength) // nozero: we are reading directly into the prealloced slice
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		return nil, fmt.Errorf("unable to read data from shaders lump: %w", err)
	}

	type RawShader struct {
		Shader       [MaxQPath]byte
		SurfaceFlags Flag
		ContentFlags Flag
	}

	rawShaderSize := unsafe.Sizeof(RawShader{[MaxQPath]byte{}, 0, 0})
	numShaders := uintptr(lump.FileLength) / rawShaderSize
	shaders := make(Shaders, numShaders) // nozero: we are reading directly into the prealloced slice
	for i := uintptr(0); i < numShaders; i++ {
		tmp := *(*RawShader)(unsafe.Pointer(&data[i*rawShaderSize]))
		shaders[i] = Shader{util.CToGoString(tmp.Shader[:]), tmp.SurfaceFlags, tmp.ContentFlags}
	}

	return &shaders, nil
}

func ReadShadersLump(r *bytes.Reader, header *Header) (*Shaders, error) {
	return loadShadersString(r, &header.Lumps[LumpShaders])
}
