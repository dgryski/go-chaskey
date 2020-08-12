// +build ignore

package main

import (
	"strconv"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func permute(r, v0, v1, v2, v3 Register) {
	i := GP64()
	MOVQ(r, i)
	l := newLoop()
	Label(l.begin())
	ADDL(v1, v0)
	ADDL(v3, v2)
	ROLL(Imm(5), v1)
	ROLL(Imm(8), v3)
	XORL(v0, v1)
	XORL(v2, v3)
	ROLL(Imm(16), v0)
	ADDL(v3, v0)
	ADDL(v1, v2)
	ROLL(Imm(13), v3)
	ROLL(Imm(7), v1)
	XORL(v0, v3)
	XORL(v2, v1)
	ROLL(Imm(16), v2)
	DECQ(i)
	JNZ(LabelRef(l.begin()))
}

func main() {
	Package("github.com/dgryski/go-chaskey")

	TEXT("chaskeyCore", NOSPLIT, "func(h *H, m []byte, tag []byte)")

	v0, v1, v2, v3 := GP32(), GP32(), GP32(), GP32()

	reg_h := GP64()
	Load(Param("h"), reg_h)

	MOVL(Mem{Base: reg_h, Disp: 0}, v0)
	MOVL(Mem{Base: reg_h, Disp: 4}, v1)
	MOVL(Mem{Base: reg_h, Disp: 8}, v2)
	MOVL(Mem{Base: reg_h, Disp: 12}, v3)

	reg_r := GP64()
	MOVQ(Mem{Base: reg_h, Disp: 48}, reg_r)

	reg_m := GP64()
	reg_m_len := GP64()
	Load(Param("m").Base(), reg_m)
	Load(Param("m").Len(), reg_m_len)

	CMPQ(reg_m_len, Imm(16))
	l := newLoop()
	JLE(LabelRef(l.end()))
	Label(l.begin())
	XORL(Mem{Base: reg_m, Disp: 0}, v0)
	XORL(Mem{Base: reg_m, Disp: 4}, v1)
	XORL(Mem{Base: reg_m, Disp: 8}, v2)
	XORL(Mem{Base: reg_m, Disp: 12}, v3)

	permute(reg_r, v0, v1, v2, v3)

	ADDQ(Imm(16), reg_m)
	SUBQ(Imm(16), reg_m_len)
	CMPQ(reg_m_len, Imm(16))
	JG(LabelRef(l.begin()))
	Label(l.end())

	lkey := GP64()
	MOVQ(reg_h, lkey)

	CMPQ(reg_m_len, Imm(16))
	JNZ(LabelRef("switch"))

	ADDQ(Imm(16), lkey)

	XORL(Mem{Base: reg_m, Disp: 0}, v0)
	XORL(Mem{Base: reg_m, Disp: 4}, v1)
	XORL(Mem{Base: reg_m, Disp: 8}, v2)
	XORL(Mem{Base: reg_m, Disp: 12}, v3)

	JMP(LabelRef("afterSwitch"))

	Label("switch")

	ADDQ(Imm(32), lkey)

	lb := GP64()
	Load(Param("tag").Base(), lb)

	MOVQ(U32(0), Mem{Base: lb})
	MOVQ(U32(0), Mem{Base: lb, Disp: 8})

	var labels []string
	for i := 0; i < 16; i++ {
		labels = append(labels, "sw"+strconv.Itoa(i))
	}

	for i := 0; i < 16; i++ {
		CMPQ(reg_m_len, Imm(uint64(i)))
		JE(LabelRef(labels[i]))
	}

	char := GP8()
	for i := 15; i > 0; i-- {
		Label(labels[i])
		MOVB(Mem{Base: reg_m, Disp: i - 1}, char)
		MOVB(char, Mem{Base: lb, Disp: i - 1})
	}

	Label(labels[0])

	ADDQ(reg_m_len, lb)
	MOVB(Imm(1), Mem{Base: lb})
	SUBQ(reg_m_len, lb)

	XORL(Mem{Base: lb}, v0)
	XORL(Mem{Base: lb, Disp: 4}, v1)
	XORL(Mem{Base: lb, Disp: 8}, v2)
	XORL(Mem{Base: lb, Disp: 12}, v3)

	Label("afterSwitch")

	XORL(Mem{Base: lkey}, v0)
	XORL(Mem{Base: lkey, Disp: 4}, v1)
	XORL(Mem{Base: lkey, Disp: 8}, v2)
	XORL(Mem{Base: lkey, Disp: 12}, v3)

	permute(reg_r, v0, v1, v2, v3)

	XORL(Mem{Base: lkey}, v0)
	XORL(Mem{Base: lkey, Disp: 4}, v1)
	XORL(Mem{Base: lkey, Disp: 8}, v2)
	XORL(Mem{Base: lkey, Disp: 12}, v3)

	ret := GP64()
	Load(Param("tag").Base(), ret)

	MOVL(v0, Mem{Base: ret})
	MOVL(v1, Mem{Base: ret, Disp: 4})
	MOVL(v2, Mem{Base: ret, Disp: 8})
	MOVL(v3, Mem{Base: ret, Disp: 12})

	RET()
	Generate()
}

type loop string

var loops int

func newLoop() loop {
	loops++
	return loop("loop" + strconv.Itoa(loops-1))
}

func (l loop) begin() string { return string(l) + "_begin" }
func (l loop) end() string   { return string(l) + "_end" }
