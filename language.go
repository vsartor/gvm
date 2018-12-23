package main

// Instruction tokens
const (
	IThalt = iota
	ITset
	ITadd
	ITsub
	ITmul
	ITdiv
	ITrem
	ITjmp
	// ITcmp
	// ITjeq
	// ITjnq
	// ITjgt
	// ITjlt
	// ITjge
	// ITjle
	ITshow
)

var (
	it2str map[int64]string
	str2it map[string]int64
)

func init() {
	it2str = make(map[int64]string)
	str2it = make(map[string]int64)

	it2str[IThalt] = "halt"
	it2str[ITset] = "set"
	it2str[ITadd] = "add"
	it2str[ITsub] = "sub"
	it2str[ITmul] = "mul"
	it2str[ITdiv] = "div"
	it2str[ITrem] = "rem"
	it2str[ITjmp] = "jmp"
	it2str[ITshow] = "show"

	for itok, stok := range it2str {
		str2it[stok] = itok
	}
}
