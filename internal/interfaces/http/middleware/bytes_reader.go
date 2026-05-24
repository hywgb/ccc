package middleware

import "bytes"

// bytesReader is a small helper to wrap a byte slice as an io.Reader without
// pulling in unnecessary buffers, mirroring bytes.NewReader so callers do not
// need to import bytes for a single line.
func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }
