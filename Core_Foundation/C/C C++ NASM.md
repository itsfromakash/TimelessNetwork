හරි Akash. මේක **C/C++ → Compiler → Assembly → ABI → Object File → Static Library → Linker → EXE** කියන සම්පූර්ණ native pipeline එකක්. මේ pipeline එක හොඳට තේරුම් ගත්තොත්, පස්සේ **Rust FFI, Go cgo, Windows API, DLLs, kernel/user-mode interfaces, reverse engineering, compiler internals** වගේ දේවල් ඉගෙනගන්න foundation එක ඉතා ශක්තිමත් වෙනවා.

මම මේක **Scratch → Intermediate → Advanced → Grand Master** ලෙස layers වලින් ගන්නවා. විශේෂයෙන් වැදගත් point එකක්: ඔයා දීලා තියෙන Assembly එක **Windows x64 ABI** සඳහායි. Linux/macOS System V ABI එකට ඒ register mapping එක වෙනස්.

---

# 🧠 මුලින්ම සම්පූර්ණ Mental Model එක

ඔයාගේ project එක මෙහෙමයි:

```text
                    YOUR SOURCE CODE
                          │
             ┌────────────┴────────────┐
             │                         │
          C / C++                   NASM
             │                         │
             ▼                         ▼
      C/C++ Compiler                NASM
             │                         │
             ▼                         ▼
       math_lib.obj               asm_func.obj
             │                         │
             └────────────┬────────────┘
                          ▼
                    Static Librarian
                          │
                          ▼
                  MathStaticLib.lib
                          │
                          │
        ┌─────────────────┴─────────────────┐
        │                                   │
      main.cpp                         core_engine.cpp
        │                                   │
        └───────────────┬───────────────────┘
                        ▼
                   C++ Compiler
                        │
                        ▼
                    .obj files
                        │
                        ▼
                      Linker
                        │
             ┌──────────┴──────────┐
             │                     │
       MainApp.obj          MathStaticLib.lib
             │                     │
             └──────────┬──────────┘
                        ▼
                   MainApp.exe
```

මේකේ ඉතා වැදගත් concept එක:

> **Compiler function එකක් call කරන්නේ නැහැ. Compiler එක machine code/object code හදනවා. Linker තමයි වෙනත් object files/libraries වල තියෙන function symbols resolve කරන්නේ.**

---

# MODULE 1 — Absolute Fundamentals & C/C++ Header Magic

---

## 1.1 `math_lib.h` කියන්නේ මොකක්ද?

```c
#ifndef MATH_LIB_H
#define MATH_LIB_H

#ifdef __cplusplus
extern "C" {
#endif

int asm_add(int a, int b);
int multiply_by_two(int x);

#ifdef __cplusplus
}
#endif

#endif
```

මුලින් `.h` file එකේ fundamental purpose එක තේරුම් ගමු.

### Header = Declaration information

Header එක කියන්නේ compiler එකට:

> "මේ functions exist කරනවා. ඒවාගේ names සහ parameter types මෙන්න."

කියලා කියන file එක.

උදාහරණයක්:

```c
int asm_add(int a, int b);
```

මේක **function definition එකක් නෙමෙයි.**

මේක **function declaration / prototype**.

---

# 1.2 Function declaration එක කඩලා බලමු

```c
int asm_add(int a, int b);
```

### `int`

Return type.

Function එක අවසානයේ integer value එකක් return කරනවා.

---

### `asm_add`

Function name.

---

### `(int a, int b)`

Parameters.

```text
asm_add
   │
   ├── parameter 1 → int a
   │
   └── parameter 2 → int b
```

---

### `;`

Declaration එක අවසන්.

Function body එක නැති නිසා `;` තියෙනවා.

---

Compare:

```c
int asm_add(int a, int b);
```

vs

```c
int asm_add(int a, int b)
{
    return a + b;
}
```

පළවෙනි එක:

```text
DECLARATION
```

දෙවෙනි එක:

```text
DEFINITION
```

---

# 1.3 Include Guards

දැන්:

```c
#ifndef MATH_LIB_H
#define MATH_LIB_H
```

මේක beginner කෙනෙකුට initially weird.

අපි ඒක step-by-step බලමු.

---

## `#ifndef`

මේක C compiler instruction එකක් නෙමෙයි.

මේක **preprocessor directive**.

`#ifndef` =

> "if not defined"

```c
#ifndef MATH_LIB_H
```

අර්ථය:

> `MATH_LIB_H` කියන macro එක තවම define කරලා නැත්නම් පහළ code process කරන්න.

---

ඊළඟ:

```c
#define MATH_LIB_H
```

දැන් `MATH_LIB_H` defined.

ඒ නිසා header එක පස්සේ නැවත include කළත්:

```c
#ifndef MATH_LIB_H
```

false වෙනවා.

ඒකෙන් header content එක නැවත process වෙන්නේ නැහැ.

---

## ඇයි මේක ඕන?

Suppose:

```cpp
#include "math_lib.h"
#include "math_lib.h"
```

හෝ indirect inclusion:

```text
main.cpp
   │
   ├── include A.h
   │       │
   │       └── include math_lib.h
   │
   └── include B.h
           │
           └── include math_lib.h
```

Header guard නැත්නම්:

```c
int asm_add(int a, int b);
int asm_add(int a, int b);
```

වගේ declarations නැවත නැවත preprocessor output එකට එන්න පුළුවන්.

බොහෝ declarations duplicate වීම problems ඇති කරනවා; definitions නම් තවත් බරපතළ **redefinition** errors ඇති කළ හැක.

Header guard එකේ mental model එක:

```text
First include
     ↓
MATH_LIB_H not defined
     ↓
process header
     ↓
#define MATH_LIB_H
     ↓
Done

Second include
     ↓
MATH_LIB_H already defined
     ↓
skip entire header
```

---

# 1.4 `__cplusplus` කියන්නේ මොකක්ද?

```c
#ifdef __cplusplus
```

මේක ඉතා වැදගත්.

`__cplusplus` කියන්නේ C++ compiler එකක් source එක compile කරන විට predefined කරන macro එකක්.

Conceptually:

```text
C compiler
    ↓
__cplusplus absent

C++ compiler
    ↓
__cplusplus defined
```

ඒ නිසා:

```c
#ifdef __cplusplus
```

කියන්නේ:

> "දැනට මේ header එක C++ compiler එකකින් process වෙනවද?"

---

# 1.5 `extern "C"` — මෙතනින් real compiler magic එක පටන් ගන්නවා

```cpp
extern "C" {
    int asm_add(int a, int b);
    int multiply_by_two(int x);
}
```

මේකේ primary purpose එක:

> **C linkage request කිරීම.**

එහෙම කියන්නේ function එකේ external symbol name එක C++ compiler එකට C++ name mangling rules වලට යටත් නොකරන්න කියන එක.

---

# 1.6 Name Mangling කියන්නේ මොකක්ද?

C++ එකේ functions වල:

```cpp
int asm_add(int a, int b);
```

වගේ function එකක් තියෙනවා.

නමුත් C++ language එක function overloading support කරනවා.

උදාහරණ:

```cpp
int add(int a, int b);
double add(double a, double b);
int add(int a);
```

මේ තුනම `add`.

Linker එකට කොහොමද differentiate කරන්නේ?

ඒ නිසා compiler එක function information එක symbol name එකට encode කරනවා.

මේක:

> **Name Mangling**

---

## GCC/Clang Itanium-style example

ඔයා සඳහන් කළ:

```text
_Z11asm_addii
```

මේකේ roughly:

```text
_Z
  │
  └── mangled C++ symbol

11
  │
  └── name length

asm_add
  │
  └── function name

ii
  │
  └── int, int
```

ඒ නිසා conceptually:

```text
C++ source:

int asm_add(int, int);

        ↓

_Z7asm_addii
```

**නමුත් මෙතන exact spelling ගැන වැදගත් correction එකක් තියෙනවා.**

`asm_add` characters 7ක් නිසා GCC/Clang Itanium ABI-style mangled form එක:

```text
_Z7asm_addii
```

වෙනවා — `_Z11asm_addii` නෙමෙයි.

---

## Windows MSVC නම්?

ඔයාගේ project එක:

```cmake
set(CMAKE_ASM_NASM_FLAGS "-f win64")
```

සහ Windows native toolchain එකක් target කරන නිසා MSVC C++ compiler එක භාවිතා කළොත් mangling format එක වෙනස්.

MSVC C++ symbol එක conceptually මෙවැනි ආකාරයක්:

```text
?asm_add@@YAHHH@Z
```

ඒ නිසා:

| Compiler/ABI              | Example             |
| ------------------------- | ------------------- |
| C                         | `asm_add`           |
| GCC/Clang C++ Itanium ABI | `_Z7asm_addii`      |
| MSVC C++                  | `?asm_add@@YAHHH@Z` |

ඒ නිසා `_Z7asm_addii` example එක **GCC/Clang-style** එකක්.

---

# 1.7 Assembly එකට මේක වැදගත් ඇයි?

ඔයාගේ NASM:

```asm
global asm_add
```

කියනකොට exported symbol එක:

```text
asm_add
```

Assembly object file එකේ symbol table එකේ තියෙනවා.

C++ compiler එක `extern "C"` නැතුව:

```cpp
int asm_add(int, int);
```

compile කළොත්, C++ compiler එක linker එකට:

```text
"I need ?asm_add@@YAHHH@Z"
```

වගේ symbol එකක් ඉල්ලන්න පුළුවන්.

Assembly object එක කියනවා:

```text
"I provide asm_add"
```

Linker:

```text
Requested:
?asm_add@@YAHHH@Z

Available:
asm_add

        ↓

NO MATCH
```

ඒ නිසා:

```text
unresolved external symbol
```

---

# 1.8 `extern "C"` දාපු විට

C++ compiler එකට:

```cpp
extern "C" int asm_add(int, int);
```

කියනවා:

> "මේ function එක C linkage එකෙන් external symbol එකක්."

එවිට expected external symbol:

```text
asm_add
```

Assembly:

```text
asm_add
```

Match:

```text
C++ object
      │
      │ requests
      ▼
   asm_add
      ▲
      │ provides
      │
NASM object
```

💡 **මෙතන `extern "C"` Assembly calling convention එක වෙනස් කරන්නේ නැහැ. එය primarily linkage/name mangling එක control කරනවා.**

---

# 1.9 `math_lib.c`

දැන්:

```c
#include "math_lib.h"

int multiply_by_two(int x) {
    return asm_add(x, x);
}
```

---

## Preprocessor stage

```c
#include "math_lib.h"
```

literally:

> header content එක source file එකට ඇතුළත් කරලා preprocessing කරන්න.

Conceptually:

```text
math_lib.c
   +
math_lib.h
   ↓
preprocessed C source
```

---

## Function

```c
int multiply_by_two(int x)
```

අර්ථය:

```text
function name:
multiply_by_two

input:
int x

output:
int
```

Body:

```c
{
    return asm_add(x, x);
}
```

---

Execution:

```text
multiply_by_two(5)
       │
       ▼
asm_add(5, 5)
       │
       ▼
5 + 5
       │
       ▼
10
       │
       ▼
return 10
```

---

# MODULE 2 — x86_64 NASM & Windows x64 ABI

දැන් CPU level එකට යමු. 🔥

ඔයාගේ Assembly:

```asm
bits 64
default rel

section .text
global asm_add

asm_add:
    mov rax, rcx
    add rax, rdx
    ret
```

---

# 2.1 `bits 64`

```asm
bits 64
```

NASM assembler එකට කියන්නේ:

> මේ source එක 64-bit mode සඳහා assemble කරන්න.

ඒක නිසා registers:

```text
RAX
RBX
RCX
RDX
RSI
RDI
RSP
RBP
R8-R15
```

වැනි 64-bit registers භාවිතා කළ හැක.

---

# 2.2 `default rel`

```asm
default rel
```

`rel` = relative addressing.

විශේෂයෙන් memory operands වලදී NASM එකට RIP-relative addressing prefer කරන්න කියන directive එකක්.

උදාහරණයක් conceptually:

```asm
mov rax, [some_data]
```

`default rel` තිබෙන 경우 assembler එක RIP-relative encoding භාවිතා කළ හැක.

x86-64 position-independent code / relocatable code design වල මේක වැදගත්.

**නමුත් ඔයාගේ current function එකේ memory operand එකක් නැති නිසා `default rel` practical effect එකක් නැහැ.**

---

# 2.3 RIP-relative addressing

x86-64:

```text
RIP = Instruction Pointer
```

Example:

```asm
mov rax, [rel value]
```

conceptually:

```text
effective address =
current RIP + displacement
```

ඒක absolute address එකක් hard-code කරනවාට වඩා relocation-friendly.

Windows PE/COFF සහ Linux ELF දෙකේම x86-64 code generation වල මේ concept එක වැදගත්.

---

# 2.4 `section .text`

```asm
section .text
```

`.text` = executable machine code section.

Object/EXE එකේ conceptually:

```text
.text
    ↓
machine instructions

.data
    ↓
initialized writable data

.bss
    ↓
zero-initialized/uninitialized storage
```

ඔයාගේ:

```asm
mov
add
ret
```

instructions `.text` section එකට යනවා.

---

# 2.5 `global asm_add`

```asm
global asm_add
```

NASM එකට:

> `asm_add` symbol එක external linker එකට visible/exported symbol එකක් කරන්න.

`global` නැත්නම් symbol එක local/private වගේ object-level visibility එකකට සීමා වෙන්න පුළුවන්.

---

# 2.6 Windows x64 Calling Convention

මේක **අතිශය වැදගත්**.

Windows x64 ABI එකේ integer/pointer arguments සඳහා first four register arguments:

```text
1st → RCX
2nd → RDX
3rd → R8
4th → R9
```

උදාහරණය:

```c
int f(int a, int b, int c, int d);
```

Windows x64:

```text
a → RCX
b → RDX
c → R8
d → R9
```

---

## System V AMD64 — Linux/macOS

Linux/macOS x86-64 වල:

```text
1st → RDI
2nd → RSI
3rd → RDX
4th → RCX
5th → R8
6th → R9
```

ඒ නිසා:

```c
asm_add(5, 5)
```

Windows:

```text
RCX = 5
RDX = 5
```

System V:

```text
RDI = 5
RSI = 5
```

---

# 2.7 මේ දෙක mix කළොත්?

Suppose Linux ABI code එක:

```asm
mov rax, rdi
add rax, rsi
```

Windows caller එකකින් call කළොත් wrong.

Windows caller:

```text
RCX = first argument
RDX = second argument
```

ඒ නිසා Assembly:

```asm
mov rax, rcx
add rax, rdx
```

හරි.

---

# 2.8 Return value

Integer return value:

```text
RAX / EAX
```

32-bit `int` සඳහා effectively:

```text
EAX
```

Return:

```asm
mov eax, ...
ret
```

---

# 2.9 දැන් Assembly line-by-line

```asm
asm_add:
```

මේක label එකක්.

Machine code execution එතනින් පටන් ගන්නවා.

---

### First instruction

```asm
mov rax, rcx
```

Meaning:

```text
RAX ← RCX
```

If:

```text
RCX = 5
```

then:

```text
RAX = 5
```

---

### Second

```asm
add rax, rdx
```

Meaning:

```text
RAX = RAX + RDX
```

If:

```text
RAX = 5
RDX = 5
```

then:

```text
RAX = 10
```

---

### `ret`

```asm
ret
```

CPU stack එකෙන් return address එක pop කරලා ඒ address එකට execution return කරනවා.

Conceptually:

```text
CALL asm_add
      │
      ├── push return address
      │
      └── jump asm_add

asm_add:
      ...
      ret
      │
      └── pop return address
          jump back
```

---

# 2.10 `int` එකකට `RAX` use කරන එක vs `EAX`

මේක subtle.

C:

```c
int
```

සාමාන්‍ය Windows x64 environment එකේ:

```text
32 bits
```

ඒ කියන්නේ:

```text
4 bytes
```

නමුත් ඔයා Assembly එකේ:

```asm
RAX
RCX
RDX
```

use කරනවා.

මේවා:

```text
64 bits
8 bytes
```

---

## Register relationship

```text
RAX = 64-bit
┌────────────────────────────────────────────────┐
│                     RAX                        │
└────────────────────────────────────────────────┘
                        │
                        └── lower 32 bits = EAX
```

Similarly:

```text
RCX → ECX
RDX → EDX
```

---

## 32-bit operation

```asm
mov eax, ecx
add eax, edx
```

මේක mathematically:

```text
EAX = ECX + EDX
```

32-bit integer arithmetic සඳහා more semantically precise.

---

## Important x86 rule

32-bit destination register එකකට write කළොත් upper 32 bits of corresponding 64-bit register zero වෙනවා.

Example:

```asm
mov eax, 10
```

results:

```text
RAX = 00000000_0000000A
```

---

But:

```asm
mov ax, 10
```

upper bits zero වෙන්නේ නැහැ.

```asm
mov al, 10
```

upper bits zero වෙන්නේ නැහැ.

32-bit write එක special.

---

# 2.11 ඔයාගේ code එක technically works ද?

ඔව්, low 32-bit result එක සඳහා current function එක සාමාන්‍ය `int` use case එකේ වැඩ කරන්න පුළුවන්:

```asm
mov rax, rcx
add rax, rdx
ret
```

නමුත් semantically cleaner implementation:

```asm
mov eax, ecx
add eax, edx
ret
```

ඇයි?

C function:

```c
int asm_add(int a, int b);
```

කියන්නේ 32-bit integers.

ඒ නිසා Assembly එක:

```text
32-bit input values
32-bit arithmetic
32-bit return
```

වලට align වෙනවා.

---

## REX prefixes ගැන

64-bit operand instructions වලදී x86-64 encoding එකට REX prefix අවශ්‍ය විය හැක.

```asm
add rax, rdx
```

vs

```asm
add eax, edx
```

64-bit version එකට 64-bit operand-size encoding එකක් අවශ්‍යයි.

32-bit version එක සාමාන්‍යයෙන් smaller instruction encoding එකක්.

ඒ නිසා `int` සඳහා:

```asm
mov eax, ecx
add eax, edx
```

**cleaner සහ slightly smaller**.

---

# MODULE 3 — Stack, Shadow Space & Register Preservation

දැන් CPU registers වලින් stack architecture එකට යමු. 🧠

Call:

```c
multiply_by_two(5)
```

ඇතුළේ:

```c
return asm_add(x, x);
```

එනිසා:

```text
multiply_by_two(5)
        │
        ▼
asm_add(5,5)
```

---

# 3.1 Windows x64 Shadow Space

Windows x64 ABI එකේ caller එකක් function call එකක් කරන විට callee සඳහා:

```text
32 bytes
```

**shadow space / home space**

reserve කරනවා.

එය:

```text
4 × 8 bytes
```

because:

```text
4 register argument slots
```

වගේ conceptually හිතන්න පුළුවන්.

---

# 3.2 කවුද allocate කරන්නේ?

**Caller.**

උදාහරණ:

```text
multiply_by_two
      │
      │ CALL asm_add
      ▼
   asm_add
```

`multiply_by_two` තමයි `asm_add` call කරන්න කලින් shadow space reserve කරන්නේ.

---

# 3.3 කවුද clean කරන්නේ?

Caller තමයි.

Callee:

```asm
ret
```

කළාම shadow space automatically cleanup වෙන්නේ නැහැ.

Caller තමන් reserve කළ stack space එක restore කරනවා.

---

# 3.4 Stack diagram

Simplified view:

```text
Higher addresses
────────────────────────────

Caller stack data

Return address
────────────────────────────  ← [RSP at callee entry + 0]

Shadow/Home Space
────────────────────────────
8 bytes
────────────────────────────
8 bytes
────────────────────────────
8 bytes
────────────────────────────
8 bytes
────────────────────────────

Lower addresses
```

Call එකේදී CPU:

```asm
call asm_add
```

කරන විට return address එක stack එකට push වෙනවා.

---

# 3.5 Callee entry RSP

`CALL` return address එක 8 bytes push කරන නිසා function entry එකේ:

```text
RSP = return-address location
```

ඒ නිසා Windows x64 ABI stack alignment discussion එකේ:

```text
callee entry RSP ≡ 8 (mod 16)
```

වීම common invariant එක.

Caller call කරන්න කලින් stack properly arranged කරලා තියෙනවා.

---

# 3.6 16-byte stack alignment

x86-64 Windows ABI එකේ stack alignment requirement එක important.

සාමාන්‍ය rule:

> A `CALL` එක execute කිරීමට පෙර caller's `RSP` 16-byte aligned වෙලා තිබිය යුතුයි.

Conceptually:

```text
RSP % 16 = 0
```

before `CALL`.

`CALL`:

```text
push return address
```

කරන නිසා callee entry:

```text
RSP % 16 = 8
```

වෙනවා.

---

# 3.7 Shadow space vs alignment

මේ දෙක confuse කරන්න එපා.

```text
Shadow space
    ↓
32 bytes
    ↓
ABI requirement for callee's home/register arguments

Stack alignment
    ↓
16-byte boundary
    ↓
ABI/ISA-compatible stack alignment requirement
```

දෙකම වෙනස් concepts.

---

# 3.8 Volatile Registers

Windows x64 ABI එකේ caller-saved / volatile registers:

```text
RAX
RCX
RDX
R8
R9
R10
R11
```

Floating point/vector registers වලත් corresponding volatility rules තියෙනවා.

Volatile කියන්නේ:

> Function call එකකින් පස්සේ මේ registers වල old values preserve වෙයි කියලා caller එකට guarantee එකක් නැහැ.

---

Example:

```asm
mov rax, 100
call something
```

`something` return වුණාට පස්සේ:

```text
RAX = 100
```

වෙයි කියලා assume කරන්න බැහැ.

---

# 3.9 Non-Volatile Registers

Windows x64 integer non-volatile registers:

```text
RBX
RBP
RDI
RSI
RSP
R12
R13
R14
R15
```

`RSP` special stack pointer එකක්.

Function එකක්:

```asm
mov rbx, 123
```

කරලා return වෙනවා නම් original `RBX` preserve කරන්න ඕන.

---

## Wrong Assembly

```asm
my_func:
    mov rbx, 999
    ret
```

Caller එකේ:

```text
RBX = important value
```

තිබුණොත් `my_func` ඒක destroy කරනවා.

එය ABI violation.

Result:

```text
random bugs
corrupted state
crashes
incorrect program behavior
```

වෙන්න පුළුවන්.

---

# 3.10 Correct preservation

```asm
my_func:
    push rbx

    mov rbx, 999

    pop rbx
    ret
```

දැන්:

```text
old RBX
   ↓
stack

RBX = 999

...

restore old RBX
```

---

# 3.11 ඔයාගේ `asm_add` එකට save කිරීම අවශ්‍යද?

නැහැ.

ඔයා භාවිතා කරන්නේ:

```text
RAX
RCX
RDX
```

මේවා volatile.

ඒ නිසා preserve කරන්න අවශ්‍ය නැහැ.

---

# MODULE 4 — Object Files, Symbols & CMake Pipeline

දැන් source code CPU instructions වලට convert වෙන journey එක බලමු.

---

# 4.1 `math_lib.c`

Source:

```c
#include "math_lib.h"

int multiply_by_two(int x) {
    return asm_add(x, x);
}
```

C compiler එක මේක compile කරනවා.

Conceptually:

```text
math_lib.c
   ↓
Preprocessor
   ↓
Compiler frontend
   ↓
C AST / semantic analysis
   ↓
Intermediate Representation
   ↓
Optimization
   ↓
Machine code generation
   ↓
COFF object
```

Windows native toolchain එකේ:

```text
math_lib.obj
```

---

# 4.2 Object file කියන්නේ EXE එකක් නෙමෙයි

මේ distinction එක මතක තියාගන්න.

`.obj`:

```text
incomplete machine-code package
```

එහි තිබිය හැක:

```text
machine instructions
symbol table
relocations
sections
debug information
```

නමුත් external references තවම unresolved වෙලා තිබිය හැක.

---

# 4.3 `math_lib.obj` තුළ

Compiler එක `multiply_by_two` machine code එක generate කරනවා.

ඒ code එක `asm_add` call කරන නිසා object file එක conceptually කියනවා:

```text
I DEFINE:
    multiply_by_two

I NEED:
    asm_add
```

---

# 4.4 NASM

Command conceptually:

```text
nasm -f win64 src/asm_func.asm -o asm_func.obj
```

`-f win64`:

> Windows 64-bit COFF object format generate කරන්න.

Result:

```text
asm_func.obj
```

---

# 4.5 NASM object එක

එහි:

```text
DEFINED SYMBOL:
    asm_add
```

වගේ symbol එකක් තියෙනවා.

Assembly:

```asm
global asm_add
```

නිසා linker එකට symbol visible.

---

# 4.6 Static Library

CMake:

```cmake
add_library(MathStaticLib STATIC
    src/asm_func.asm
    src/math_lib.c
)
```

`STATIC` කියන්නේ static library.

Windows MSVC ecosystem එකේ result:

```text
MathStaticLib.lib
```

Conceptually:

```text
MathStaticLib.lib
       │
       ├── asm_func.obj
       │
       └── math_lib.obj
```

Static library basically object files collection/archive එකක්.

---

# 4.7 `.lib` කියන්නේ එකම දෙයක් නෙමෙයි

Windows `.lib` extension එක දෙවර්ගයකට භාවිතා වෙන්න පුළුවන්:

```text
Static library
Import library
```

ඔයාගේ:

```text
MathStaticLib.lib
```

static library.

DLL එකකට associated import `.lib` එකක්ත් තියෙන්න පුළුවන්.

---

# 4.8 Executable

```cmake
add_executable(MainApp
    src/main.cpp
    src/core_engine.cpp
)
```

CMake compiler එකෙන්:

```text
main.cpp
    ↓
main.obj

core_engine.cpp
    ↓
core_engine.obj
```

---

# 4.9 `target_link_libraries`

```cmake
target_link_libraries(MainApp PRIVATE MathStaticLib)
```

මේක කියන්නේ:

> `MainApp` target එක link කරන විට `MathStaticLib` library එකත් linker inputs වලට include කරන්න.

---

# 4.10 Linker Symbol Resolution

මේක entire pipeline එකේ heart එක. ❤️

Suppose:

```text
main.obj
```

needs:

```text
multiply_by_two
```

and:

```text
math_lib.obj
```

provides:

```text
multiply_by_two
```

එම object එක needs:

```text
asm_add
```

and:

```text
asm_func.obj
```

provides:

```text
asm_add
```

Linker:

```text
main.obj
    │
    │ undefined: multiply_by_two
    ▼
math_lib.obj
    │
    │ provides multiply_by_two
    │
    │ undefined: asm_add
    ▼
asm_func.obj
    │
    │ provides asm_add
    ▼
RESOLVED
```

---

# 4.11 Final executable

Linker:

```text
main.obj
core_engine.obj
math_lib.obj
asm_func.obj
runtime libraries
system libraries
        │
        ▼
    MainApp.exe
```

---

# 4.12 Relocation

මෙතන advanced concept එකක් තියෙනවා.

`math_lib.obj` තුළ compiler එක:

```asm
call asm_add
```

වගේ machine code generate කළත් `asm_add` final executable memory address එක compile කරන වෙලාවේ known නැහැ.

ඒ නිසා object file එකේ:

```text
relocation entry
```

තියෙනවා.

Linker final layout එක තීරණය කළ පස්සේ:

```text
asm_add address = 0x00007FF...
```

වගේ actual address relationship එක resolve කරලා machine code/relocation data fix කරනවා.

---

# 4.13 CMake ඇත්තටම compiler එකද?

නැහැ.

මේ distinction එක **grand-master level එකටත්** වැදගත්.

CMake:

```text
Build system generator / build configuration system
```

එය සාමාන්‍යයෙන් compiler එක වෙනුවට machine code generate කරන්නේ නැහැ.

CMake:

```text
ඔයා කියනවා:
    මේ targets තියෙනවා
    මේ languages තියෙනවා
    මේ libraries link කරන්න
```

CMake:

```text
MSVC / Ninja / Visual Studio / Makefiles
```

වගේ actual build system/toolchain එකකට instructions generate කරනවා.

---

# 4.14 `project(...)`

```cmake
project(NativePipelineProject C CXX ASM_NASM)
```

මෙහි:

```text
NativePipelineProject
```

project name.

Languages:

```text
C
CXX
ASM_NASM
```

`CXX` = C++.

`ASM_NASM` = NASM assembler language.

---

# 4.15 C standard

```cmake
set(CMAKE_C_STANDARD 11)
```

C language standard:

```text
C11
```

---

# 4.16 C++ standard

```cmake
set(CMAKE_CXX_STANDARD 17)
```

C++:

```text
C++17
```

---

# 4.17 `enable_language`

```cmake
enable_language(ASM_NASM)
```

NASM language support enable කරනවා.

Note එකක්: `project(... ASM_NASM)` already enables that language in normal CMake usage, so explicit:

```cmake
enable_language(ASM_NASM)
```

මෙතන generally redundant.

---

# 4.18 NASM flags

```cmake
set(CMAKE_ASM_NASM_FLAGS "-f win64")
```

NASM object format:

```text
win64
```

එනම් Windows x64 COFF.

Modern CMake projects වල target-specific/source-specific configuration use කිරීම global `CMAKE_ASM_NASM_FLAGS` වලට වඩා maintainable විය හැකි අවස්ථා තියෙනවා.

---

# 4.19 Include directories

```cmake
target_include_directories(MathStaticLib PUBLIC include)
```

`include` directory එක compiler header search path එකට add කරනවා.

ඒ නිසා:

```c
#include "math_lib.h"
```

කියද්දී compiler එක:

```text
include/math_lib.h
```

find කරන්න පුළුවන්.

---

`PUBLIC` කියන්නේ dependency propagation.

```text
MathStaticLib
    │
    └── PUBLIC include directory
              │
              ▼
       consumers may inherit
```

---

# 4.20 MainApp include path

```cmake
target_include_directories(MainApp PRIVATE include)
```

`MainApp` itself `include` directory එක use කරනවා.

`PRIVATE`:

> MainApp consumer එකකට මේ include path propagate කරන්න අවශ්‍ය නැහැ.

---

# MODULE 5 — Grand Master Optimizations & Real-World Edge Cases

දැන් optimization.

---

# 5.1 Current Assembly

```asm
asm_add:
    mov rax, rcx
    add rax, rdx
    ret
```

For C `int`:

```asm
asm_add:
    mov eax, ecx
    add eax, edx
    ret
```

මේක semantically better.

---

# 5.2 Can we remove the `mov`?

ඔව්.

```asm
asm_add:
    lea eax, [rcx + rdx]
    ret
```

නමුත් මෙතන register width mismatch issue තියෙන නිසා:

```asm
lea eax, [ecx + edx]
```

කියලා ලියන්න බැහැ? Actually x86 addressing operands allow 32-bit base/index registers? In 64-bit mode address generation uses 64-bit registers; `lea eax,[rcx+rdx]` computes 64-bit effective address then truncates to 32 bits, which works for low bits, but using RCX/RDX assumes meaningful upper bits. For ABI int args, upper bits aren't semantically part of int. So `lea eax,[rcx+rdx]` could still produce same low32 modulo 2^32 because addition modulo 2^64 then truncation, regardless upper bits, but address-generation overflow semantics modulo 2^64. Yet `lea` with two registers uses REX, 4-byte dest. Could be one instruction vs mov+add. But for signed int addition, C overflow is UB. For all defined additions, low32 is right. However, `lea eax,[ecx+edx]` invalid. Better maybe `lea eax,[rcx+rdx]`, but not necessarily "maximum performance". On modern x86, `lea` may have different latency/ports and `add` is 1 µop. `lea` one instruction one µop often. But `mov eax,ecx` is eliminated by move elimination on many CPUs; `add eax,edx` is dependency chain. `lea eax,[rcx+rdx]` directly combines. Yet upper bits irrelevant. But ABI doesn't guarantee upper bits? Windows x64 says integer args in regs, args smaller than 64 are right-justified, upper bits unspecified. `lea eax,[rcx+rdx]` still low32 correct due modulo arithmetic. Good. Could even `lea eax,[ecx+edx]` no. Could `add ecx, edx; mov eax, ecx` modifies volatile RCX so okay, 2 instructions and move. `lea` one instruction.

But "maximum performance" for a function is context-dependent; likely `lea eax,[rcx+rdx]` or `add ecx, edx; mov eax,ecx`? Let's be precise: `lea eax,[rcx+rdx]` uses 64-bit address generation, destination 32 bit. It does not access memory. But address size 64 and arithmetic modulo 2^64. Correct low 32. However `lea eax,[rcx+rdx]` encoding is 7? `lea eax,[rcx+rdx]` likely 4? REX + opcode + modrm + sib = 4? REX.W absent but REX.R/B needed = 1 + 1+1+1 = 4. `mov eax,ecx` 2, add eax,edx 2 => 4 total. So same code size. `lea` one instruction. Could also `lea eax,[rcx+rdx]` likely one µop but port restrictions. `add eax,edx` after move elimination maybe 1 µop. `lea` might be 1 µop. Fine.

Could `lea eax,[rcx+rdx]` be considered incorrect for signed overflow? It computes low32, same modulo 2^32. But if C semantics require UB, any result is allowed on overflow, so no issue. For non-overflow, exact. Good.

---

# 5.3 Calling `printf` from NASM

Now important.

Suppose:

```c
void asm_print(int x);
```

and Assembly:

```asm
extern printf

section .rdata
fmt db "Value = %d", 13, 10, 0

section .text
global asm_print

asm_print:
    ...
```

Windows x64 ABI.

`printf` signature:

```c
int printf(const char *format, ...);
```

First argument:

```text
RCX
```

Second integer argument:

```text
RDX
```

---

## Need shadow space

At call site:

```text
32-byte shadow space
```

must be available.

Need alignment too.

Assume function entry:

```text
RSP % 16 = 8
```

We need before CALL:

```text
RSP % 16 = 0
```

Allocate:

```text
40 bytes
```

because:

```text
32 shadow
+
8 alignment padding
=
40
```

---

## Assembly

```asm
bits 64
default rel

extern printf

section .rdata
fmt db "Value = %d", 13, 10, 0

section .text
global asm_print

asm_print:
    sub rsp, 40

    lea rcx, [rel fmt]
    mov edx, 123

    call printf

    add rsp, 40
    ret
```

---

# 5.4 මේක line-by-line

### `extern printf`

```asm
extern printf
```

Assembly එක කියනවා:

> `printf` මේ object file එකේ define කරලා නැහැ. Linker එකෙන් වෙනත් object/library එකකින් resolve කරන්න.

---

### `.rdata`

```asm
section .rdata
```

Read-only initialized data සඳහා.

---

### Format string

```asm
fmt db "Value = %d", 13, 10, 0
```

`db` = Define Bytes.

Memory:

```text
V
a
l
u
e
...
%
d
CR
LF
NUL
```

`printf` C string එකක් expect කරන නිසා final:

```text
0
```

NUL terminator.

---

### Function entry

```asm
asm_print:
```

---

### Stack allocation

```asm
sub rsp, 40
```

40 bytes reserve.

```text
32 bytes → shadow space
8 bytes  → alignment padding
```

---

### First argument

```asm
lea rcx, [rel fmt]
```

`RCX`:

```text
printf first argument
```

එය format string address එක.

`LEA`:

> memory load එකක් නොකර effective address calculate කරන instruction එක.

---

### Second argument

```asm
mov edx, 123
```

`printf` සඳහා second argument.

Windows x64:

```text
RCX = arg1
RDX = arg2
```

---

### Call

```asm
call printf
```

CPU:

```text
push return address
jump printf
```

---

### Cleanup

```asm
add rsp, 40
```

Own stack allocation restore.

---

### Return

```asm
ret
```

---

# ⚠️ Variadic functions: `printf`

`printf` variadic function එකක්:

```c
printf(const char *format, ...);
```

Windows x64 ABI එකේ floating-point arguments සහ varargs සඳහා additional ABI details තියෙනවා. උදාහරණයක් ලෙස floating-point varargs register passing වල integer-register duplication rules relevant වෙන්න පුළුවන්.

ඒ නිසා:

```text
printf("%f", double_value)
```

වගේ Assembly calls, `printf("%d", int)` ට වඩා ABI-sensitive.

---

# 5.5 Non-volatile register use කරන්නේ නම්

Suppose:

```asm
asm_func:
    mov r12, 123
    ret
```

❌ ABI violation.

Correct:

```asm
asm_func:
    push r12

    mov r12, 123

    pop r12
    ret
```

නමුත් stack alignment/prologue design එකත් consider කරන්න ඕන.

---

# 5.6 Stack frame එකක් අවශ්‍යමද?

නැහැ.

ඔයාගේ:

```asm
asm_add:
    mov eax, ecx
    add eax, edx
    ret
```

function එක **leaf function**.

එය:

```text
doesn't call another function
doesn't need local stack storage
doesn't modify nonvolatile registers
```

ඒ නිසා prologue නැති:

```text
leaf / frameless function
```

එකක් වෙන්න පුළුවන්.

---

# 5.7 Grand-Master version

For your exact signature:

```c
int asm_add(int a, int b);
```

Windows x64:

```asm
bits 64
default rel

section .text
global asm_add

asm_add:
    lea eax, [rcx + rdx]
    ret
```

මේක interesting.

`LEA` memory access එකක් **කරන්නේ නැහැ**.

```text
LEA = Load Effective Address
```

නම නිසා memory load කරනවා කියලා හිතන්න එපා.

It computes:

```text
EAX = low32(RCX + RDX)
```

and returns.

---

## නමුත් "maximum performance" කියන්නේ මේකමද?

**Automatically no.**

Performance depends on:

```text
CPU microarchitecture
dependency chains
inlining
caller context
register pressure
instruction throughput
latency
code size
compiler optimization
```

ඉතා වැදගත් point එක:

### C compiler එකට මේ function එක දුන්නොත්

```c
int asm_add(int a, int b)
{
    return a + b;
}
```

optimized compiler එක මේක inline කරලා call එකම eliminate කරන්න පුළුවන්.

Assembly external function එකක් නම් compiler එකට generally inline body එක දන්නේ නැහැ.

ඒ නිසා tiny Assembly function එකක් manually optimize කරනවාට වඩා:

```text
inlining
LTO
PGO
compiler optimization
```

කාලෙකදී වැඩි performance impact එකක් දෙන්න පුළුවන්.

---

# 5.8 Static Library එකේ "asm_add" call flow

අපි final CPU execution එක trace කරමු.

Suppose:

```cpp
int main()
{
    return multiply_by_two(5);
}
```

Conceptual flow:

```text
MAIN
 │
 │ call multiply_by_two
 ▼
multiply_by_two
 │
 │ x = 5
 │
 │ RCX = 5
 │
 │ RDX = 5
 │
 │ call asm_add
 ▼
asm_add
 │
 │ EAX = ECX + EDX
 │
 │ EAX = 10
 │
 │ ret
 ▼
multiply_by_two
 │
 │ return 10
 ▼
main
```

---

# 🔥 Full Memory/CPU Picture

At a conceptual level:

```text
                 MAIN.EXE
┌─────────────────────────────────────────┐
│                                         │
│ main()                                  │
│   │                                     │
│   │ call multiply_by_two                │
│   ▼                                     │
│ multiply_by_two()                       │
│   │                                     │
│   │ call asm_add                        │
│   ▼                                     │
│ asm_add()                               │
│                                         │
│   RCX = first argument                  │
│   RDX = second argument                 │
│                                         │
│   EAX = ECX + EDX                       │
│                                         │
│   RET                                   │
│                                         │
└─────────────────────────────────────────┘
```

---

# 🧩 Linker Error Debugging Master Checklist

දැන් practical debugging.

---

## Error 1 — `LNK2019`

Example:

```text
unresolved external symbol asm_add
```

First question:

> Is `asm_add` actually present in the object/library?

Check:

```text
dumpbin /symbols asm_func.obj
```

or library:

```text
dumpbin /symbols MathStaticLib.lib
```

Look for:

```text
asm_add
```

---

# Error 2 — C++ mangling mismatch

C++ asks:

```text
?asm_add@@YAHHH@Z
```

Assembly provides:

```text
asm_add
```

Solution:

```cpp
#ifdef __cplusplus
extern "C" {
#endif

int asm_add(int, int);

#ifdef __cplusplus
}
#endif
```

---

# Error 3 — Wrong ABI

Assembly expects:

```text
RDI
RSI
```

but Windows caller uses:

```text
RCX
RDX
```

Then arguments are wrong.

For Windows x64:

```text
arg1 → RCX
arg2 → RDX
arg3 → R8
arg4 → R9
```

---

# Error 4 — Wrong object format

You accidentally assemble:

```text
-f elf64
```

on a Windows MSVC pipeline.

But linker expects:

```text
COFF / win64
```

Use:

```text
-f win64
```

---

# Error 5 — Function wasn't exported

Assembly:

```asm
asm_add:
```

but no:

```asm
global asm_add
```

Then linker visibility can fail.

---

# Error 6 — Wrong symbol spelling

Check exact symbols.

Possible differences:

```text
asm_add
_asm_add
?asm_add@@YAHHH@Z
_Z7asm_addii
```

These are **different symbols**.

Linker doesn't think:

> "They're basically the same."

It compares symbol identities/encodings.

---

# Error 7 — Library isn't actually linked

CMake:

```cmake
target_link_libraries(MainApp PRIVATE MathStaticLib)
```

must connect the target dependency.

If missing:

```text
math_lib.obj
```

may contain:

```text
undefined asm_add
```

and final link fails.

---

# Error 8 — `ld.exe cannot find`

Usually investigate:

```text
library search path
library name
toolchain
architecture
target configuration
```

Example:

```text
x86 vs x64
```

Mixing:

```text
32-bit object
```

with:

```text
64-bit object
```

is a fundamental toolchain mismatch.

---

# 🧠 The 5 Most Important ABI Rules to Memorize

For **Windows x64**:

### Argument registers

```text
1 → RCX
2 → RDX
3 → R8
4 → R9
```

### Return

```text
integer → RAX/EAX
```

### Volatile

```text
RAX
RCX
RDX
R8
R9
R10
R11
```

### Non-volatile

```text
RBX
RBP
RDI
RSI
R12
R13
R14
R15
```

plus `RSP` as the stack pointer.

### Call requirements

```text
32-byte shadow space
+
16-byte stack alignment requirement
```

---

# 🆚 Windows x64 vs System V AMD64 — Cheat Sheet

| Concept          | Windows x64             | System V AMD64                          |
| ---------------- | ----------------------- | --------------------------------------- |
| 1st integer arg  | RCX                     | RDI                                     |
| 2nd              | RDX                     | RSI                                     |
| 3rd              | R8                      | RDX                                     |
| 4th              | R9                      | RCX                                     |
| 5th              | stack                   | R8                                      |
| 6th              | stack                   | R9                                      |
| Integer return   | RAX/EAX                 | RAX/EAX                                 |
| Shadow space     | **32 bytes**            | **No Windows-style 32-byte home space** |
| Stack alignment  | 16-byte ABI alignment   | 16-byte ABI alignment                   |
| Nonvolatile GPRs | RBX,RBP,RDI,RSI,R12-R15 | RBX,RBP,R12-R15,RSP                     |

**මේ table එක memorize කරන්න.** Cross-platform Assembly/FFI කරනකොට මේක core knowledge.

---

# 🏆 Final Grand-Master Mental Model

මේ entire project එක තේරුම් ගන්න එක sentence එකකට compress කළොත්:

> **C/C++ compiler එක source code එකෙන් object code සහ external symbol references හදනවා; NASM එක Assembly එක COFF object එකකට assemble කරනවා; CMake ඒ targets/dependencies build system එකට wire කරනවා; static librarian object files archive කරලා `.lib` එකක් හදනවා; linker symbol tables සහ relocations භාවිතා කරලා references resolve කර executable එක හදනවා; runtime එකේ ABI එක CPU registers, stack, return values සහ register preservation rules මගින් functions එකිනෙකට communicate කරන contract එකක් ලබා දෙනවා.**

ඔයාගේ exact pipeline එක:

```text
                 ┌──────────────────┐
                 │   math_lib.h     │
                 │ declarations     │
                 │ extern "C"       │
                 └────────┬─────────┘
                          │
             ┌────────────┴─────────────┐
             │                          │
             ▼                          ▼
        math_lib.c                 asm_func.asm
             │                          │
             ▼                          ▼
       C Compiler                     NASM
             │                          │
             ▼                          ▼
       math_lib.obj                asm_func.obj
             │                          │
             └────────────┬─────────────┘
                          ▼
                     ARCHIVER
                          │
                          ▼
                MathStaticLib.lib
                          │
                          ▼
                     LINKER ◄──── main.obj
                       ▲
                       │
                  core_engine.obj
                       │
                       ▼
                    MainApp.exe
                       │
                       ▼
                    CPU EXECUTION
                       │
                       ▼
             Windows x64 ABI
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
         RCX          RDX          RAX
        arg #1       arg #2       return
          │            │
          └──────┬─────┘
                 ▼
              asm_add
                 │
                 ▼
             EAX = ECX + EDX
                 │
                 ▼
                RET
```

### එක් crucial correction එකක් අවසානයේ:

ඔයාගේ original:

```asm
mov rax, rcx
add rax, rdx
ret
```

**low 32-bit `int` result සඳහා වැඩ කළ හැකි** වුණත්, function signature එක `int` නිසා මම production version එකට:

```asm
bits 64
default rel

section .text
global asm_add

asm_add:
    lea eax, [rcx + rdx]
    ret
```

වගේ implementation එක consider කරනවා. එහෙත් "optimization" කියන්නේ instruction count පමණක් නොවන බව මතක තියාගන්න—**ABI correctness → compiler interaction → inlining/LTO → microarchitecture → measurement** යන order එකෙන් බලන්න ඕන.

ඊළඟ level එකේ මේක තවත් deep කරන්න නම්, natural progression එක වන්නේ **`main.cpp → C++ compiler → generated `.obj`→ disassembly → CALL instruction → exact RSP values → shadow-space bytes → return address → linker relocation → PE/COFF symbol table → final`MainApp.exe`** කියන එක CPU instruction එකෙන් instruction එක trace කිරීමයි.
