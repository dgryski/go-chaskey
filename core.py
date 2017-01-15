import peachpy.x86_64


def permute(r, v0, v1, v2, v3):
    i = GeneralPurposeRegister64()
    MOV(i, r)
    with Loop() as loop:
        ADD(v0, v1)
        ROL(v1, 5)
        XOR(v1, v0)
        ROL(v0, 16)
        ADD(v2, v3)
        ROL(v3, 8)
        XOR(v3, v2)
        ADD(v0, v3)
        ROL(v3, 13)
        XOR(v3, v0)
        ADD(v2, v1)
        ROL(v1, 7)
        XOR(v1, v2)
        ROL(v2, 16)
        DEC(i)
        JNZ(loop.begin)


def MakeChaskeyCore():
    h = Argument(ptr())
    m_base = Argument(ptr())
    m_len = Argument(int64_t)
    m_cap = Argument(int64_t)

    tag_base = Argument(ptr())
    tag_len = Argument(int64_t)
    tag_cap = Argument(int64_t)

    with Function(
            "chaskeyCore",
        (h, m_base, m_len, m_cap, tag_base, tag_len, tag_cap),
            target=uarch.default) as function:

        v0, v1, v2, v3 = GeneralPurposeRegister32(), GeneralPurposeRegister32(
        ), GeneralPurposeRegister32(), GeneralPurposeRegister32()

        reg_h = GeneralPurposeRegister64()
        LOAD.ARGUMENT(reg_h, h)

        MOV(v0, [reg_h])
        MOV(v1, [reg_h + 4])
        MOV(v2, [reg_h + 8])
        MOV(v3, [reg_h + 12])

        reg_r = GeneralPurposeRegister64()
        MOV(reg_r, [reg_h + 48])
        rounds = reg_r

        reg_m = GeneralPurposeRegister64()
        reg_m_len = GeneralPurposeRegister64()
        LOAD.ARGUMENT(reg_m, m_base)
        LOAD.ARGUMENT(reg_m_len, m_len)

        loop = Loop()
        CMP(reg_m_len, 16)
        JLE(loop.end)
        with loop:
            XOR(v0, [reg_m])
            XOR(v1, [reg_m + 4])
            XOR(v2, [reg_m + 8])
            XOR(v3, [reg_m + 12])

            permute(rounds, v0, v1, v2, v3)

            ADD(reg_m, 16)
            SUB(reg_m_len, 16)
            CMP(reg_m_len, 16)
            JG(loop.begin)

        switch = Label()
        afterSwitch = Label()

        lkey = GeneralPurposeRegister64()
        MOV(lkey, reg_h)

        CMP(reg_m_len, 16)
        JNZ(switch)

        ADD(lkey, 16)

        XOR(v0, [reg_m])
        XOR(v1, [reg_m + 4])
        XOR(v2, [reg_m + 8])
        XOR(v3, [reg_m + 12])

        JMP(afterSwitch)

        LABEL(switch)
        ADD(lkey, 32)

        lb = GeneralPurposeRegister64()
        LOAD.ARGUMENT(lb, tag_base)

        MOV(qword[lb], 0)
        MOV(qword[lb + 8], 0)

        # no support for jump tables
        labels = [Label("sw%d" % i) for i in range(0, 16)]

        for i in range(0, 16):
            CMP(reg_m_len, i)
            JE(labels[i])

        char = GeneralPurposeRegister8()
        for i in range(15, 0, -1):
            LABEL(labels[i])
            MOV(char, byte[reg_m + i - 1])
            MOV(byte[lb + i - 1], char)

        LABEL(labels[0])

        ADD(lb, reg_m_len)
        MOV(byte[lb], 0x01)
        SUB(lb, reg_m_len)

        XOR(v0, [lb])
        XOR(v1, [lb + 4])
        XOR(v2, [lb + 8])
        XOR(v3, [lb + 12])

        LABEL(afterSwitch)

        XOR(v0, [lkey])
        XOR(v1, [lkey + 4])
        XOR(v2, [lkey + 8])
        XOR(v3, [lkey + 12])

        permute(rounds, v0, v1, v2, v3)

        XOR(v0, [lkey])
        XOR(v1, [lkey + 4])
        XOR(v2, [lkey + 8])
        XOR(v3, [lkey + 12])

        ret = GeneralPurposeRegister64()
        LOAD.ARGUMENT(ret, tag_base)

        MOV([ret], v0)
        MOV([ret + 4], v1)
        MOV([ret + 8], v2)
        MOV([ret + 12], v3)

        RETURN()


MakeChaskeyCore()
