package asm

type Registers [6]int
type Instruction struct {
	F        Op
	Operands [3]int
}

type Op func(Registers, int, int, int) Registers

func Addr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] + r[operB]
	return result
}
func Addi(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] + operB
	return result
}
func Mulr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] * r[operB]
	return result
}
func Muli(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] * operB
	return result
}
func Banr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] & r[operB]
	return result
}
func Bani(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] & operB
	return result
}
func Borr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] | r[operB]
	return result
}
func Bori(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA] | operB
	return result
}
func Setr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = r[operA]
	return result
}
func Seti(r Registers, operA, operB, targetReg int) Registers {
	result := r
	result[targetReg] = operA
	return result
}
func Gtir(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if operA > r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Gtri(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] > operB {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Gtrr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] > r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Eqir(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if operA == r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Eqri(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] == operB {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}
func Eqrr(r Registers, operA, operB, targetReg int) Registers {
	result := r
	if r[operA] == r[operB] {
		result[targetReg] = 1
	} else {
		result[targetReg] = 0
	}
	return result
}

var AllOps = map[string]Op{
	"addr": Addr,
	"addi": Addi,
	"mulr": Mulr,
	"muli": Muli,
	"banr": Banr,
	"bani": Bani,
	"borr": Borr,
	"bori": Bori,
	"setr": Setr,
	"seti": Seti,
	"gtir": Gtir,
	"gtri": Gtri,
	"gtrr": Gtrr,
	"eqir": Eqir,
	"eqri": Eqri,
	"eqrr": Eqrr,
}
