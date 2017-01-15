// +build amd64

package chaskey

//go:generate python -m peachpy.x86_64 core.py -S -o core_amd64.s -mabi=goasm
//go:noescape

func chaskeyCore(h *H, m []byte, tag []byte)
