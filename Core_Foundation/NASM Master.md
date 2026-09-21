# Assembly Language — Scratch → Principal Systems Engineer Masterclass

මේක **x86-64 Assembly එක CPU level එකෙන් පටන්ගෙන Windows / Linux / macOS OS boundary එක දක්වා** සම්බන්ධ කරන guide එකක්. මෙහි වැදගත්ම mental model එක:

> **CPU architecture එක එකම දෙයක්. Assembly source syntax, ABI, executable format, linker, loader, system-call interface වගේ දේවල් OS/toolchain එක අනුව වෙනස් වෙනවා.**

ඒ නිසා **`MOV` instruction එක CPU එකේ instruction එකක්** වුණත්, source code එකේ ඒක ලියන විදිහ, function එකකට argument දෙන register එක, executable එකේ symbol format එක, OS එකට service එකක් request කරන විදිහ වෙනස් වෙන්න පුළුවන්.

---

# 1. ARCHITECTURAL FOUNDATIONS — OS-Agnostic Core

## 1.1 මුලින්ම Computer එකේ layers ටික

Assembly ඉගෙනගන්න කලින් මේ hierarchy එක හොඳට තේරුම් ගන්න.

```text
┌───────────────────────────────────────────────┐
│              Your Application                 │
│          C / C++ / Rust / Go / C#             │
├───────────────────────────────────────────────┤
│             Compiler / Linker                 │
│        GCC / Clang / MSVC / LLVM / etc.      │
├───────────────────────────────────────────────┤
│                Assembly                      │
│       NASM / MASM / GAS / LLVM syntax         │
├───────────────────────────────────────────────┤
│               Machine Code                    │
│          Bytes → Opcodes + Operands            │
├───────────────────────────────────────────────┤
│             CPU Architecture                  │
│       x86-64 / ARM64 / RISC-V etc.            │
├───────────────────────────────────────────────┤
│                  CPU                        │
│ Registers / ALU / FPU / MMU / Cache / etc.   │
├───────────────────────────────────────────────┤
│                 Hardware                      │
│         RAM / SSD / GPU / Devices             │
└───────────────────────────────────────────────┘
```

ඒ අතර OS එක මේ layers අතර ඉතා වැදගත් boundary එකක්:

```text
Application
    │
    │ function call
    ▼
Libraries / Runtime
    │
    │ system call / API
    ▼
Kernel
    │
    │ drivers
    ▼
Hardware
```

---

# 1.2 CPU එකේ fundamental resources

x86-64 CPU එකක් ගැන කතා කරනකොට ප්‍රධාන state එක:

```text
                 CPU
        ┌───────────────────┐
        │ General Registers │
        │ RAX RBX RCX ...   │
        │                   │
        │ RIP               │ ← instruction pointer
        │ RFLAGS            │ ← status/control flags
        │                   │
        │ XMM/YMM/ZMM       │ ← SIMD / FP
        └─────────┬─────────┘
                  │
             Load / Store
                  │
                  ▼
        ┌───────────────────┐
        │       Cache       │
        ├───────────────────┤
        │        RAM        │
        └───────────────────┘
```

CPU එකට **C variable**, **Go variable**, **Python variable** කියලා concept එකක් hardware level එකේ නැහැ.

CPU එකට තිබෙන්නේ:

* bytes
* addresses
* registers
* instructions
* flags
* memory

---

# 1.3 Variable එකක් Assembly වලට map වෙන්නේ කොහොමද?

C:

```c
int x = 42;
```

Compiler එකට choices තියෙනවා.

### Register variable

```text
RAX = 42
```

### Stack variable

```text
[RBP-4] = 42
```

### Global variable

```text
.data
x: dd 42
```

එනම්:

```text
C variable
   │
   ├── Register
   │
   ├── Stack memory
   │
   ├── Heap memory
   │
   └── Static/Data memory
```

**Variable කියන්නේ abstract programming-language concept එකක්.**

Assembly එකේදී ඒක eventually:

> **register එකක value එකක් හෝ memory address එකකට අදාළ bytes ටිකක්**

වෙනවා.

---

# 1.4 Primitive data types → bytes

CPU එකට `int`, `bool`, `char` කියලා semantic type information සාමාන්‍යයෙන් නැහැ.

උදාහරණයක්:

```text
00000000 00000000 00000000 00101010
```

මේ bytes ටික C `int` එකක් වෙන්න පුළුවන්.

එම bytes ටික:

```text
42
```

වෙන්නත් පුළුවන්.

නැත්නම් character data එකක් වෙන්නත් පුළුවන්.

### Common sizes

| High-level type |    Typical size | x86-64 concept    |
| --------------- | --------------: | ----------------- |
| `char`          |          1 byte | byte              |
| `bool`          | 1 byte commonly | byte              |
| `short`         |               2 | word              |
| `int`           |               4 | dword             |
| `long long`     |               8 | qword             |
| pointer         |               8 | qword             |
| `float`         |               4 | IEEE-754 binary32 |
| `double`        |               8 | IEEE-754 binary64 |

Assembly/NASM terminology:

```text
DB = Define Byte       1 byte
DW = Define Word       2 bytes
DD = Define Doubleword 4 bytes
DQ = Define Quadword   8 bytes
```

Example:

```asm
section .data

a db 10
b dw 1000
c dd 100000
d dq 10000000000
```

Memory:

```text
a → [ 0A ]

b → [ E8 03 ]

c → [ A0 86 01 00 ]

d → [ ... 8 bytes ... ]
```

x86-64 is **little-endian**.

ඒ නිසා `1000 = 0x03E8` memory එකේ:

```text
low address
   ↓
┌────┬────┐
│ E8 │ 03 │
└────┴────┘
```

---

# 1.5 Registers

x86-64 general-purpose registers:

```text
RAX
RBX
RCX
RDX

RSI
RDI
RBP
RSP

R8
R9
R10
R11
R12
R13
R14
R15
```

64-bit register එකක sub-registers තියෙනවා.

```text
                RAX
        ┌──────────────────┐
        │      64 bits     │
        └──────────────────┘
              EAX
        ┌────────────┐
        │ lower 32   │
        └────────────┘
             AX
        ┌──────┬──────┐
        │ AH   │ AL   │
        └──────┴──────┘
```

උදාහරණ:

```asm
mov rax, 123
```

64-bit value.

```asm
mov eax, 123
```

32-bit value.

```asm
mov ax, 123
```

16-bit.

```asm
mov al, 123
```

8-bit.

### Important x86-64 rule

32-bit register එකකට write කිරීම:

```asm
mov eax, 123
```

`RAX` හි upper 32 bits zero වෙනවා.

```text
Before:
RAX = FFFFFFFFFFFFFFFF

mov eax, 123

After:
RAX = 000000000000007B
```

මේක x86-64 code එකේ ඉතා වැදගත් behavior එකක්.

---

# 1.6 RAM

RAM කියන්නේ huge byte array එකක් වගේ imagine කරන්න.

```text
Address
00000000 → [ byte ]
00000001 → [ byte ]
00000002 → [ byte ]
00000003 → [ byte ]
...
```

CPU එක:

```asm
mov rax, [rbx]
```

කියන්නේ:

> `RBX` තුළ තිබෙන address එකට ගිහින් memory එකෙන් 8 bytes load කරලා RAX එකට දාන්න.

```text
RBX
 │
 │ address
 ▼
RAM
┌────────────────────┐
│ 8 bytes             │
└────────────────────┘
          │
          ▼
         RAX
```

`[]` කියන්නේ NASM Intel syntax එකේ:

> **memory dereference**

---

# 1.7 Immediate vs Register vs Memory

```asm
mov rax, 10
```

`10` = immediate constant.

```asm
mov rax, rbx
```

`RBX` = register operand.

```asm
mov rax, [rbx]
```

`[RBX]` = memory operand.

```text
10
↓
constant

RBX
↓
register

[RBX]
↓
memory at address contained in RBX
```

---

# 1.8 Addressing Modes

x86-64 addressing expression:

```asm
[base + index*scale + displacement]
```

General form:

```text
[ Base + Index × Scale + Offset ]
```

Scale:

```text
1
2
4
8
```

Example:

```asm
mov rax, [rbx + rcx*8 + 16]
```

CPU calculates:

```text
effective_address =
    RBX
  + RCX × 8
  + 16
```

ඉතා practical example:

```c
array[i]
```

if `array` contains 8-byte elements:

```asm
mov rax, [array + rcx*8]
```

because:

```text
element_size = 8
index = i
address = array + i × 8
```

---

# 1.9 Opcodes

Assembly:

```asm
mov rax, 42
```

CPU එක මේ text එක execute කරන්නේ නැහැ.

Assembler එක මේක machine code bytes වලට convert කරනවා.

```text
Assembly source
      │
      ▼
Assembler
      │
      ▼
Machine-code bytes
      │
      ▼
CPU decoder
      │
      ▼
Execution
```

Machine instruction roughly:

```text
┌────────────┬──────────────────┐
│   Opcode   │     Operands     │
└────────────┴──────────────────┘
```

x86 instructions variable-length.

ඒක RISC architectures වල fixed-width instruction model එකෙන් වෙනස්.

---

# 1.10 Operand

Example:

```asm
add rax, rbx
```

මෙහි:

```text
ADD = mnemonic
RAX = destination operand
RBX = source operand
```

Conceptually:

```text
RAX = RAX + RBX
```

---

# 1.11 Flags

`RFLAGS` register එකේ CPU status/control bits තියෙනවා.

ඔබ ඉල්ලපු 4 ප්‍රධාන flags:

| Flag | Meaning       |
| ---- | ------------- |
| ZF   | Zero Flag     |
| CF   | Carry Flag    |
| SF   | Sign Flag     |
| OF   | Overflow Flag |

---

## ZF — Zero Flag

```asm
cmp rax, rbx
```

CPU internally computes:

```text
RAX - RBX
```

result = 0 නම්:

```text
ZF = 1
```

Example:

```asm
mov rax, 10
cmp rax, 10
je equal
```

`JE` = Jump if Equal.

Actually condition:

```text
ZF = 1
```

---

# 1.12 CF — Carry Flag

Unsigned arithmetic overflow/carry සඳහා වැදගත්.

Example conceptual 8-bit:

```text
255 + 1
```

```text
11111111
+00000001
---------
00000000
```

carry එකක් ඇතිවෙනවා.

```text
CF = 1
```

---

# 1.13 SF — Sign Flag

Signed result එකේ highest bit එක 1 නම්:

```text
SF = 1
```

Example 8-bit:

```text
10000000
```

two's complement signed representation එකේ negative value.

---

# 1.14 OF — Overflow Flag

Signed arithmetic overflow.

8-bit signed range:

```text
-128 → +127
```

```text
127 + 1
```

binary:

```text
01111111
+00000001
---------
10000000
```

CPU එක signed interpretation එකේ:

```text
+127 + +1
```

cannot fit.

```text
OF = 1
```

---

# 1.15 `CMP` + `Jcc`

High-level:

```c
if (a == b) {
    ...
}
```

Assembly:

```asm
cmp rax, rbx
je equal
```

Conceptually:

```text
CMP
 │
 │ RAX - RBX
 ▼
FLAGS
 │
 ├── ZF
 ├── CF
 ├── SF
 └── OF
 │
 ▼
JE
 │
 └── if ZF=1 → jump
```

Common conditional jumps:

| Assembly      | Meaning        |
| ------------- | -------------- |
| `JE` / `JZ`   | equal / zero   |
| `JNE` / `JNZ` | not equal      |
| `JG`          | signed greater |
| `JL`          | signed less    |
| `JGE`         | signed ≥       |
| `JLE`         | signed ≤       |
| `JA`          | unsigned >     |
| `JB`          | unsigned <     |
| `JAE`         | unsigned ≥     |
| `JBE`         | unsigned ≤     |
| `JC`          | carry          |
| `JO`          | overflow       |
| `JS`          | sign           |

**Signed සහ unsigned comparison එක වෙනස්.**

ඒක Assembly mastery එකේ fundamental concept එකක්.

---

# 1.16 If / Else → Jumps

C:

```c
if (x == 0) {
    A();
} else {
    B();
}
```

Assembly structure:

```asm
cmp rax, 0
jne else_branch

; A()
jmp end_if

else_branch:
; B()

end_if:
```

Control-flow graph:

```text
             CMP
              │
          ZF / condition
         /            \
       true           false
        │               │
        ▼               ▼
       A()             B()
        │               │
        └──────┬────────┘
               ▼
             END
```

---

# 1.17 Loops

C:

```c
while (x != 0) {
    x--;
}
```

Assembly:

```asm
loop_start:

    cmp rax, 0
    je loop_end

    dec rax
    jmp loop_start

loop_end:
```

CPU perspective:

```text
          ┌──────────────┐
          │ loop_start   │
          └──────┬───────┘
                 ↓
               CMP
                 │
            ┌────┴────┐
            │         │
           !=         ==
            │         │
            ▼         ▼
           DEC       END
            │
            └──────► loop_start
```

---

# 1.18 Functions

High-level:

```c
int add(int a, int b) {
    return a + b;
}
```

Conceptually:

```asm
add:
    ; arguments according to ABI
    add ...
    ret
```

Caller:

```asm
call add
```

`CALL` internally does roughly:

```text
push return_address
jump to function
```

`RET`:

```text
pop return_address
jump there
```

Conceptually:

```text
CALL
 │
 ├── save RIP of next instruction
 │
 └── jump to function

             function
                │
               RET
                │
                ▼
        saved return address
                │
                ▼
             caller
```

---

# 1.19 Stack

Stack grows toward lower addresses on x86-64.

```text
Higher addresses
────────────────────
│ caller data       │
├────────────────────┤
│ return address     │
├────────────────────┤
│ local variables    │
├────────────────────┤
│ saved registers    │
├────────────────────┤ ← RSP
│                    │
▼ stack grows down
────────────────────
Lower addresses
```

Typical function:

```asm
push rbp
mov rbp, rsp

; locals

pop rbp
ret
```

Conceptual stack frame:

```text
             HIGH ADDRESS
┌─────────────────────────┐
│ caller's stack          │
├─────────────────────────┤
│ return address          │ ← [RBP+8]
├─────────────────────────┤
│ old RBP                 │ ← [RBP]
├─────────────────────────┤
│ local variable          │ ← [RBP-8]
├─────────────────────────┤
│ another local           │
└─────────────────────────┘
             ↓
          LOW ADDRESS
```

Modern optimized compilers **may omit RBP as frame pointer**, so don't assume every function has this exact structure.

---

# 1.20 Data / BSS / Text

Typical executable memory model:

```text
Virtual Address Space

HIGH
┌──────────────────────────┐
│ Kernel space             │
├──────────────────────────┤
│ Stack ↓                  │
│                          │
│                          │
│ Heap ↑                   │
├──────────────────────────┤
│ BSS                      │
│ Data                     │
│ Read-only data           │
│ Text / Code              │
└──────────────────────────┘
LOW
```

Assembler sections:

```asm
section .text
```

machine code.

```asm
section .data
```

initialized writable data.

```asm
section .bss
```

uninitialized zero-initialized storage.

Example:

```asm
section .data
message db "Hello", 10

section .bss
buffer resb 100

section .text
```

---

# 2. THE "WHY" — Why Assembly Diverges Between OSs

මෙතනින් තමයි **really important distinction** එක එන්නේ.

## CPU architecture ≠ OS ABI ≠ assembler syntax

මේ තුන එකම දෙයක් නෙවෙයි.

```text
                x86-64 ISA
                    │
       ┌────────────┼─────────────┐
       │            │             │
     Windows      Linux         macOS
       │            │             │
      ABI          ABI           ABI
       │            │             │
      PE           ELF          Mach-O
       │            │             │
    MASM/NASM     GAS/NASM     Clang/NASM
```

CPU instruction set එක common.

නමුත් surrounding conventions වෙනස්.

---

# 2.1 Intel syntax vs AT&T syntax

### Intel/NASM

```asm
mov rax, rbx
```

### AT&T/GAS

```asm
movq %rbx, %rax
```

Meaning එක එකමයි:

```text
RAX = RBX
```

Difference:

### Intel

```text
destination, source
```

### AT&T

```text
source, destination
```

ඒක beginner කෙනෙක්ට ඉතා confusing.

---

# 2.2 Memory syntax

Intel:

```asm
mov rax, [rbx + rcx*8]
```

AT&T:

```asm
movq (%rbx,%rcx,8), %rax
```

Architecture එක වෙනස් වෙලා නෑ.

**Notation එක විතරයි වෙනස්.**

---

# 2.3 NASM

NASM = Netwide Assembler.

සාමාන්‍යයෙන් Intel-style syntax.

```asm
section .text
global _start

_start:
    mov rax, 60
    xor rdi, rdi
    syscall
```

---

# 2.4 MASM

Microsoft Macro Assembler.

Windows ecosystem එකේ common.

Example:

```asm
mov rax, 60
```

MASM directives සහ syntax NASM එකට identical නෙවෙයි.

---

# 2.5 GAS

GNU Assembler.

Linux toolchain එකේ historically dominant.

Default syntax:

```asm
movq %rax, %rbx
```

GAS එක Intel syntax support කරන්නත් පුළුවන්.

---

# 2.6 LLVM

LLVM ecosystem එකට තමන්ගේ assembly parsing / MC infrastructure තියෙනවා.

Clang → LLVM IR → machine code / object.

Conceptual pipeline:

```text
C/C++
   ↓
Clang
   ↓
LLVM IR
   ↓
Optimization
   ↓
Machine instruction selection
   ↓
Object file
   ↓
Linker
   ↓
Executable
```

---

# 2.7 Executable format

Assembly source file එක executable එකක් නෙවෙයි.

Assembler:

```text
.asm
 ↓
object file
```

Linker:

```text
object + libraries
 ↓
executable
```

OS loader:

```text
executable
 ↓
virtual memory mappings
 ↓
process
```

Major formats:

| OS      | Executable format |
| ------- | ----------------- |
| Windows | PE / COFF         |
| Linux   | ELF               |
| macOS   | Mach-O            |

---

# 2.8 PE

Windows:

```text
.exe
.dll
```

PE structure contains things such as:

```text
DOS header
PE signature
COFF header
Optional header
Section table
.text
.rdata
.data
.rsrc
.reloc
...
```

Windows loader ඒ information use කරලා process address space එක construct කරනවා.

---

# 2.9 ELF

Linux:

```text
ELF Header
Program Headers
Sections
.text
.rodata
.data
.bss
...
```

Important distinction:

> **Sections are primarily linker/object organization; program headers describe loadable runtime segments.**

Dynamic ELF executable එකක්:

```text
ELF
 │
 ├── code
 ├── data
 ├── dynamic metadata
 └── interpreter
          │
          ▼
       dynamic loader
```

---

# 2.10 Mach-O

macOS executable format:

```text
Mach-O Header
Load Commands
Segments
Sections
Symbol information
...
```

macOS uses Mach-O rather than ELF.

---

# 2.11 Symbol naming

Source:

```c
int main()
```

Object file එකේ symbol name එක platform/toolchain ABI අනුව වෙනස් විය හැක.

Historically macOS x86-64 C symbols are commonly represented with leading underscore:

```text
main
 ↓
_main
```

ඒ නිසා NASM Mach-O code එකේ:

```asm
global _main
_main:
```

දකින්න පුළුවන්.

Linux ELF:

```asm
global main
main:
```

Windows toolchain/configuration අනුව symbol decoration rules වෙනස් විය හැක.

---

# 3. WINDOWS vs LINUX vs macOS

දැන් වැදගත් distinction එක:

> **Calling convention ≠ system-call convention.**

Function call:

```text
your function → another function
```

System call:

```text
user process → kernel
```

---

# 3.1 Linux x86-64 System Call ABI

Linux x86-64 standard syscall mechanism:

```asm
syscall
```

Registers:

| Purpose        | Register |
| -------------- | -------- |
| syscall number | `RAX`    |
| arg1           | `RDI`    |
| arg2           | `RSI`    |
| arg3           | `RDX`    |
| arg4           | `R10`    |
| arg5           | `R8`     |
| arg6           | `R9`     |
| return value   | `RAX`    |

Important:

**4th argument = R10**, not RCX.

Because `SYSCALL` itself has architectural behavior involving RCX/R11.

---

## Linux example: write

Linux x86-64:

```text
SYS_write = 1
```

Conceptually:

```asm
mov rax, 1        ; syscall number
mov rdi, 1        ; fd = stdout
mov rsi, message  ; buffer
mov rdx, length   ; count
syscall
```

Kernel:

```text
user process
    │
    │ syscall
    ▼
CPU privilege transition
    │
    ▼
Linux kernel
    │
    ▼
sys_write(...)
```

---

# 3.2 Linux exit

```text
SYS_exit = 60
```

```asm
mov rax, 60
mov rdi, 0
syscall
```

---

# 3.3 `int 0x80`

32-bit Linux historically:

```asm
int 0x80
```

32-bit syscall ABI differs.

64-bit Linux normally uses:

```asm
syscall
```

So:

```text
32-bit historical:
int 0x80

64-bit:
syscall
```

`sysenter` is also a historical/architectural fast system-call mechanism, but Linux x86-64 user programs conventionally use `syscall`.

---

# 3.4 Windows system calls

Windows is fundamentally different from Linux here.

Windows application programmers generally don't write:

```asm
mov rax, syscall_number
syscall
```

because **Windows native syscall numbers are not a stable application ABI**.

Instead:

```text
Application
    │
    ▼
Win32 API
    │
    ▼
Kernel32 / KernelBase
    │
    ▼
NTDLL
    │
    ▼
Windows kernel
```

For example:

```text
WriteFile()
```

is an API-level contract.

Native NT functions such as:

```text
NtWriteFile
```

are lower-level.

The exact syscall interface is an implementation detail and syscall numbers can change between Windows versions/builds.

Therefore a professional Windows Assembly program normally uses documented APIs rather than hardcoding kernel syscall numbers.

---

# 3.5 Windows x64 ABI

Microsoft x64 ABI:

### Integer/pointer arguments

```text
1 → RCX
2 → RDX
3 → R8
4 → R9
5+ → stack
```

### Floating-point

```text
1 → XMM0
2 → XMM1
3 → XMM2
4 → XMM3
```

Microsoft documents these registers and volatility rules. ([Microsoft Learn][1])

---

# 3.6 Windows shadow space

This is one of the most important differences.

Before calling a function:

```text
caller reserves 32 bytes
```

called function can use it as register-argument home/shadow space.

```text
STACK

higher addresses
──────────────────
5th argument
4th argument
3rd argument
2nd argument
1st argument
──────────────────
32-byte shadow space
──────────────────
return address
──────────────────
lower addresses
```

Microsoft x64 ABI requires caller-provided space for four register parameters. ([Microsoft Learn][1])

---

# 3.7 Windows volatile registers

Microsoft x64:

```text
volatile:
RAX
RCX
RDX
R8
R9
R10
R11
XMM0-XMM5
```

nonvolatile:

```text
RBX
RBP
RDI
RSI
RSP
R12-R15
XMM6-XMM15
```

Microsoft's current documentation lists these preservation rules. ([Microsoft Learn][1])

---

# 3.8 System V AMD64 ABI

Linux and macOS x86-64 generally use the System V AMD64 family of conventions, although the OS ABIs are not literally identical in every detail.

Integer/pointer arguments:

```text
1 → RDI
2 → RSI
3 → RDX
4 → RCX
5 → R8
6 → R9
```

Return:

```text
integer/pointer → RAX
floating point → XMM0
```

This is dramatically different from Windows:

```text
Windows:
RCX RDX R8 R9

SysV:
RDI RSI RDX RCX R8 R9
```

---

# 3.9 The six-argument comparison

Suppose:

```c
foo(a,b,c,d,e,f);
```

### Windows x64

```text
a → RCX
b → RDX
c → R8
d → R9
e → stack
f → stack
```

### Linux/macOS SysV

```text
a → RDI
b → RSI
c → RDX
d → RCX
e → R8
f → R9
```

This single difference explains a huge amount of C/Assembly interoperability bugs.

---

# 3.10 Red Zone vs Shadow Space

Another major difference.

### System V AMD64

Has a **128-byte red zone** below RSP available to leaf functions in normal user-space code.

Conceptually:

```text
RSP
 ↓
┌─────────────────────┐
│ return address etc. │
├─────────────────────┤
│                     │
│ 128-byte red zone   │
│                     │
└─────────────────────┘
```

A leaf function can sometimes use it without adjusting RSP.

### Windows x64

No SysV-style red zone.

Instead Windows uses:

```text
32-byte shadow/home space
```

This is a very important cross-platform distinction.

---

# 3.11 Stack alignment

Both ABIs care about stack alignment.

For System V AMD64, at a normal function call boundary, the ABI requires the stack to be arranged so the callee sees the expected alignment, commonly described as:

```text
RSP + 8 ≡ 0 (mod 16)
```

at function entry.

Windows x64 requires stack alignment to support its calling convention and SIMD usage, with exceptions around leaf/prologue/epilogue situations. Microsoft explicitly documents 16-byte alignment requirements for non-leaf regions. ([Microsoft Learn][1])

---

# 3.12 Full ABI comparison

| Feature            | Windows x64                | Linux x86-64      | macOS x86-64            |
| ------------------ | -------------------------- | ----------------- | ----------------------- |
| Common ABI family  | Microsoft x64              | System V AMD64    | SysV-derived Darwin ABI |
| arg1 integer       | RCX                        | RDI               | RDI                     |
| arg2               | RDX                        | RSI               | RSI                     |
| arg3               | R8                         | RDX               | RDX                     |
| arg4               | R9                         | RCX               | RCX                     |
| arg5               | stack                      | R8                | R8                      |
| arg6               | stack                      | R9                | R9                      |
| integer return     | RAX                        | RAX               | RAX                     |
| FP return          | XMM0                       | XMM0              | XMM0                    |
| shadow space       | 32 B                       | No                | No                      |
| red zone           | No                         | 128 B             | generally 128 B         |
| executable         | PE                         | ELF               | Mach-O                  |
| assembler commonly | MASM/NASM                  | GAS/NASM          | Clang/GAS/NASM          |
| kernel interface   | NT system services / Win32 | Linux syscall ABI | Darwin syscall ABI      |

---

# 3.13 Memory Management

මෙතන OS එක CPU architecture එකට එකතු කරන major functionality එකක් තමයි **virtual memory**.

Application එකට:

```text
0x0000000000400000
```

වගේ virtual address එකක් පේන්න පුළුවන්.

ඒක direct physical RAM address එකක් නෙවෙයි.

```text
CPU virtual address
       │
       ▼
     MMU
       │
       ▼
 Page Tables
       │
       ▼
Physical address
       │
       ▼
    RAM
```

---

# 3.14 Pages

Typical x86-64 page size:

```text
4 KiB
```

Virtual memory:

```text
Virtual Address Space

0x0000...
┌──────────────┐
│ code         │
├──────────────┤
│ read-only    │
├──────────────┤
│ data         │
├──────────────┤
│ heap ↑       │
│              │
│              │
│ stack ↓      │
├──────────────┤
│ kernel       │
└──────────────┘
```

Exact layout, addresses, ASLR behavior, and user/kernel split are OS/version/configuration dependent.

---

# 3.15 User mode vs Kernel mode

CPU privilege levels x86-64 වල:

```text
Ring 3 → User mode
Ring 0 → Kernel mode
```

Typical:

```text
Application
Ring 3
   │
   │ syscall
   ▼
Kernel
Ring 0
```

Application එකට arbitrary kernel memory access කරන්න බැහැ.

මේ isolation එක OS security architecture එකේ foundation එකක්.

---

# 3.16 Page permissions

Memory page එකට permissions තියෙන්න පුළුවන්:

```text
R = Read
W = Write
X = Execute
```

Typical code:

```text
.text → R-X
```

Data:

```text
.data → R-W
```

Modern systems generally enforce W^X / NX-related protections to reduce executable writable memory.

---

# 3.17 Windows memory model

Windows process:

```text
User Virtual Address Space
│
├── EXE image
├── DLLs
├── heap
├── stacks
├── mapped files
└── shared regions
```

Kernel virtual address space is isolated from normal user-mode access.

PE loader maps executable sections and resolves imports/relocations as necessary.

---

# 3.18 Linux memory model

Linux process commonly contains:

```text
┌─────────────────────────┐
│ stack                   │
├─────────────────────────┤
│ mmap regions / DLLs     │
├─────────────────────────┤
│ heap                    │
├─────────────────────────┤
│ BSS                     │
├─────────────────────────┤
│ data                    │
├─────────────────────────┤
│ rodata                  │
├─────────────────────────┤
│ text                    │
└─────────────────────────┘
```

ELF loader maps loadable segments.

ASLR randomizes relevant mappings.

---

# 3.19 macOS memory model

macOS similarly uses per-process virtual address spaces with:

```text
Mach-O image
dyld
shared libraries
heap
stack
mmap regions
kernel mappings
```

But the executable/loader infrastructure is Darwin/Mach-O based.

---

# 4. REAL-WORLD CODE MATRIX

දැන් එකම:

```text
Hello, World!
```

OS තුනේ.

මම **x86-64** use කරනවා. ඒක Windows/Linux/macOS තුනම direct compare කරන්න පහසු නිසා.

---

# 4.1 Linux x86-64 — pure syscall

NASM source:

```asm
section .data

message db "Hello, World!", 10
message_len equ $ - message


section .text

global _start

_start:
    mov rax, 1
    mov rdi, 1
    mov rsi, message
    mov rdx, message_len
    syscall

    mov rax, 60
    xor rdi, rdi
    syscall
```

---

## Line-by-line

### Data section

```asm
section .data
```

initialized writable data.

---

```asm
message db "Hello, World!", 10
```

`message` = label.

`db` = Define Byte.

Characters bytes ලෙස store කරනවා.

```text
H e l l o ,   W o r l d !
```

අවසාන:

```text
10
```

newline (`\n`).

---

```asm
message_len equ $ - message
```

`$` = current assembly location.

So:

```text
$ - message
```

= message length.

`equ` = assembly-time constant.

Runtime instruction එකක් නෙවෙයි.

---

# 4.2 `_start`

```asm
global _start
```

Linker එකට `_start` symbol එක externally visible කරන්න.

Linux ELF executable එකේ entry point එක `_start` ලෙස set කරන්න පුළුවන්.

---

```asm
_start:
```

program entry label.

CRT `main()` අවශ්‍ය නැහැ.

---

# 4.3 Linux write

```asm
mov rax, 1
```

Linux x86-64:

```text
RAX = syscall number
```

`1 = write`.

---

```asm
mov rdi, 1
```

first syscall argument:

```text
RDI = file descriptor
```

`1` = stdout.

---

```asm
mov rsi, message
```

second argument:

```text
RSI = buffer address
```

---

```asm
mov rdx, message_len
```

third:

```text
RDX = number of bytes
```

---

```asm
syscall
```

CPU transitions through the Linux syscall mechanism into kernel mode.

Conceptually:

```text
RAX = 1
RDI = stdout
RSI = buffer
RDX = length

        syscall
           │
           ▼
       Linux kernel
           │
           ▼
        write()
```

Return value comes back in `RAX`.

---

# 4.4 Linux exit

```asm
mov rax, 60
```

Linux x86-64:

```text
60 = exit
```

```asm
xor rdi, rdi
```

```text
RDI = 0
```

exit status 0.

`xor rdi,rdi` common zeroing idiom.

---

```asm
syscall
```

Process terminates.

---

# 4.5 Linux build

```bash
nasm -f elf64 hello.asm -o hello.o
ld hello.o -o hello
./hello
```

Pipeline:

```text
hello.asm
   │
 NASM
   ↓
hello.o
   │
  ld
   ↓
ELF executable
   │
 Linux loader
   ↓
 process
```

---

# 4.6 Windows x64

Windows application-level Assembly should generally use the **documented Windows API**, rather than hard-coded NT syscall numbers.

We'll use:

```text
GetStdHandle
WriteFile
ExitProcess
```

NASM:

```asm
default rel

extern GetStdHandle
extern WriteFile
extern ExitProcess

global main

section .data

message db "Hello, World!", 13, 10
message_len equ $ - message


section .bss

bytes_written resq 1


section .text

main:

    sub rsp, 40

    mov ecx, -11
    call GetStdHandle

    mov rcx, rax
    mov rdx, message
    mov r8d, message_len
    lea r9, [rel bytes_written]

    mov qword [rsp+32], 0

    call WriteFile

    xor ecx, ecx
    call ExitProcess
```

---

# 4.7 Why Windows code looks completely different

Notice:

Linux:

```asm
mov rax, 1
mov rdi, 1
mov rsi, message
mov rdx, message_len
syscall
```

Windows:

```asm
call GetStdHandle
call WriteFile
call ExitProcess
```

Reason:

```text
Linux application
       │
       ▼
syscall ABI
       │
       ▼
kernel

Windows application
       │
       ▼
Win32 API
       │
       ▼
Windows user-mode libraries
       │
       ▼
NT kernel interface
```

---

# 4.8 Windows `GetStdHandle`

```asm
mov ecx, -11
call GetStdHandle
```

Windows x64:

```text
first argument → RCX
```

`STD_OUTPUT_HANDLE` is `-11`.

Because `ECX` is used, the value is placed in the low 32 bits and zero-extension semantics apply to R64 register state. The API's parameter is a signed integer constant.

Return:

```text
RAX = HANDLE
```

---

# 4.9 Windows `WriteFile`

C prototype conceptually:

```c
BOOL WriteFile(
    HANDLE       hFile,
    LPCVOID      lpBuffer,
    DWORD        nNumberOfBytesToWrite,
    LPDWORD      lpNumberOfBytesWritten,
    LPOVERLAPPED lpOverlapped
);
```

Five arguments.

Windows x64:

```text
arg1 → RCX
arg2 → RDX
arg3 → R8
arg4 → R9
arg5 → stack
```

Therefore:

```asm
mov rcx, rax
```

handle.

```asm
mov rdx, message
```

buffer.

```asm
mov r8d, message_len
```

length.

```asm
lea r9, [rel bytes_written]
```

address where Windows can store number of bytes written.

---

# 4.10 The fifth argument

```asm
mov qword [rsp+32], 0
```

This is the fifth argument location in our call setup.

But notice the Windows ABI requirement:

```text
32 bytes shadow space
```

We reserve:

```asm
sub rsp, 40
```

Why 40?

```text
32 bytes shadow space
+
8 bytes alignment adjustment
=
40
```

At function entry from the Microsoft CRT, the return-address state means subtracting 40 gives the alignment needed before calls in this simple function.

---

# 4.11 Why `default rel`?

```asm
default rel
```

NASM directive that makes suitable memory references use RIP-relative addressing by default in 64-bit mode.

This is useful for modern x86-64 position-friendly code.

---

# 4.12 Why `ExitProcess`

```asm
xor ecx, ecx
call ExitProcess
```

Windows API:

```text
ExitProcess(0)
```

First argument:

```text
RCX = 0
```

Unlike Linux:

```asm
syscall
```

Windows application code normally doesn't directly invoke a stable public syscall-number ABI.

---

# 4.13 Windows linking

Using MinGW-w64 toolchain, conceptually:

```bash
nasm -f win64 hello.asm -o hello.obj
gcc hello.obj -o hello.exe -lkernel32
```

Here the linker/CRT environment supplies the normal PE executable startup machinery around `main`.

So this example is:

```text
NASM
 ↓
COFF object
 ↓
MinGW linker/CRT
 ↓
PE executable
 ↓
Windows loader
 ↓
main
```

It is intentionally an **API-level Windows example**, because that is the portable professional model for Windows user-space Assembly.

---

# 4.14 macOS x86-64 — direct Darwin syscall

macOS x86-64 has a Darwin syscall ABI.

A classic direct-syscall example:

```asm
global _main

section .data

message db "Hello, World!", 10
message_len equ $ - message


section .text

_main:

    mov rax, 0x2000004
    mov rdi, 1
    lea rsi, [rel message]
    mov rdx, message_len
    syscall

    mov rax, 0x2000001
    xor rdi, rdi
    syscall
```

---

# 4.15 macOS symbol naming

Notice:

```asm
global _main
```

not:

```asm
global main
```

Mach-O/Darwin C symbol naming traditionally uses a leading underscore.

So:

```text
C:
main

Mach-O symbol:
_main
```

This is one of the visible differences between Linux ELF and macOS Mach-O toolchains.

---

# 4.16 macOS syscall number

Classic x86-64 Darwin syscall encoding:

```text
0x20000000
```

indicates the UNIX syscall class, with the syscall number encoded in the low portion.

For example:

```text
0x2000004
```

corresponds to:

```text
write
```

and:

```text
0x2000001
```

corresponds to:

```text
exit
```

**Important:** direct Darwin syscall programming is much less desirable for normal application development than using the documented libc/system APIs, and syscall interfaces should not be treated as a stable cross-version application ABI.

---

# 4.17 macOS build

With NASM and Apple's linker environment:

```bash
nasm -f macho64 hello.asm -o hello.o
ld hello.o -o hello -e _main -lSystem -syslibroot $(xcrun --sdk-path macosx)
./hello
```

Depending on the installed SDK/toolchain, exact linker flags can vary. For a simple direct-syscall executable, you can also deliberately avoid libc and use an appropriate Mach-O entry-point configuration.

The important architecture is:

```text
hello.asm
    │
   NASM
    ↓
Mach-O object
    │
    ↓
Apple linker
    ↓
Mach-O executable
    │
    ↓
macOS loader
    ↓
_main
```

---

# 4.18 Three Hello Worlds side-by-side

## Linux

```asm
mov rax, 1
mov rdi, 1
mov rsi, message
mov rdx, message_len
syscall
```

## macOS

```asm
mov rax, 0x2000004
mov rdi, 1
mov rsi, message
mov rdx, message_len
syscall
```

## Windows

```asm
mov ecx, -11
call GetStdHandle

mov rcx, rax
mov rdx, message
mov r8d, message_len
lea r9, [rel bytes_written]

call WriteFile
```

Notice something extremely important:

### CPU instruction

All three use:

```asm
mov
```

and potentially:

```asm
call
```

or:

```asm
syscall
```

because CPU architecture is the same.

But the **contract surrounding those instructions** is different.

---

# 4.19 One function, three ABIs

Suppose:

```c
long calculate(
    long a,
    long b,
    long c,
    long d,
    long e,
    long f
);
```

### Windows

```text
RCX = a
RDX = b
R8  = c
R9  = d

[stack] = e
[stack] = f
```

### Linux

```text
RDI = a
RSI = b
RDX = c
RCX = d
R8  = e
R9  = f
```

### macOS x86-64

```text
RDI = a
RSI = b
RDX = c
RCX = d
R8  = e
R9  = f
```

Therefore the same machine code cannot simply be copied between OSes and expected to call the same C function correctly.

---

# 4.20 Why Assembly itself isn't really "Windows Assembly"

මේක ඉතා වැදගත් conceptual correction එකක්.

Technically:

```text
x86-64 Assembly
```

is the language/instruction representation of the CPU ISA.

But:

```text
Windows Assembly program
Linux Assembly program
macOS Assembly program
```

කියනකොට usually අදහස් කරන්නේ:

```text
x86-64 ISA
+
assembler syntax
+
ABI
+
object format
+
linker
+
loader
+
OS API/syscall ABI
```

---

# 4.21 Complete mental model

```text
                         SOURCE
                           │
                           ▼
                 ┌──────────────────┐
                 │ Assembly syntax  │
                 │ NASM/MASM/GAS    │
                 └────────┬─────────┘
                          │
                          ▼
                 ┌──────────────────┐
                 │    Assembler     │
                 └────────┬─────────┘
                          │
                          ▼
                    Object File
                          │
          ┌───────────────┼────────────────┐
          │               │                │
          ▼               ▼                ▼
         PE              ELF             Mach-O
       Windows           Linux            macOS
          │               │                │
          └───────────────┼────────────────┘
                          ▼
                       Linker
                          │
                          ▼
                     Executable
                          │
                          ▼
                        Loader
                          │
                          ▼
                       Process
                          │
                 ┌────────┴────────┐
                 │                 │
              User Mode        Kernel Mode
                 │                 │
                 │ syscall/API     │
                 └────────►────────┘
                          │
                          ▼
                       Hardware
```

---

# 4.22 The deepest distinction: ISA vs ABI vs OS API

මේ තුන හොඳට වෙන් කරගන්න.

### ① ISA — Instruction Set Architecture

CPU එක තේරුම් ගන්න instruction vocabulary:

```text
MOV
ADD
SUB
MUL
DIV
CMP
JMP
CALL
RET
PUSH
POP
SYSCALL
...
```

x86-64 CPU එකට අදාල.

---

### ② ABI — Application Binary Interface

Programs/functions එකිනෙකා සමඟ binary level එකෙන් communicate කරන rules.

Includes:

```text
argument registers
return registers
stack layout
stack alignment
volatile registers
nonvolatile registers
structure layout
calling convention
symbol conventions
```

---

### ③ OS API / syscall ABI

Application එක OS එකෙන් service request කරන contract.

Linux:

```text
syscall
```

Windows:

```text
Win32 / NT interfaces
```

macOS:

```text
Darwin/POSIX interfaces
```

---

# 4.23 Compiler එක මේ සියල්ල combine කරන හැටි

C:

```c
int add(int a, int b)
{
    return a + b;
}
```

Compiler knows:

```text
target architecture = x86-64
target OS = Windows/Linux/macOS
target ABI = corresponding ABI
```

Then Windows version may conceptually become:

```asm
add:
    lea eax, [rcx + rdx]
    ret
```

Linux/macOS version:

```asm
add:
    lea eax, [rdi + rsi]
    ret
```

Notice:

**function body logic එක same.**

නමුත් arguments arrive කරන registers වෙනස්.

---

# 4.24 Why compiler matters so much

ඔබ:

```c
printf("Hello");
```

ලියනකොට compiler එකෙන් generated machine code එක:

```text
C source
  ↓
Compiler
  ↓
ABI-compliant function call
  ↓
Library
  ↓
OS interface
  ↓
Kernel
```

Assembly developer කෙනෙක්ට මේ entire chain එක understand කරන්න පුළුවන්.

ඒක තමයි **systems programming**.

---

# 4.25 Final comparison matrix

| Layer                 | Windows x64                       | Linux x86-64                    | macOS x86-64                  |
| --------------------- | --------------------------------- | ------------------------------- | ----------------------------- |
| CPU ISA               | x86-64                            | x86-64                          | x86-64                        |
| Registers             | Same                              | Same                            | Same                          |
| `MOV` instruction     | Same                              | Same                            | Same                          |
| `CMP`/`Jcc`           | Same                              | Same                            | Same                          |
| `CALL`/`RET`          | Same                              | Same                            | Same                          |
| Common syntax         | Intel/MASM/NASM                   | AT&T/GAS/NASM                   | AT&T/Intel/NASM               |
| ABI                   | Microsoft x64                     | SysV AMD64                      | Darwin/SysV-derived           |
| arg1                  | RCX                               | RDI                             | RDI                           |
| arg2                  | RDX                               | RSI                             | RSI                           |
| arg3                  | R8                                | RDX                             | RDX                           |
| arg4                  | R9                                | RCX                             | RCX                           |
| arg5                  | stack                             | R8                              | R8                            |
| arg6                  | stack                             | R9                              | R9                            |
| Shadow space          | 32 B                              | No                              | No                            |
| Red zone              | No                                | 128 B                           | 128 B generally               |
| Object/executable     | PE/COFF                           | ELF                             | Mach-O                        |
| C symbol example      | toolchain-dependent               | `main`                          | `_main`                       |
| Direct syscall ABI    | not stable/public app ABI         | `syscall`                       | Darwin syscall mechanism      |
| syscall #1 equivalent | API, not stable syscall number    | `write=1`                       | `write=0x2000004`             |
| exit                  | `ExitProcess` API                 | `exit=60`                       | `exit=0x2000001`              |
| Dynamic loader        | Windows loader                    | ELF dynamic loader              | `dyld`                        |
| Kernel                | NT kernel                         | Linux kernel                    | XNU/Darwin kernel             |
| VM                    | hardware page tables + Windows VM | hardware page tables + Linux VM | hardware page tables + XNU VM |

---

# 4.26 The one diagram I want you to remember

```text
                    ┌──────────────────────┐
                    │      CPU ISA         │
                    │      x86-64          │
                    │                      │
                    │ MOV ADD CMP JMP CALL │
                    │ RAX RBX RCX RDX ...  │
                    └──────────┬───────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
          WINDOWS           LINUX            macOS
              │                │                │
          Microsoft          SysV             Darwin
             ABI              ABI               ABI
              │                │                │
            PE/COFF           ELF             Mach-O
              │                │                │
          Win32/NT          Linux syscall    Darwin API
              │                │                │
              ▼                ▼                ▼
           Windows           Linux            XNU
           Kernel            Kernel           Kernel
              │                │                │
              └────────────────┼────────────────┘
                               ▼
                         Hardware / MMU
                               │
                               ▼
                              RAM
```

## 🔥 Principal-level takeaway

Assembly mastery කියන්නේ **mnemonics memorize කිරීම නෙවෙයි**.

ඔබට මේ chain එක mental model එකක් විදිහට තිබිය යුතුයි:

```text
High-level language
        ↓
Compiler
        ↓
ABI
        ↓
Assembly
        ↓
Machine-code encoding
        ↓
CPU decoder
        ↓
Registers / ALU / flags
        ↓
Virtual address
        ↓
MMU + page tables
        ↓
Physical memory
```

ඊට පස්සේ OS boundary එක:

```text
Application
    ↓
ABI
    ↓
Library / API
    ↓
Syscall mechanism
    ↓
CPU privilege transition
    ↓
Kernel
    ↓
Driver
    ↓
Hardware
```

**මෙතනින් ඔබට එක fundamental truth එකක් පැහැදිලි වෙන්න ඕනේ:**

> **Windows, Linux, macOS තුනේ CPU එක එකම x86-64 instruction set එක execute කරනවා. වෙනස් වෙන්නේ CPU එකේ `MOV`, `ADD`, `JMP` වගේ instructions නොවෙයි; ඒ instructions වලින් build කරන software contract එක — ABI, object format, linker/loader rules, system interface, runtime libraries — තමයි.**

ඒ නිසා **NASM ඉගෙනීම = CPU ඉගෙනීමේ එක් කොටසක්**.
**Systems Assembly engineer වීමට = ISA + ABI + executable formats + linker + loader + virtual memory + OS kernel interface** සියල්ල එකට තේරුම් ගත යුතුයි. 💻⚙️

Microsoft x64 register/ABI details සඳහා Microsoftගේ current documentation ද බලන්න පුළුවන්. ([Microsoft Learn][1])

[1]: https://learn.microsoft.com/en-us/cpp/build/x64-calling-convention?view=msvc-170&utm_source=chatgpt.com "x64 Calling Convention | Microsoft Learn"
