//go:build amd64 && !tinygo
// +build amd64,!tinygo

package chaskey

//go:generate go run asm.go -out core_amd64.s
//go:noescape

func chaskeyCore(h *H, m []byte, tag []byte)
