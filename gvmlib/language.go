package gvmlib

// Instruction tokens
const (
	IThalt = iota
	ITset
	ITpush
	ITpop
	ITinc
	ITdec
	ITmov
	ITadd
	ITsub
	ITmul
	ITdiv
	ITrem
	ITcmp
	ITjmp
	ITjeq
	ITjne
	ITjgt
	ITjlt
	ITjge
	ITjle
	ITshow
	ITcall
	ITret
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
	it2str[ITpush] = "push"
	it2str[ITpop] = "pop"
	it2str[ITinc] = "inc"
	it2str[ITdec] = "dec"
	it2str[ITmov] = "mov"
	it2str[ITadd] = "add"
	it2str[ITsub] = "sub"
	it2str[ITmul] = "mul"
	it2str[ITdiv] = "div"
	it2str[ITrem] = "rem"
	it2str[ITcmp] = "cmp"
	it2str[ITjmp] = "jmp"
	it2str[ITjeq] = "jeq"
	it2str[ITjne] = "jne"
	it2str[ITjgt] = "jgt"
	it2str[ITjlt] = "jlt"
	it2str[ITjge] = "jge"
	it2str[ITjle] = "jle"
	it2str[ITshow] = "show"
	it2str[ITcall] = "call"
	it2str[ITret] = "ret"

	for itok, stok := range it2str {
		str2it[stok] = itok
	}
}
