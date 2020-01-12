# Go Virtual Machine

A simple virtual machine written in Go that started out as a way to pass time
during the holidays.

## The GVM

The GVM is a very simple virtual machine where you can play with with integer
registers, a call stack and a stack.

## GBF: GVM Binary File

The GVM executes GVM binary files, which are produced by the compiler.

These are files written with little endian encoding, containing a header used
by the GVM to check binary compatibility between itself and the compiled file
version. The header is then followed by the number of elements in the code
array, followed by the code array itself.

## GSM: GVM Assembly Language

The GVM's assembly language is what you use to write code that will be compiled
into actual bytecode. The proposed extension is 'gsm', though, of course, this
makes no difference.

There are some quality of life simplifications for the written language that
make it not an _exact_ map to the GVM compiled bytecode, but all the exceptions
are laid out below.

There are a few examples of GSM programs located in the `examples` folder at
the root of this repository. They hopefully can serve as self-explanatory and
reasonable examples of GSM usage.

### Registers

Integer registers are referred to as `r<n>` where `<n>` indicates the number
of the register. For example `r1` refers to the first integer register and
`r16` refers to the sixteenth integer register.

### Labels

Labels are used to simplify control flow, being essentially an abstraction for
the following instruction address.

They are written as tokens ending with `:`, with the name of the label being
everything before the `:`.

One can also use sublabels. The basic idea is that after any given label `lab`
is defined, any labels starting with a `.` will be a sublabel of `lab`, being
effectively a short alias for `lab.<sublabel>` and the sublabel. For example,
in the code:
```
main:
    set r1 2
.second
    set r2 5
    jmp .second
```
the label `.second` is effectively an independent label named `main.second`.

Note that this implies a scope generated by common labels. To illustrate this,
consider the example below:
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
which is effectively the same as
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

Note also that if a label named `main` is defined, an instruction is added
implicitly by the compiler to the start of the code sequence that jumps to
this label, so that it can be used as a logical entry point for the program.

### Comments

Comments start with the character `;`. As with usual comments, both the
character `;` and everything that follows it will be ignored by the compiler.

### Instructions

This section needs to be written.
