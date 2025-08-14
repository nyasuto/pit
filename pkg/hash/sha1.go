package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

// SHA1 represents a 20-byte SHA-1 digest in binary form.
// Keep it as a fixed-size array to be binary-safe and allocation-friendly.
type SHA1 [20]byte

// String returns the lowercase hex representation (40 chars).
func (h SHA1) String() string {
	buf := make([]byte, 40)
	hex.Encode(buf, h[:])
	return string(buf)
}

// Bytes returns a copy of the digest bytes.
func (h SHA1) Bytes() []byte {
	b := make([]byte, len(h))
	copy(b, h[:])
	return b
}

// IsZero reports whether the digest is the zero value.
func (h SHA1) IsZero() bool {
	var z SHA1
	return h == z
}

// Short returns the first n hex characters (like Git's short hash).
// If n is out of range, it clamps to [1, 40].
func (h SHA1) Short(n int) string {
	if n < 1 {
		n = 1
	}
	if n > 40 {
		n = 40
	}
	s := h.String()
	return s[:n]
}

// Parse converts a 40-hex-character string into a SHA1.
func Parse(s string) (SHA1, error) {
	if len(s) != 40 {
		return SHA1{}, errors.New("sha1: invalid length")
	}
	var h SHA1
	b, err := hex.DecodeString(s)
	if err != nil {
		return SHA1{}, err
	}
	copy(h[:], b)
	return h, nil
}

// FromBytes constructs a SHA1 from a 20-byte slice.
func FromBytes(b []byte) (SHA1, error) {
	if len(b) != 20 {
		return SHA1{}, errors.New("sha1: invalid byte length")
	}
	var h SHA1
	copy(h[:], b)
	return h, nil
}

// Hash computes the SHA-1 of the provided data (non-streaming).
// Note: For empty input, this correctly returns the digest of empty string:
// da39a3ee5e6b4b0d3255bfef95601890afd80709
func Hash(data []byte) SHA1 {
	sum := sha1.Sum(data)
	return SHA1(sum)
}

// Hasher provides an incremental API compatible with io.Writer.
type Hasher struct {
	h hash.Hash
}

// New returns a new streaming SHA-1 hasher.
func New() *Hasher {
	return &Hasher{h: sha1.New()}
}

// Write feeds data to the hasher.
func (w *Hasher) Write(p []byte) (int, error) {
	return w.h.Write(p)
}

// Sum finalizes and returns the SHA-1 digest.
func (w *Hasher) Sum() SHA1 {
	d := w.h.Sum(nil)
	var h SHA1
	copy(h[:], d)
	return h
}

// SumReader hashes all bytes read from r.
func SumReader(r io.Reader) (SHA1, error) {
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return SHA1{}, err
	}
	var out SHA1
	copy(out[:], h.Sum(nil))
	return out, nil
}

// MarshalText implements encoding.TextMarshaler for pretty JSON/TOML/YAML.
func (h SHA1) MarshalText() ([]byte, error) {
	buf := make([]byte, 40)
	hex.Encode(buf, h[:])
	return buf, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (h *SHA1) UnmarshalText(text []byte) error {
	if len(text) != 40 {
		return errors.New("sha1: invalid text length")
	}
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	copy(h[:], b)
	return nil
}
