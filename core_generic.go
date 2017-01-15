// +build !amd64

package chaskey

import "encoding/binary"

func chaskeyCore(h *H, m []byte, tag []byte) {

	v := h.k

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

}
