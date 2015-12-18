// Package chaskey implements the Chaskey MAC
/*

http://mouha.be/chaskey/

https://eprint.iacr.org/2014/386.pdf
https://eprint.iacr.org/2015/1182.pdf

*/
package chaskey

import "encoding/binary"

// H holds keys for an instance of chaskey
type H struct {
	k  [4]uint32
	k1 [4]uint32
	k2 [4]uint32
	r  int
}

// New returns a new 8-round chaskey hasher.
func New(k [4]uint32) *H { return newH(k, 8) }

// New12 returns a new 12-round chaskey hasher.
func New12(k [4]uint32) *H { return newH(k, 12) }

func newH(k [4]uint32, rounds int) *H {

	h := H{
		k: k,
		r: rounds,
	}

	timestwo(h.k1[:], k[:])
	timestwo(h.k2[:], h.k1[:])

	return &h
}

// MAC computes the chaskey MAC of a message m.  The returned byte slice will be a subslice of tag, if provided.
func (h *H) MAC(m, tag []byte) []byte {

	v := h.k

	if len(tag) < 16 {
		tag = make([]byte, 16)
	}

	for ; len(m) > 16; m = m[16:] {

		v[0] ^= binary.LittleEndian.Uint32(m[0:])
		v[1] ^= binary.LittleEndian.Uint32(m[4:])
		v[2] ^= binary.LittleEndian.Uint32(m[8:])
		v[3] ^= binary.LittleEndian.Uint32(m[12:])

		// permute
		for i := 0; i < h.r; i++ {
			// round
			v[0] += v[1]
			v[1] = rotl32(v[1], 5)
			v[1] ^= v[0]
			v[0] = rotl32(v[0], 16)
			v[2] += v[3]
			v[3] = rotl32(v[3], 8)
			v[3] ^= v[2]
			v[0] += v[3]
			v[3] = rotl32(v[3], 13)
			v[3] ^= v[0]
			v[2] += v[1]
			v[1] = rotl32(v[1], 7)
			v[1] ^= v[2]
			v[2] = rotl32(v[2], 16)
		}
	}

	var l [4]uint32
	var lastblock [4]uint32

	if len(m) == 16 {
		l = h.k1

		lastblock[0] = binary.LittleEndian.Uint32(m[0:])
		lastblock[1] = binary.LittleEndian.Uint32(m[4:])
		lastblock[2] = binary.LittleEndian.Uint32(m[8:])
		lastblock[3] = binary.LittleEndian.Uint32(m[12:])

	} else {
		l = h.k2
		var lb [16]byte
		copy(lb[:], m)

		lb[len(m)] = 0x01

		lastblock[0] = binary.LittleEndian.Uint32(lb[0:])
		lastblock[1] = binary.LittleEndian.Uint32(lb[4:])
		lastblock[2] = binary.LittleEndian.Uint32(lb[8:])
		lastblock[3] = binary.LittleEndian.Uint32(lb[12:])
	}

	v[0] ^= lastblock[0]
	v[1] ^= lastblock[1]
	v[2] ^= lastblock[2]
	v[3] ^= lastblock[3]

	v[0] ^= l[0]
	v[1] ^= l[1]
	v[2] ^= l[2]
	v[3] ^= l[3]

	// permute
	for i := 0; i < h.r; i++ {
		// round
		v[0] += v[1]
		v[1] = rotl32(v[1], 5)
		v[1] ^= v[0]
		v[0] = rotl32(v[0], 16)
		v[2] += v[3]
		v[3] = rotl32(v[3], 8)
		v[3] ^= v[2]
		v[0] += v[3]
		v[3] = rotl32(v[3], 13)
		v[3] ^= v[0]
		v[2] += v[1]
		v[1] = rotl32(v[1], 7)
		v[1] ^= v[2]
		v[2] = rotl32(v[2], 16)
	}

	v[0] ^= l[0]
	v[1] ^= l[1]
	v[2] ^= l[2]
	v[3] ^= l[3]

	binary.LittleEndian.PutUint32(tag[0:], v[0])
	binary.LittleEndian.PutUint32(tag[4:], v[1])
	binary.LittleEndian.PutUint32(tag[8:], v[2])
	binary.LittleEndian.PutUint32(tag[12:], v[3])

	return tag[:16]
}

func rotl32(x uint32, b uint) uint32 {
	return (x >> (32 - b)) | (x << b)
}

func timestwo(out []uint32, in []uint32) {
	var C = [2]uint32{0x00, 0x87}
	out[0] = (in[0] << 1) ^ C[in[3]>>31]
	out[1] = (in[1] << 1) | (in[0] >> 31)
	out[2] = (in[2] << 1) | (in[1] >> 31)
	out[3] = (in[3] << 1) | (in[2] >> 31)
}
