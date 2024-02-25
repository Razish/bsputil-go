package util

import (
	"bytes"
	"unsafe"
)

func Int32AsString(i int32) string {
	bytes := make([]byte, 4)
	bytes[0] = (byte)(i & 0xff)
	bytes[1] = (byte)((i >> 8) & 0xff)
	bytes[2] = (byte)((i >> 16) & 0xff)
	bytes[3] = (byte)((i >> 24) & 0xff)

	return *(*string)(unsafe.Pointer(&bytes))
}

func CToGoString(b []byte) string {
	i := bytes.IndexByte(b, 0)
	if i < 0 {
		i = len(b)
	}
	return string(b[:i])
}
