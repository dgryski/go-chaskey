// +build !amd64

package chaskey

import (
	"encoding/binary"
	"math/bits"
)

func chaskeyCore(h *H, m []byte, tag []byte) {

	v := h.k

	for ; len(m) > 16; m = m[16:] {

		v[0] ^= binary.LittleEndian.Uint32(m[0:4])
		v[1] ^= binary.LittleEndian.Uint32(m[4:8])
		v[2] ^= binary.LittleEndian.Uint32(m[8:12])
		v[3] ^= binary.LittleEndian.Uint32(m[12:16])

		// permute
		for i := 0; i < h.r; i++ {
			// round
			v[0] += v[1]
			v[1] = bits.RotateLeft32(v[1], 5)
			v[1] ^= v[0]
			v[0] = bits.RotateLeft32(v[0], 16)
			v[2] += v[3]
			v[3] = bits.RotateLeft32(v[3], 8)
			v[3] ^= v[2]
			v[0] += v[3]
			v[3] = bits.RotateLeft32(v[3], 13)
			v[3] ^= v[0]
			v[2] += v[1]
			v[1] = bits.RotateLeft32(v[1], 7)
			v[1] ^= v[2]
			v[2] = bits.RotateLeft32(v[2], 16)
		}
	}

	var l [4]uint32
	var lastblock [4]uint32

	if len(m) == 16 {
		l = h.k1

		lastblock[0] = binary.LittleEndian.Uint32(m[0:4])
		lastblock[1] = binary.LittleEndian.Uint32(m[4:8])
		lastblock[2] = binary.LittleEndian.Uint32(m[8:12])
		lastblock[3] = binary.LittleEndian.Uint32(m[12:16])

	} else {
		l = h.k2
		var lb [16]byte
		copy(lb[:], m)

		lb[len(m)] = 0x01

		lastblock[0] = binary.LittleEndian.Uint32(lb[0:4])
		lastblock[1] = binary.LittleEndian.Uint32(lb[4:8])
		lastblock[2] = binary.LittleEndian.Uint32(lb[8:12])
		lastblock[3] = binary.LittleEndian.Uint32(lb[12:16])
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
		v[1] = bits.RotateLeft32(v[1], 5)
		v[1] ^= v[0]
		v[0] = bits.RotateLeft32(v[0], 16)
		v[2] += v[3]
		v[3] = bits.RotateLeft32(v[3], 8)
		v[3] ^= v[2]
		v[0] += v[3]
		v[3] = bits.RotateLeft32(v[3], 13)
		v[3] ^= v[0]
		v[2] += v[1]
		v[1] = bits.RotateLeft32(v[1], 7)
		v[1] ^= v[2]
		v[2] = bits.RotateLeft32(v[2], 16)
	}

	v[0] ^= l[0]
	v[1] ^= l[1]
	v[2] ^= l[2]
	v[3] ^= l[3]

	_ = tag[15]

	binary.LittleEndian.PutUint32(tag[0:4], v[0])
	binary.LittleEndian.PutUint32(tag[4:8], v[1])
	binary.LittleEndian.PutUint32(tag[8:12], v[2])
	binary.LittleEndian.PutUint32(tag[12:16], v[3])

}
