package parse

import (
	"bsputil/util"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"unsafe"
)

const MAX_QPATH = 64

type Flag uint32

type Shader struct {
	Shader       string `json:"shader"`
	SurfaceFlags Flag   `json:"surfaceFlags"`
	ContentFlags Flag   `json:"contentFlags"`
}

type Shaders []Shader

func (shaders Shaders) Write(enc *json.Encoder) {
	for _, shader := range shaders {
		enc.Encode(shader)
	}
}

func loadShadersString(r *bytes.Reader, lump *Lump) (*Shaders, error) {
	r.Seek(int64(lump.FileOffset), io.SeekStart)
	data := make([]byte, lump.FileLength)
	binary.Read(r, binary.LittleEndian, data)

	type RawShader struct {
		Shader       [MAX_QPATH]byte
		SurfaceFlags Flag
		ContentFlags Flag
	}

	rawShaderSize := unsafe.Sizeof(RawShader{})
	numShaders := uintptr(lump.FileLength) / rawShaderSize
	shaders := make(Shaders, numShaders)
	for i := uintptr(0); i < numShaders; i++ {
		tmp := *(*RawShader)(unsafe.Pointer(&data[i*rawShaderSize]))
		shaders[i] = Shader{util.CToGoString(tmp.Shader[:]), tmp.SurfaceFlags, tmp.ContentFlags}
	}

	return &shaders, nil
}

func ReadShadersLump(r *bytes.Reader, header *Header) (*Shaders, error) {
	return loadShadersString(r, &header.Lumps[LUMP_SHADERS])
}
