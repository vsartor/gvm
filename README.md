GVM
===

A simple virtual machine written in Go for no purpose other than passing some time during the holidays.

## Virtual Machine

It is a register based virtual machine, with 16 integer registers and 16 floating point registers. From now on, if the type of the register is not specified, it is an integer register.

## Assembly Language

The assembly language is written in text format in files with the *.gsm extension. There are simplifications for the written language that make it not an exact map to the VM instructions when compiled, but all exceptions are laid out here.

In general, there will be examples located in the very well named `examples` folder at the root of the repository. They should be self-explanatory and a reasonable example of usage of the language.

### Registers

Registers are referred to as `rx` where `x` indicates the number of the register. For example `r1` refers to the first register and `r16` refers to the sixteenth register. The same convention is used for floating point registers, however the letter `f` is used instead of `r`.

### Labels

Labels are tokens ending with `:`, and the name of the label is everything before the `:`. Labels are used for jumping instructions, so in the gsm file you jump to labels instead of literal code positions.

One can also use sublabels. The basic idea is that after any given label `lab` is defined, any labels starting with a `.` will be a sublabel, being renamed internally to be a concatenation between `lab` and the sublabel. For example, in the code
```
main:
    set r1 2
.second
    set r2 5
    jmp .second
```
the label `.second` is renamed to `main.second`.

Note however, that references to sublabels such as `.second` will also expand as a sublabel of the current scope. For example, the code
```
main1:
    set r1 2
.second
    set r2 5
    jmp .second
main2:
    set r1 2
.second
    set r2 5
    jmp .second
```
is equivalent to
```
main1:
    set r1 2
main1.second
    set r2 5
    jmp main1.second
main2:
    set r1 2
main2.second
    set r2 5
    jmp main2.second
```

If the label `main` is defined, an instruction is added implicitly to the start of the code sequence that jumps to this label.

### Comments

Comments start with the character `;`. Both the character `;` and everything that follows it, will be ignored by the compiler.

### Instructions

Token|Description
-----|-----------
`halt`| Stops program execution.
`set rx val`| Sets register `x` to value `val`.
`add rx ry`| Adds register `x` to register `y`.
`sub rx ry`| Subtracts register `x` to register `y`.
`mul rx ry`| Multiplies register `x` to register `y`.
`div rx ry`| Divides register `y` by register `x`.
`rem rx ry`| Remainder of division of register `y` by register `x`.
`cmp rx ry`| Compares register `y` in relation to `x`.
`jmp lab` | Jumps to position at label `lab`.
`jeq lab` | Jumps to position at label `lab` if comparison was equal.
`jne lab` | Jumps to position at label `lab` if comparison was not equal.
`jgt lab` | Jumps to position at label `lab` if `ry` was larger than `rx`.
`jlt lab` | Jumps to position at label `lab` if `ry` was lesser than `rx`.
`jge lab` | Jumps to position at label `lab` if `ry` was larger or equal than `rx`.
`jle lab` | Jumps to position at label `lab` if `ry` was lesser or equal than `rx`.
`show rx` | Displays the register `x`.

### Details

* At the very beginning of the program, a label named `_zero` is implicitly defined.
* After finishing parsing the program, a `halt` instruction is always added at the end.

### Object Files

Object (compiled) files are binary files with little endian encoding and extension *.gbf.

Object files start with an arbitrary header that will change as incompatible changes are introduced, to help avoid problematic errors with incompatible binary formats. The header is then followed by the number of elements in the code array. Then, a straight-up binary encoding of the code array is present.
