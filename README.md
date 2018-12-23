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

### Comments

Comments start with the character `;`. Both the character `;` and everything that follows it, will be ignored by the compiler.

### Instructions

Token|Description
-----|-----------
`halt`| Stops program execution
`set rx val`| Sets register `x` to value `val`
`add rx ry`| Adds value of register `x` to register `y`
`show rx` | Displays the value of the register `x`

### Details

* After finishing parsing the program, a `halt` instruction is always added at the end.

### Object Files

Object (compiled) files are binary files with little endian encoding and extension *.gbf.

Object files start with an arbitrary header that will change as incompatible changes are introduced, to help avoid problematic errors with incompatible binary formats. The header is then followed by the number of elements in the code array. Then, a straight-up binary encoding of the code array is present.
