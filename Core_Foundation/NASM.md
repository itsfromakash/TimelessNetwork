**“NASM + x86-64 registers + C/C++ datatypes + Windows/Linux/macOS ABI + calling convention + stack + memory + instruction operand sizes + data definitions + SIMD registers + flags + caller/callee saved + object formats”** 

එක වැදගත් correction එකක් මුලින්ම:

> **NASM එකට `int`, `char`, `float`, `double` වගේ C/C++ datatype system එකක් නැහැ.**
> NASM assembly එකේ අපි mostly **byte/word/dword/qword**, memory operands, registers, instructions, symbols, labels වගේ low-level concepts භාවිතා කරනවා.
> C/C++ datatype එක NASM එකට translate වෙන්නේ **size + representation + ABI classification** ලෙස.

පහත table එක **x86-64 NASM master reference** එකක් ලෙස use කරන්න. 🧠

---

# 1. 🧠 C/C++ → Compiler → ABI → NASM Mental Model

```text
C / C++
   │
   │  int a, double b, char c...
   ▼
Compiler
   │
   │  ABI rules apply
   ▼
x86-64 machine-level representation
   │
   ▼
Registers / Stack / Memory
   │
   ▼
NASM Assembly
   │
   ▼
Assembler
   │
   ▼
.obj / .o
   │
   ▼
Linker
   │
   ▼
Executable
```

**ABI** කියන්නේ වෙනම assembly instruction එකක් හෝ compiler stage එකක් නෙවෙයි.

ABI කියන්නේ:

```text
"Function එකක් call කරනකොට
arguments කොහෙද තියෙන්නේ?
return value කොහෙද?
කොයි registers preserve කරන්න ඕනද?
stack එක කොහොමද?
alignment කොහොමද?
memory representation කොහොමද?"
```

වගේ binary-level contract එක.

---

# 2. 📦 NASM Fundamental Data Sizes

| Name    |     Size | Bits | Common use             |
| ------- | -------: | ---: | ---------------------- |
| `BYTE`  |   1 byte |    8 | `char`, byte data      |
| `WORD`  |  2 bytes |   16 | `short`, 16-bit value  |
| `DWORD` |  4 bytes |   32 | `int`, 32-bit value    |
| `QWORD` |  8 bytes |   64 | `long long`, pointer   |
| `TWORD` | 10 bytes |   80 | x87 extended precision |
| `OWORD` | 16 bytes |  128 | XMM-sized memory       |
| `YWORD` | 32 bytes |  256 | YMM-sized memory       |
| `ZWORD` | 64 bytes |  512 | ZMM-sized memory       |

### Size relationship

```text
BYTE
  8 bits
    ↓
WORD
 16 bits
    ↓
DWORD
 32 bits
    ↓
QWORD
 64 bits
    ↓
OWORD
128 bits
    ↓
YWORD
256 bits
    ↓
ZWORD
512 bits
```

---

# 3. 🔥 NASM Register Master Table

| Family            | 64-bit | 32-bit | 16-bit | 8-bit low |
| ----------------- | ------ | ------ | ------ | --------- |
| A                 | `RAX`  | `EAX`  | `AX`   | `AL`      |
| B                 | `RBX`  | `EBX`  | `BX`   | `BL`      |
| C                 | `RCX`  | `ECX`  | `CX`   | `CL`      |
| D                 | `RDX`  | `EDX`  | `DX`   | `DL`      |
| Source Index      | `RSI`  | `ESI`  | `SI`   | `SIL`     |
| Destination Index | `RDI`  | `EDI`  | `DI`   | `DIL`     |
| R8                | `R8`   | `R8D`  | `R8W`  | `R8B`     |
| R9                | `R9`   | `R9D`  | `R9W`  | `R9B`     |
| R10               | `R10`  | `R10D` | `R10W` | `R10B`    |
| R11               | `R11`  | `R11D` | `R11W` | `R11B`    |
| R12               | `R12`  | `R12D` | `R12W` | `R12B`    |
| R13               | `R13`  | `R13D` | `R13W` | `R13B`    |
| R14               | `R14`  | `R14D` | `R14W` | `R14B`    |
| R15               | `R15`  | `R15D` | `R15W` | `R15B`    |

---

# 4. ⚠️ Special Register Width Rules

`RAX` family:

```text
RAX
┌──────────────────────────────────────────────┐
│                  64 bits                     │
└──────────────────────────────────────────────┘

                    EAX
                    ┌───────────────────────────┐
                    │       32 bits             │
                    └───────────────────────────┘

                            AX
                            ┌──────────────────┐
                            │     16 bits      │
                            └──────────────────┘

                                  AL
                                  ┌────────────┐
                                  │  8 bits    │
                                  └────────────┘
```

For example:

```asm
mov eax, 10
```

means:

```text
EAX = 10
```

and importantly in x86-64:

```text
RAX upper 32 bits → zeroed
```

So:

```asm
mov rax, 0xFFFFFFFFFFFFFFFF
mov eax, 10
```

results in:

```text
RAX = 0x000000000000000A
```

But:

```asm
mov rax, 0xFFFFFFFFFFFFFFFF
mov ax, 10
```

does **not** zero the upper bits.

Result:

```text
RAX = 0xFFFFFFFFFFFF000A
```

---

# 5. 🪟🐧🍎 Integer / Pointer ABI Register Table

## x86-64

| Argument            | Windows x64 | Linux SysV AMD64 | macOS Intel x86-64 |
| ------------------- | ----------- | ---------------- | ------------------ |
| 1st integer/pointer | `RCX`       | `RDI`            | `RDI`              |
| 2nd                 | `RDX`       | `RSI`            | `RSI`              |
| 3rd                 | `R8`        | `RDX`            | `RDX`              |
| 4th                 | `R9`        | `RCX`            | `RCX`              |
| 5th                 | Stack       | `R8`             | `R8`               |
| 6th                 | Stack       | `R9`             | `R9`               |
| 7th+                | Stack       | Stack            | Stack              |

### 32-bit views

| Argument | Windows | Linux | macOS |
| -------- | ------- | ----- | ----- |
| 1st      | `ECX`   | `EDI` | `EDI` |
| 2nd      | `EDX`   | `ESI` | `ESI` |
| 3rd      | `R8D`   | `EDX` | `EDX` |
| 4th      | `R9D`   | `ECX` | `ECX` |
| 5th      | Stack   | `R8D` | `R8D` |
| 6th      | Stack   | `R9D` | `R9D` |

---

# 6. 🧩 ABI Mapping + Register Width Together

| Parameter | Windows 64 | Windows 32 | Windows 16 | Windows 8 | Linux 64 | Linux 32 | Linux 16 | Linux 8 | macOS 64 | macOS 32 | macOS 16 | macOS 8 |
| --------- | ---------- | ---------- | ---------- | --------- | -------- | -------- | -------- | ------- | -------- | -------- | -------- | ------- |
| 1st       | `RCX`      | `ECX`      | `CX`       | `CL`      | `RDI`    | `EDI`    | `DI`     | `DIL`   | `RDI`    | `EDI`    | `DI`     | `DIL`   |
| 2nd       | `RDX`      | `EDX`      | `DX`       | `DL`      | `RSI`    | `ESI`    | `SI`     | `SIL`   | `RSI`    | `ESI`    | `SI`     | `SIL`   |
| 3rd       | `R8`       | `R8D`      | `R8W`      | `R8B`     | `RDX`    | `EDX`    | `DX`     | `DL`    | `RDX`    | `EDX`    | `DX`     | `DL`    |
| 4th       | `R9`       | `R9D`      | `R9W`      | `R9B`     | `RCX`    | `ECX`    | `CX`     | `CL`    | `RCX`    | `ECX`    | `CX`     | `CL`    |
| 5th       | Stack      | Stack      | Stack      | Stack     | `R8`     | `R8D`    | `R8W`    | `R8B`   | `R8`     | `R8D`    | `R8W`    | `R8B`   |
| 6th       | Stack      | Stack      | Stack      | Stack     | `R9`     | `R9D`    | `R9W`    | `R9B`   | `R9`     | `R9D`    | `R9W`    | `R9B`   |
| 7th+      | Stack      | Stack      | Stack      | Stack     | Stack    | Stack    | Stack    | Stack   | Stack    | Stack    | Stack    | Stack   |

> ⚠️ **මෙය simple integer/pointer arguments සඳහා.** Structs, unions, vectors, floating-point, variadic functions, `_Bool`, complex aggregates වගේ cases වල ABI classification rules තවත් complex.

---

# 7. 🧱 C/C++ Type → Typical Size → NASM Representation

| C/C++ Type           | Typical size on x86-64 | Bits | NASM-oriented representation          |
| -------------------- | ---------------------: | ---: | ------------------------------------- |
| `char`               |                 1 byte |    8 | `BYTE`                                |
| `signed char`        |                      1 |    8 | `BYTE`                                |
| `unsigned char`      |                      1 |    8 | `BYTE`                                |
| `short`              |                      2 |   16 | `WORD`                                |
| `unsigned short`     |                      2 |   16 | `WORD`                                |
| `int`                |                      4 |   32 | `DWORD`                               |
| `unsigned int`       |                      4 |   32 | `DWORD`                               |
| `long long`          |                      8 |   64 | `QWORD`                               |
| `unsigned long long` |                      8 |   64 | `QWORD`                               |
| pointer              |                      8 |   64 | `QWORD`                               |
| `float`              |                      4 |   32 | `DWORD` memory representation / `XMM` |
| `double`             |                      8 |   64 | `QWORD` memory representation / `XMM` |
| `long double`        |          ABI-dependent |    — | ABI-specific                          |
| `bool` / `_Bool`     |             commonly 1 |    8 | `BYTE`                                |
| `size_t`             |                      8 |   64 | `QWORD`                               |
| `intptr_t`           |                      8 |   64 | `QWORD`                               |
| `uintptr_t`          |                      8 |   64 | `QWORD`                               |

⚠️ `long` කියන්නේ **Windows x64 සහ Linux/macOS x86-64 අතර වෙනස් වෙන්න පුළුවන්**.

| Platform     | C/C++ `long` |
| ------------ | -----------: |
| Windows x64  |       32-bit |
| Linux x86-64 |       64-bit |
| macOS x86-64 |       64-bit |

මේක ඉතා වැදගත්.

---

# 8. 🔢 NASM Data Declaration Directives

| NASM directive | Size per element | Example     |
| -------------- | ---------------: | ----------- |
| `DB`           |           1 byte | `db 10`     |
| `DW`           |          2 bytes | `dw 1000`   |
| `DD`           |          4 bytes | `dd 100000` |
| `DQ`           |          8 bytes | `dq 100000` |
| `DT`           |         10 bytes | `dt 1.0`    |
| `DO`           |         16 bytes | `do ...`    |
| `DY`           |         32 bytes | `dy ...`    |
| `DZ`           |         64 bytes | `dz ...`    |

Example:

```asm
section .data

myByte  db  10
myWord  dw  1000
myDword dd  100000
myQword dq  10000000000
```

---

# 9. 🧮 Integer Constants

| Format      | Example     | Meaning         |
| ----------- | ----------- | --------------- |
| Decimal     | `123`       | decimal         |
| Hexadecimal | `0x7B`      | hex             |
| Binary      | `0b1111011` | binary          |
| Octal       | `0o173`     | octal           |
| Character   | `'A'`       | character value |
| String      | `"Hello"`   | bytes           |

Example:

```asm
mov eax, 100
mov eax, 0x64
mov eax, 0b1100100
```

all represent `100`.

---

# 10. 🧮 Signed vs Unsigned

**CPU register එකට “signed” / “unsigned” කියලා separate mode එකක් නැහැ.**

Example:

```text
11111111
```

8-bit value එක:

| Interpretation          | Value |
| ----------------------- | ----: |
| Unsigned                | `255` |
| Signed two's complement |  `-1` |

Assembly instruction එක සහ comparison instruction එක අනුව interpretation එක වෙනස් වෙනවා.

| Purpose             | Instructions             |
| ------------------- | ------------------------ |
| Signed comparison   | `JG`, `JL`, `JGE`, `JLE` |
| Unsigned comparison | `JA`, `JB`, `JAE`, `JBE` |
| Equal               | `JE` / `JZ`              |
| Not equal           | `JNE` / `JNZ`            |

---

# 11. 🧮 Multiplication / Division Register Rules

| Operation    | Important registers      |
| ------------ | ------------------------ |
| `mul` 8-bit  | `AL` → result `AX`       |
| `mul` 16-bit | `AX` → result `DX:AX`    |
| `mul` 32-bit | `EAX` → result `EDX:EAX` |
| `mul` 64-bit | `RAX` → result `RDX:RAX` |
| `div` 8-bit  | `AX` / divisor           |
| `div` 16-bit | `DX:AX` / divisor        |
| `div` 32-bit | `EDX:EAX` / divisor      |
| `div` 64-bit | `RDX:RAX` / divisor      |

Example:

```asm
mov rax, 10
mov rbx, 20
mul rbx
```

64-bit result:

```text
RDX:RAX
```

---

# 12. 🧠 Main General-Purpose Register Roles

මේවා **historical / conventional roles**. Modern x86-64 instructions වල registers බොහෝ විට general-purpose ලෙස භාවිතා කළ හැක.

| Register | Traditional role           |
| -------- | -------------------------- |
| `RAX`    | accumulator / return value |
| `RBX`    | base register              |
| `RCX`    | counter / argument         |
| `RDX`    | data / argument            |
| `RSI`    | source index               |
| `RDI`    | destination index          |
| `RSP`    | stack pointer              |
| `RBP`    | frame/base pointer         |
| `R8-R15` | additional GPRs            |

---

# 13. 🏠 Special GPRs

| Register | Size | Purpose             |
| -------- | ---: | ------------------- |
| `RSP`    |   64 | Stack Pointer       |
| `RBP`    |   64 | Frame/Base Pointer  |
| `RIP`    |   64 | Instruction Pointer |
| `RFLAGS` |   64 | CPU flags           |

⚠️ `RIP` normal `mov rax, rip` වගේ register operand එකක් ලෙස සාමාන්‍යයෙන් භාවිතා කරන්න බැහැ. RIP-relative addressing තිබෙනවා.

Example:

```asm
mov rax, [rel message]
```

---

# 14. 🚩 RFLAGS Important Bits

| Flag | Meaning          |
| ---- | ---------------- |
| `CF` | Carry Flag       |
| `PF` | Parity Flag      |
| `AF` | Auxiliary Carry  |
| `ZF` | Zero Flag        |
| `SF` | Sign Flag        |
| `OF` | Overflow Flag    |
| `DF` | Direction Flag   |
| `IF` | Interrupt Enable |

Example:

```asm
cmp eax, ebx
je equal
```

`cmp` internally subtraction-like operation එකක් කරලා flags set කරනවා.

`JE` → `ZF = 1`.

---

# 15. 🧮 Common Integer Instructions

| Instruction | Meaning                     |
| ----------- | --------------------------- |
| `MOV`       | copy                        |
| `LEA`       | calculate effective address |
| `ADD`       | addition                    |
| `SUB`       | subtraction                 |
| `INC`       | increment                   |
| `DEC`       | decrement                   |
| `IMUL`      | signed multiplication       |
| `MUL`       | unsigned multiplication     |
| `IDIV`      | signed division             |
| `DIV`       | unsigned division           |
| `NEG`       | negate                      |
| `CMP`       | compare                     |
| `TEST`      | bitwise AND for flags       |
| `AND`       | bitwise AND                 |
| `OR`        | bitwise OR                  |
| `XOR`       | bitwise XOR                 |
| `NOT`       | bitwise NOT                 |
| `SHL`       | shift left                  |
| `SHR`       | logical shift right         |
| `SAR`       | arithmetic shift right      |
| `ROL`       | rotate left                 |
| `ROR`       | rotate right                |

---

# 16. 📏 Operand Size Prefixes

NASM code එකේ instruction එකට operands වල size එකෙන් CPU එකට operation width එක තේරෙන අවස්ථා ගොඩක් තියෙනවා.

```asm
mov al, 10       ; 8-bit
mov ax, 10       ; 16-bit
mov eax, 10      ; 32-bit
mov rax, 10      ; 64-bit
```

| Operand | Width |
| ------- | ----: |
| `AL`    |     8 |
| `AX`    |    16 |
| `EAX`   |    32 |
| `RAX`   |    64 |

---

# 17. 🧠 Memory Operand Sizes

Memory එකට:

```asm
mov byte  [rax], 10
mov word  [rax], 10
mov dword [rax], 10
mov qword [rax], 10
```

| NASM keyword | Bits |
| ------------ | ---: |
| `byte`       |    8 |
| `word`       |   16 |
| `dword`      |   32 |
| `qword`      |   64 |
| `oword`      |  128 |
| `yword`      |  256 |
| `zword`      |  512 |

---

# 18. 🧱 Memory Addressing Formula

x86-64 addressing basically:

```text
base + index × scale + displacement
```

NASM:

```asm
[base + index*scale + displacement]
```

| Component      | Meaning            |
| -------------- | ------------------ |
| `base`         | base register      |
| `index`        | index register     |
| `scale`        | `1`, `2`, `4`, `8` |
| `displacement` | constant offset    |

Example:

```asm
mov eax, [rax + rcx*4 + 8]
```

Meaning:

```text
address =
    RAX
  + RCX × 4
  + 8
```

---

# 19. 📚 Stack Registers

| Register | Purpose                           |
| -------- | --------------------------------- |
| `RSP`    | top of stack                      |
| `RBP`    | optional stack-frame base         |
| `RIP`    | next/current instruction location |

Typical:

```asm
push rbp
mov rbp, rsp
```

and:

```asm
mov rsp, rbp
pop rbp
ret
```

Modern optimized compilers **always use `RBP` as frame pointer** කියලා හිතන්න එපා.

---

# 20. 📞 `CALL` / `RET`

```asm
call function
```

conceptually:

```text
push return address
jump function
```

Then:

```asm
ret
```

conceptually:

```text
pop return address
jump return address
```

So:

```text
Caller
  │
  │ CALL
  ▼
Function
  │
  │ RET
  ▼
Caller continues
```

---

# 21. 🪟 Windows x64 ABI — Register Preservation

Windows x64 එකේ GPRs broadly:

| Register | Classification |
| -------- | -------------- |
| `RAX`    | volatile       |
| `RCX`    | volatile       |
| `RDX`    | volatile       |
| `R8`     | volatile       |
| `R9`     | volatile       |
| `R10`    | volatile       |
| `R11`    | volatile       |
| `RBX`    | nonvolatile    |
| `RBP`    | nonvolatile    |
| `RSI`    | nonvolatile    |
| `RDI`    | nonvolatile    |
| `R12`    | nonvolatile    |
| `R13`    | nonvolatile    |
| `R14`    | nonvolatile    |
| `R15`    | nonvolatile    |
| `RSP`    | nonvolatile    |

**Volatile** = function call එකෙන් පස්සේ callerට preserve වෙයි කියලා assume කරන්න බැහැ.

**Nonvolatile** = callee එක modify කරනවා නම් original value restore කරන්න ඕන.

---

# 22. 🐧 Linux SysV AMD64 — Register Preservation

| Register | Classification |
| -------- | -------------- |
| `RAX`    | caller-saved   |
| `RCX`    | caller-saved   |
| `RDX`    | caller-saved   |
| `RSI`    | caller-saved   |
| `RDI`    | caller-saved   |
| `R8`     | caller-saved   |
| `R9`     | caller-saved   |
| `R10`    | caller-saved   |
| `R11`    | caller-saved   |
| `RBX`    | callee-saved   |
| `RBP`    | callee-saved   |
| `R12`    | callee-saved   |
| `R13`    | callee-saved   |
| `R14`    | callee-saved   |
| `R15`    | callee-saved   |
| `RSP`    | callee-saved   |

---

# 23. 🍎 macOS Intel x86-64

macOS Intel uses a System V-derived ABI.

| Register | Typical classification |
| -------- | ---------------------- |
| `RAX`    | caller-saved           |
| `RCX`    | caller-saved           |
| `RDX`    | caller-saved           |
| `RSI`    | caller-saved           |
| `RDI`    | caller-saved           |
| `R8`     | caller-saved           |
| `R9`     | caller-saved           |
| `R10`    | caller-saved           |
| `R11`    | caller-saved           |
| `RBX`    | callee-saved           |
| `RBP`    | callee-saved           |
| `R12`    | callee-saved           |
| `R13`    | callee-saved           |
| `R14`    | callee-saved           |
| `R15`    | callee-saved           |
| `RSP`    | callee-saved           |

---

# 24. 🏠 Windows Shadow Space

Windows x64:

```text
Caller stack

┌──────────────────────┐
│ Return Address       │
├──────────────────────┤
│ Shadow Space 8       │
├──────────────────────┤
│ Shadow Space 16      │
├──────────────────────┤
│ Shadow Space 24      │
├──────────────────────┤
│ Shadow Space 32      │
├──────────────────────┤
│ Additional args      │
└──────────────────────┘
```

Total:

```text
32 bytes
```

Linux/macOS SysV:

```text
No Windows-style mandatory 32-byte shadow space.
```

---

# 25. 🧮 Floating-Point Registers

| Register family |   Width |
| --------------- | ------: |
| `XMM0-XMM15`    | 128-bit |
| `YMM0-YMM15`    | 256-bit |
| `ZMM0-ZMM31`    | 512-bit |

Conceptually:

```text
ZMM0
┌─────────────────────────────────────────────┐
│                  512 bits                    │
└─────────────────────────────────────────────┘
                       │
                       └── YMM0 = lower 256
                                │
                                └── XMM0 = lower 128
```

---

# 26. 🌊 SIMD Register Hierarchy

| 512-bit | 256-bit | 128-bit |
| ------- | ------- | ------- |
| `ZMM0`  | `YMM0`  | `XMM0`  |
| `ZMM1`  | `YMM1`  | `XMM1`  |
| `ZMM2`  | `YMM2`  | `XMM2`  |
| ...     | ...     | ...     |
| `ZMM15` | `YMM15` | `XMM15` |

With AVX-512, additional `ZMM16-ZMM31` exist on supported CPUs/OS environments.

---

# 27. 🔢 Floating-Point C/C++ Types

| C/C++         |           Typical size | Common register |
| ------------- | ---------------------: | --------------- |
| `float`       |                 32-bit | `XMM`           |
| `double`      |                 64-bit | `XMM`           |
| `long double` | ABI/platform dependent | ABI dependent   |

Example:

```c
double add(double a, double b);
```

On Windows x64:

```text
a → XMM0
b → XMM1
return → XMM0
```

Linux/macOS x86-64:

```text
a → XMM0
b → XMM1
return → XMM0
```

---

# 28. 🪟🐧🍎 FP Argument Table

| FP argument | Windows x64 | Linux SysV | macOS Intel |
| ----------- | ----------- | ---------- | ----------- |
| 1st         | `XMM0`      | `XMM0`     | `XMM0`      |
| 2nd         | `XMM1`      | `XMM1`     | `XMM1`      |
| 3rd         | `XMM2`      | `XMM2`     | `XMM2`      |
| 4th         | `XMM3`      | `XMM3`     | `XMM3`      |
| 5th         | stack*      | `XMM4`     | `XMM4`      |
| 6th         | stack*      | `XMM5`     | `XMM5`      |
| 7th         | stack*      | `XMM6`     | `XMM6`      |
| 8th         | stack*      | `XMM7`     | `XMM7`      |

* Windows x64 has important positional/mirroring rules for mixed integer/FP arguments and variadic calls, so this simplified row shouldn't be used as a complete ABI algorithm.

---

# 29. ↩️ Return Registers

| Return         | Windows x64 | Linux x86-64 | macOS x86-64 |
| -------------- | ----------- | ------------ | ------------ |
| 8-bit integer  | `AL`        | `AL`         | `AL`         |
| 16-bit integer | `AX`        | `AX`         | `AX`         |
| 32-bit integer | `EAX`       | `EAX`        | `EAX`        |
| 64-bit integer | `RAX`       | `RAX`        | `RAX`        |
| pointer        | `RAX`       | `RAX`        | `RAX`        |
| `float`        | `XMM0`      | `XMM0`       | `XMM0`       |
| `double`       | `XMM0`      | `XMM0`       | `XMM0`       |

---

# 30. 🧬 C Type → Register Example

### Windows

```c
void test(
    char a,
    short b,
    int c,
    long long d
);
```

Conceptually:

```text
a → first integer argument register family → RCX
b → second → RDX
c → third → R8
d → fourth → R9
```

Relevant widths:

```text
a → CL
b → DX
c → R8D
d → R9
```

**Important:** මෙතන `CL`, `DX`, `R8D`, `R9` කියන්නේ ABI එක වෙන වෙනම registers 4ක් තෝරනවා කියන එක නොවෙයි. ඒවා respective argument registers වල width views.

---

# 31. 🧬 Linux/macOS Same Example

```c
void test(
    char a,
    short b,
    int c,
    long long d
);
```

Conceptually:

```text
a → RDI family → DIL
b → RSI family → SI
c → RDX family → EDX
d → RCX family → RCX
```

---

# 32. 🧠 NASM Instruction Categories

| Category      | Examples                                   |
| ------------- | ------------------------------------------ |
| Data movement | `MOV`, `MOVZX`, `MOVSX`, `MOVSXD`, `LEA`   |
| Arithmetic    | `ADD`, `SUB`, `INC`, `DEC`, `IMUL`, `IDIV` |
| Logic         | `AND`, `OR`, `XOR`, `NOT`                  |
| Shift         | `SHL`, `SHR`, `SAR`                        |
| Rotate        | `ROL`, `ROR`, `RCL`, `RCR`                 |
| Compare       | `CMP`, `TEST`                              |
| Branch        | `JMP`, `JE`, `JNE`, `JG`, `JL`             |
| Function      | `CALL`, `RET`                              |
| Stack         | `PUSH`, `POP`                              |
| String        | `MOVSB`, `MOVSW`, `MOVSD`, `MOVSQ`         |
| SIMD          | `MOVDQA`, `ADDPS`, `ADDPD`, etc.           |
| SSE           | `MOVSS`, `ADDSS`, `MULSS`                  |
| AVX           | `VMOVSS`, `VADDPS`, etc.                   |
| AVX-512       | `VADDPS`, masks, ZMM etc.                  |
| Atomic        | `LOCK`, `XCHG`, `CMPXCHG`                  |
| Control       | `CPUID`, `RDTSC`, etc.                     |

---

# 33. 🔀 Conditional Jumps

| Instruction   | Condition            |
| ------------- | -------------------- |
| `JE` / `JZ`   | equal / zero         |
| `JNE` / `JNZ` | not equal            |
| `JG`          | signed greater       |
| `JGE`         | signed greater/equal |
| `JL`          | signed less          |
| `JLE`         | signed less/equal    |
| `JA`          | unsigned above       |
| `JAE`         | unsigned above/equal |
| `JB`          | unsigned below       |
| `JBE`         | unsigned below/equal |
| `JC`          | carry                |
| `JNC`         | no carry             |
| `JO`          | overflow             |
| `JNO`         | no overflow          |
| `JS`          | sign                 |
| `JNS`         | no sign              |

---

# 34. 🔍 `MOVZX` vs `MOVSX`

Very important for C/C++ interoperability.

### Zero extension

```asm
movzx eax, cl
```

```text
CL = 0xFF

EAX = 0x000000FF
```

### Sign extension

```asm
movsx eax, cl
```

If `CL = 0xFF`:

```text
EAX = 0xFFFFFFFF
```

which represents:

```text
-1
```

---

# 35. 🧠 Important Extension Instructions

| Instruction | Meaning                   |
| ----------- | ------------------------- |
| `MOVZX`     | zero extend               |
| `MOVSX`     | sign extend               |
| `MOVSXD`    | sign-extend DWORD → QWORD |
| `CWDE`      | AX → EAX                  |
| `CDQE`      | EAX → RAX                 |
| `CWD`       | AX → DX:AX                |
| `CDQ`       | EAX → EDX:EAX             |
| `CQO`       | RAX → RDX:RAX             |

---

# 36. 🏷️ NASM Sections

| Section   | Typical purpose                                  |
| --------- | ------------------------------------------------ |
| `.text`   | executable code                                  |
| `.data`   | initialized writable data                        |
| `.bss`    | uninitialized/zero-initialized storage           |
| `.rodata` | read-only data, platform/object-format dependent |

NASM:

```asm
section .text
```

```asm
section .data
```

```asm
section .bss
```

---

# 37. 📝 NASM Symbols / Labels

```asm
global add

section .text

add:
    mov eax, edi
    add eax, esi
    ret
```

| Element  | Meaning                 |
| -------- | ----------------------- |
| `global` | export symbol to linker |
| `add:`   | label/symbol            |
| `mov`    | instruction             |
| `eax`    | register                |
| `edi`    | register                |
| `ret`    | instruction             |

---

# 38. 📦 Windows vs Linux/macOS Assembly Object Files

| Platform    | Object format | Typical object |
| ----------- | ------------- | -------------- |
| Windows     | PE/COFF       | `.obj`         |
| Linux       | ELF           | `.o`           |
| macOS Intel | Mach-O        | `.o`           |

Static libraries:

| Platform | Static library |
| -------- | -------------- |
| Windows  | `.lib`         |
| Linux    | `.a`           |
| macOS    | `.a`           |

---

# 39. 🔗 Complete Build Pipeline

```text
C / C++
   │
   ▼
Compiler
   │
   ▼
Assembly
   │
   ├──────────────┐
   ▼              ▼
NASM           Compiler ASM
   │              │
   ▼              ▼
.obj / .o      .obj / .o
   │              │
   └──────┬───────┘
          ▼
     Static Library
       .lib / .a
          │
          ▼
        Linker
          │
          ▼
   Executable
```

Platform:

|                 | Windows       | Linux      | macOS Intel  |
| --------------- | ------------- | ---------- | ------------ |
| Object          | `.obj`        | `.o`       | `.o`         |
| Static lib      | `.lib`        | `.a`       | `.a`         |
| Executable      | `.exe` / PE   | ELF        | Mach-O       |
| ABI             | Microsoft x64 | SysV AMD64 | SysV-derived |
| Common compiler | MSVC/Clang    | GCC/Clang  | Clang        |

---

# 40. 🧩 C → NASM Example — Windows

C:

```c
int add(int a, int b)
{
    return a + b;
}
```

ABI:

```text
a → RCX family
b → RDX family
return → RAX family
```

Because `int = 32-bit`:

```text
a → ECX
b → EDX
return → EAX
```

NASM:

```asm
global add

section .text

add:
    mov eax, ecx
    add eax, edx
    ret
```

---

# 41. 🧩 C → NASM Example — Linux/macOS

```asm
global add

section .text

add:
    mov eax, edi
    add eax, esi
    ret
```

Because:

```text
Linux/macOS

1st → RDI family → EDI for int
2nd → RSI family → ESI for int

return → RAX family → EAX
```

---

# 42. 🔥 The Most Important Table

මේ table එක තමයි ඔයාට දැන් **memorize කරන්න ඕන core table**:

| Concept        | Windows x64      | Linux x86-64 | macOS Intel x86-64 |
| -------------- | ---------------- | ------------ | ------------------ |
| 1st integer    | `RCX`            | `RDI`        | `RDI`              |
| 2nd integer    | `RDX`            | `RSI`        | `RSI`              |
| 3rd integer    | `R8`             | `RDX`        | `RDX`              |
| 4th integer    | `R9`             | `RCX`        | `RCX`              |
| 5th integer    | Stack            | `R8`         | `R8`               |
| 6th integer    | Stack            | `R9`         | `R9`               |
| integer return | `RAX`            | `RAX`        | `RAX`              |
| 1st FP         | `XMM0`           | `XMM0`       | `XMM0`             |
| 2nd FP         | `XMM1`           | `XMM1`       | `XMM1`             |
| 3rd FP         | `XMM2`           | `XMM2`       | `XMM2`             |
| 4th FP         | `XMM3`           | `XMM3`       | `XMM3`             |
| FP return      | `XMM0`           | `XMM0`       | `XMM0`             |
| Shadow space   | 32 bytes         | None         | None               |
| Object         | `.obj` / PE-COFF | `.o` / ELF   | `.o` / Mach-O      |
| Static library | `.lib`           | `.a`         | `.a`               |

---

# 43. 🚨 One Critical Correction About Your Original Tables

ඔයා දුන්න table එකේ:

```text
1st integer / pointer
Windows → RCX
Linux   → RDI
macOS   → RDI
```

**මේක correct.**

නමුත්:

```text
1st = ECX
```

කියන එක technically shorthand එකක්.

වඩා precise:

```text
Windows ABI:
1st integer argument → RCX register family

if argument is 32-bit:
                       ↓
                      ECX
```

ඒ වගේම:

```text
Linux:
1st integer argument → RDI register family

if argument is 32-bit:
                       ↓
                      EDI
```

එහෙම නැතුව:

```text
RCX
ECX
CX
CL
```

කියන්නේ registers 4ක් කියලා හිතන්න එපා.

එක register family එකක්:

```text
             RAX
              │
       ┌──────┴──────┐
       │             │
      EAX            │
       │             │
      AX              │
       │             │
      AL              │
```

---

# 🏆 Final Mental Model

ඔයා මේ entire subject එකේ layers 6ක් වෙන වෙනම තියාගන්න:

```text
LEVEL 1
C/C++ TYPE

char
short
int
long long
pointer
float
double
       │
       ▼
LEVEL 2
SIZE / REPRESENTATION

8
16
32
64
       │
       ▼
LEVEL 3
ABI

Windows:
RCX RDX R8 R9

Linux/macOS:
RDI RSI RDX RCX R8 R9
       │
       ▼
LEVEL 4
REGISTER VIEW

64 → RAX
32 → EAX
16 → AX
 8 → AL
       │
       ▼
LEVEL 5
NASM INSTRUCTION

mov
add
sub
cmp
call
ret
       │
       ▼
LEVEL 6
OBJECT / LINKER

.obj / .o
.lib / .a
PE / ELF / Mach-O
```