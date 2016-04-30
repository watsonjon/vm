package main

import "fmt"

// set of operations
const (
	PUSH = iota
	ADD
	PRINT
	HALT
	JMPLT
	SUB
	MUL
	STORE
	LOAD
	POP
	JMP
	JMPGT
	JMPEQ
)

type op struct {
	name  string
	nargs int
}

var ops = map[int]op{
	PUSH:  op{"push", 1},
	POP:   op{"pop", 0},
	ADD:   op{"add", 0},
	PRINT: op{"print", 0},
	HALT:  op{"halt", 0},
	JMP:   op{"jmp", 1},
	JMPLT: op{"jmplt", 2},
	JMPGT: op{"jmpgt", 2},
	JMPEQ: op{"jmpeq", 2},
	SUB:   op{"sub", 0},
	MUL:   op{"mul", 0},
	STORE: op{"store", 1},
	LOAD:  op{"load", 1},
}

// VM is our virtual machine
type VM struct {
	code []int
	pc   int

	stack []int
	sp    int

	mem []int
}

func (v *VM) trace() {
	addr := v.pc
	op := ops[v.code[v.pc]]
	args := v.code[v.pc+1 : v.pc+op.nargs+1]
	stack := v.stack[0 : v.sp+1]
	fmt.Printf("%04d: %s\t%v\t%v\n", addr, op.name, args, stack)
}

func (v *VM) dump() {
	fmt.Println("DATA:")
	for addr, data := range v.mem {
		fmt.Printf("%04d: %v\n", addr, data)
	}
}

func (v *VM) run(c []int) {

	v.stack = make([]int, 100)
	v.mem = make([]int, 10)
	v.sp = -1

	v.code = c
	v.pc = 0

	for {
		v.trace()

		// Fetch
		op := v.code[v.pc]
		v.pc++

		// Decode
		switch op {
		case PUSH:
			val := v.code[v.pc]
			v.pc++

			v.sp++
			v.stack[v.sp] = val
		case POP:
			v.sp--
		case ADD:
			b := v.stack[v.sp]
			v.sp--
			a := v.stack[v.sp]
			v.sp--

			v.sp++
			v.stack[v.sp] = a + b
		case SUB:
			b := v.stack[v.sp]
			v.sp--
			a := v.stack[v.sp]
			v.sp--

			v.sp++
			v.stack[v.sp] = a - b
		case MUL:
			b := v.stack[v.sp]
			v.sp--
			a := v.stack[v.sp]
			v.sp--

			v.sp++
			v.stack[v.sp] = a * b
		case PRINT:
			val := v.stack[v.sp]
			v.sp--
			fmt.Println(val)
		case JMP:
			addr := v.code[v.pc]
			v.pc = addr
		case JMPLT:
			val := v.code[v.pc]
			v.pc++
			addr := v.code[v.pc]
			v.pc++

			if v.stack[v.sp] < val {
				v.pc = addr
			}
		case JMPGT:
			val := v.code[v.pc]
			v.pc++
			addr := v.code[v.pc]
			v.pc++

			if v.stack[v.sp] > val {
				v.pc = addr
			}
		case JMPEQ:
			val := v.code[v.pc]
			v.pc++
			addr := v.code[v.pc]
			v.pc++

			if v.stack[v.sp] == val {
				v.pc = addr
			}
		case LOAD:
			addr := v.code[v.pc]
			v.pc++
			val := v.mem[addr]
			v.sp++
			v.stack[v.sp] = val
		case STORE:
			val := v.stack[v.sp]
			v.sp--
			addr := v.code[v.pc]
			v.pc++
			// append(v.mem, val)
			v.mem[addr] = val
		case HALT:
			v.dump()
			return
		}
	}
}

func main() {
	code := []int{
		PUSH, 2,
		PUSH, 3,
		// JMP, 33,
		ADD,
		PRINT,
		PUSH, 15,
		PUSH, 5,
		SUB,
		PRINT,
		PUSH, 3,
		PUSH, 3,
		MUL,
		JMPLT, 10, 14,
		POP,
		PUSH, 88,
		STORE, 0,
		PUSH, 7,
		POP,
		LOAD, 0,
		PRINT,
		HALT,
	}

	v := &VM{}
	v.run(code)
}
