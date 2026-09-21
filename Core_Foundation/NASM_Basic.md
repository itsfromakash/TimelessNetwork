හරි Akash. මේක **Assembly ඉගෙනගන්න foundation එකේම foundation එක**. මේ දෙක හරියට තේරුණොත් C/C++ → Compiler → Assembly → CPU → RAM → OS/ABI කියන chain එක එකම picture එකක් වගේ පේන්න පටන් ගන්නවා. 🔥

මම **Grade 1 ළමයෙකුට කියන තරම් සරලව පටන්ගෙන, Grand Master level එකට** layer-by-layer ගෙන යන්නම්.

---

# 1. x86-64 CPU එකේ Fundamental Resources / ප්‍රධාන State එක මොනවාද?

මුලින්ම **CPU එක කියන්නේ මොකක්ද?**

CPU එක කියන්නේ:

> **Instructions execute කරන electronic machine එකක්.**

ඒකට program එකේ හැම දෙයක්ම "මතක" තියෙන්නේ නැහැ.

CPU එකට ඉතාම වැදගත් internal state එකක් තියෙනවා.

සරලව:

```text
CPU STATE
│
├── General-Purpose Registers
│
├── Instruction Pointer (RIP)
│
├── RFLAGS
│
├── Segment Registers
│
├── Control Registers
│
├── Debug Registers
│
├── SIMD / Floating-Point State
│
└── Memory-related architectural state
```

හැබැයි මේවා එකම මට්ටමේ දේවල් නෙවෙයි.

අපි එකින් එක බලමු.

---

# 2. මුලින්ම "State" කියන්නේ මොකක්ද?

මේක **ඉතාම වැදගත් concept එකක්.**

Imagine කරන්න:

ඔයා chess game එකක් play කරනවා.

Board එකේ:

```text
King → E4
Queen → D5
Pawn → F2
...
```

මේ වෙලාවේ board එකේ තියෙන arrangement එක තමයි **state**.

CPU එකටත් ඒ වගේම current condition එකක් තියෙනවා.

උදාහරණයක්:

```text
RAX = 25
RBX = 10
RIP = 0x401020
RFLAGS = ...
```

මේ වෙලාවේ CPU එකේ current condition එක මෙන්න මේක.

ඒ කියන්නේ:

> **CPU state = CPU එකේ execution එක continue කරන්න අවශ්‍ය current architectural information.**

---

# 3. CPU එකේ ප්‍රධානම resource එක — Registers

මෙතනින් Assembly ලෝකය පටන් ගන්නවා.

**Register** කියන්නේ CPU ඇතුළේ තියෙන ඉතාම වේගවත් storage location එකක්.

RAM:

```text
CPU
 ↓
RAM
```

Register:

```text
CPU
 ┌─────────────────┐
 │ RAX = 123       │
 │ RBX = 456       │
 │ RCX = 789       │
 └─────────────────┘
```

Registers CPU එකේම architectural state එකේ කොටසක්.

---

# 4. x86-64 General Purpose Registers

x86-64 වල ප්‍රධාන 16 General-Purpose Registers තියෙනවා.

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

මේ හැම එකකටම **64-bit** width එකක් තියෙනවා.

```text
64 bits = 8 bytes
```

ඒ නිසා:

```text
RAX
└── 64 bits
```

---

# 5. Register එකක් කියන්නේ variable එකක්ද?

මෙතනින් තමයි Assembly වලට ගොඩක් beginners confuse වෙන්නේ.

**නැහැ.**

Register එකක් variable එකක් නෙවෙයි.

Register එකක් කියන්නේ CPU එකේ **storage location** එකක්.

Programming language එකේ:

```c
int x = 10;
```

මෙතන:

```text
x
```

කියන්නේ source-level variable එකක්.

Assembly එකේ:

```asm
mov rax, 10
```

මෙතන:

```text
RAX = 10
```

CPU state එකේ register එකක් value එකක් hold කරනවා.

---

# 6. RIP — Instruction Pointer

මේක **ඉතාම වැදගත්**.

x86-64 CPU එකේ:

```text
RIP
```

= **Instruction Pointer Register**

එය next instruction එකේ address එක track කරනවා.

Imagine:

```asm
mov rax, 10
add rax, 20
sub rax, 5
```

Memory එකේ instructions තිබෙනවා:

```text
0x401000   mov rax, 10
0x401007   add rax, 20
0x40100A   sub rax, 5
```

CPU එක:

```text
RIP = 0x401000
```

එතකොට CPU එක කියනවා:

> "මම දැන් 0x401000 address එකේ instruction එක execute කරන්න ඕන."

ඒ instruction execute කරලා CPU එක RIP update කරනවා.

```text
RIP = next instruction address
```

ඒ නිසා program execution එකේ basic idea එක:

```text
RIP
 ↓
Fetch instruction
 ↓
Decode
 ↓
Execute
 ↓
Update RIP
 ↓
Next instruction
```

---

# 7. RFLAGS — CPU එකේ condition information

තවත් major state එකක්:

```text
RFLAGS
```

මේක CPU එකේ **status/condition flags** තියෙන register එක.

උදාහරණ:

```text
ZF = Zero Flag
CF = Carry Flag
SF = Sign Flag
OF = Overflow Flag
```

උදාහරණයක්:

```asm
mov rax, 10
sub rax, 10
```

Result:

```text
RAX = 0
```

CPU එකට:

```text
ZF = 1
```

වෙන්න පුළුවන්.

ඊළඟට:

```asm
jz somewhere
```

`jz` = Jump if Zero.

CPU එක ZF බලනවා.

```text
ZF == 1?
   │
   ├── YES → jump
   │
   └── NO  → continue
```

---

# 8. RSP — Stack Pointer

මේක Assembly + C + OS + ABI වල **super important**.

```text
RSP
```

= Stack Pointer

Stack කියන්නේ memory එකේ special usage pattern එකක්.

Conceptually:

```text
Higher Address
────────────────
│              │
│   stack      │
│              │
├──────────────┤
│      ↓       │
│     RSP      │
└──────────────┘
Lower Address
```

Function calls වලදී:

```text
call
push
pop
ret
```

වගේ instructions එක්ක RSP ගොඩක් use වෙනවා.

---

# 9. RBP — Base Pointer / Frame Pointer

```text
RBP
```

සාමාන්‍යයෙන් stack frame එකේ reference point එකක් ලෙස භාවිතා කළ හැක.

උදාහරණයක්:

```text
Function stack frame

Higher address
────────────────
│ argument      │
│ return addr   │
│ old RBP       │ ← RBP reference
│ local var     │
│ local var     │
────────────────
Lower address
```

හැබැයි modern compilers **RBP හැම වෙලාවෙම frame pointer ලෙස භාවිතා කරන්නේ නැහැ**.

Optimization නිසා RBP එක general-purpose register එකක් ලෙසත් භාවිතා කළ හැක.

---

# 10. SIMD / Floating Point State

x86-64 CPU එකේ integer registers විතරක් නෙවෙයි.

තවත් register state එකක් තියෙනවා:

```text
XMM
YMM
ZMM
```

උදාහරණ:

```text
XMM0
XMM1
...
```

මේවා:

* floating-point
* SIMD
* vector operations

සඳහා භාවිතා වෙනවා.

උදාහරණයක්:

```text
4 numbers එකවර process කිරීම
```

වගේ දේවල්.

---

# 11. Control Registers

තවත් deep layer එකක්:

```text
CR0
CR2
CR3
CR4
...
```

මේවා normal application variables සඳහා නෙවෙයි.

Operating system / CPU control සඳහා වැදගත්.

විශේෂයෙන්:

```text
CR3
```

page-table / address-translation context එක සමඟ ඉතා වැදගත්.

ඒ නිසා:

```text
User program
   ↓
General registers
   ↓
OS / kernel
   ↓
Control registers
   ↓
Memory management hardware
```

වගේ deeper relationship එකක් තියෙනවා.

---

# 12. Segment Registers

x86 architecture එකේ historical + architectural state එකේ කොටසක්:

```text
CS
DS
ES
FS
GS
SS
```

Modern x86-64 programming වල segmentation බොහෝ දුරට simplified වෙලා තිබුණත්:

```text
FS
GS
```

විශේෂයෙන් OS/runtime implementations වල වැදගත් වෙන්න පුළුවන්.

---

# 13. දැන් CPU state එක එක picture එකකට දාමු

```text
                x86-64 CPU
┌───────────────────────────────────────────┐
│                                           │
│ General Purpose Registers                 │
│                                           │
│ RAX RBX RCX RDX RSI RDI                  │
│ RBP RSP R8 R9 R10 R11 R12 R13 R14 R15   │
│                                           │
│ Instruction Pointer                       │
│ RIP                                       │
│                                           │
│ Status / Control                          │
│ RFLAGS                                    │
│                                           │
│ Segment State                              │
│ CS DS ES FS GS SS                         │
│                                           │
│ Control Registers                          │
│ CR0 CR2 CR3 CR4 ...                       │
│                                           │
│ Debug Registers                            │
│ DR0 ...                                   │
│                                           │
│ Floating / SIMD State                     │
│ XMM / YMM / ZMM                           │
│                                           │
└───────────────────────────────────────────┘
```

මේ තමයි "CPU එකේ state" කියන concept එකේ broad picture එක.

---

# 14. දැන් දෙවැනි ප්‍රශ්නය — Variable එක Assembly වලට map වෙන්නේ කොහොමද?

මෙතනින් **Compiler Design** පටන් ගන්නවා. 🔥

Source code එකක් ගමු.

```c
int x = 10;
```

Beginner කෙනෙක් හිතන විදිහ:

```text
x → RAM එකේ 10
```

මේක **හැමවිටම true නෑ.**

ඇත්තටම compiler එකට freedom තියෙනවා.

`x` කියන source-level variable එක:

```text
Register
        හෝ
Stack
        හෝ
Memory
        හෝ
Optimized away
```

වෙන්න පුළුවන්.

---

# 15. Source Variable කියන්නේ මොකක්ද?

C code:

```c
int x = 10;
```

මෙතන `x` කියන්නේ:

> Programmer විසින් value එකකට දුන්න symbolic name එක.

CPU එකට:

```text
"x" කියලා variable එකක්
```

කියලා directly තේරෙන්නේ නැහැ.

CPU එක දන්නේ:

```text
registers
memory addresses
instructions
flags
```

වගේ low-level things.

---

# 16. මේ transformation එක තේරුම් ගන්න

```text
C source code

int x = 10;
```

↓

```text
Compiler / Optimizer

"මේ x එක කොහෙ තියාගන්නද?"
```

↓

Possible:

```text
RAX
```

↓

Assembly:

```asm
mov eax, 10
```

↓

CPU:

```text
EAX = 10
```

---

# 17. `int x = 10` → `EAX`

උදාහරණයක්:

```c
int x = 10;
int y = x + 5;
```

Compiler එකට මෙහෙම generate කරන්න පුළුවන්:

```asm
mov eax, 10
add eax, 5
```

Execution:

```text
Initially:

EAX = ??????
```

First:

```asm
mov eax, 10
```

Result:

```text
EAX = 10
```

Then:

```asm
add eax, 5
```

Result:

```text
EAX = 15
```

මෙතන:

```text
x
```

සහ

```text
y
```

CPU එකේ actual named objects වශයෙන් තිබුණේ නැහැ.

Compiler එක ඒ concepts registers වලට map කළා.

---

# 18. එතකොට `x` තියෙන්නේ කොහෙද?

මෙන්න **Grand Master question එක.**

Answer:

> **එක fixed location එකක තියෙනවා කියලා කියන්න බැහැ.**

Compiler optimization අනුව:

```text
x
│
├── RAX
├── RBX
├── stack memory
├── another memory location
└── nowhere
```

---

# 19. "Nowhere" කියන්නේ කොහොමද?

උදාහරණයක්:

```c
int x = 10;
return x;
```

Compiler එකට:

```asm
mov eax, 10
ret
```

වගේ generate කරන්න පුළුවන්.

`x` කියන value එක වෙනම location එකකට store කරන්න අවශ්‍ය නැහැ.

Compiler එක direct result එක generate කරනවා.

ඒ නිසා:

```text
Source:

x = 10
return x
```

machine level එකේ conceptually:

```text
return 10
```

වගේ වෙන්න පුළුවන්.

---

# 20. Variable → Register

ඉතා common mapping එක:

```text
C variable
     ↓
Compiler
     ↓
Register
```

Example:

```c
int a = 10;
int b = 20;
int c = a + b;
```

Conceptually:

```text
a → EAX
b → ECX
c → EDX
```

Assembly:

```asm
mov eax, 10
mov ecx, 20
add eax, ecx
```

Result:

```text
EAX = 30
```

Compiler එක මේ register allocation එක decide කරනවා.

---

# 21. Variable → Stack

Register එකක් හැම variable එකකටම දෙන්න බැහැ.

ඒකට හේතු:

* registers limited
* function calls
* register pressure
* ABI requirements
* address-taking
* lifetime
* optimization decisions

එතකොට variable එක stack එකේ තිබෙන්න පුළුවන්.

Example conceptual assembly:

```asm
mov dword [rbp-4], 10
```

මෙතන:

```text
[rbp-4]
```

කියන්නේ memory address එකක්.

Conceptually:

```text
RBP = 0x1000

RBP - 4
 ↓
0x0FFC
```

එතකොට:

```text
memory[0x0FFC] = 10
```

---

# 22. `[rbp-4]` කියන්නේ මොකක්ද?

මේක ඉතාම වැදගත් Assembly syntax එකක්.

```asm
mov dword [rbp-4], 10
```

`rbp`:

```text
register
```

`rbp - 4`:

```text
address calculation
```

`[ ... ]`:

```text
memory dereference
```

ඒ නිසා:

```asm
[rbp-4]
```

කියන්නේ:

> RBP register එකේ value එකෙන් 4 bytes අඩු කරලා ලැබෙන memory address එකේ තිබෙන data.

---

# 23. `[ ]` නැතිනම්?

```asm
mov rax, rbp
```

මෙහි:

```text
RAX = RBP
```

එනම් register → register.

නමුත්:

```asm
mov rax, [rbp]
```

මෙහි:

```text
RAX = memory[RBP]
```

ඒ කියන්නේ:

```text
RBP
 ↓
address
 ↓
RAM
 ↓
value
 ↓
RAX
```

**මේ distinction එක Assembly වල fundamental.**

---

# 24. Variable → Memory

Source:

```c
int x = 42;
```

Compiler එකට stack variable එකක් අවශ්‍ය නම් conceptually:

```asm
mov dword [rbp-4], 42
```

Memory:

```text
Stack

Address
0x1000
────────────────
0x0FFC → 42
────────────────
```

`x` කියන්නේ source language එකේ name එක.

Machine level එකේ:

```text
x
 ↓
stack offset
 ↓
memory address
 ↓
bytes
```

---

# 25. Variable එකේ datatype එකත් disappear වෙන්න පුළුවන්

C:

```c
int x = 42;
```

`int` කියන්නේ source-level type information.

CPU එකට:

```text
"මෙන්න int එකක්"
```

කියලා object/type system එකක් නැහැ.

CPU එකට basically:

```text
bits
```

තියෙනවා.

ඒ නිසා:

```text
int
float
pointer
struct
```

වගේ source-level concepts compiler එක machine operations වලට translate කරනවා.

---

# 26. `int` එකේ size එක Assembly එකේ කොහොමද දන්නේ?

Instruction operand size එකෙන්.

```asm
mov eax, 42
```

`EAX` = 32-bit.

ඒ නිසා 32-bit operation.

Memory case එක:

```asm
mov dword [rbp-4], 42
```

`dword` = 32-bit = 4 bytes.

```text
byte  = 8 bits
word  = 16 bits
dword = 32 bits
qword = 64 bits
```

---

# 27. Variable → Register → Memory

Variable එක execution එක අතරතුර move වෙන්නත් පුළුවන්.

උදාහරණයක්:

```text
C variable x
      ↓
Stack
      ↓
load
      ↓
Register
      ↓
calculation
      ↓
Register
      ↓
store
      ↓
Stack
```

Conceptually:

```asm
mov eax, [rbp-4]
add eax, 10
mov [rbp-4], eax
```

මෙතන:

```text
memory → register → computation → memory
```

---

# 28. මේක CPU එකේදී ඇත්තටම සිදුවෙන්නේ?

Suppose:

```asm
mov eax, [rbp-4]
```

CPU execution එකේ high-level conceptual flow:

```text
1. Decode instruction
       ↓
2. Identify source:
       [RBP - 4]
       ↓
3. Calculate effective address
       ↓
4. Access memory hierarchy
       ↓
5. Obtain 32-bit value
       ↓
6. Put value into EAX
```

හැබැයි modern CPU එක මේක **literally sequential simple steps** ලෙස execute කරනවා කියලා හිතන්න එපා.

Modern x86 CPU:

```text
Fetch
Decode
Rename
Dispatch
Schedule
Execute
Retire
```

වගේ complex out-of-order machinery භාවිතා කරනවා.

ඒක පස්සේ අපි deep dive කරන්න ඕන topic එකක්.

---

# 29. Local variable එකක් සහ global variable එකක්

මේ දෙකත් වෙනස්.

### Local

```c
void f() {
    int x = 10;
}
```

`x` සාමාන්‍යයෙන් function scope එකට සම්බන්ධ.

Compiler එකට:

```text
register
or
stack
```

දෙන්න පුළුවන්.

---

### Global

```c
int x = 10;
```

මේක static storage duration එකක් තියෙන object එකක්.

සාමාන්‍යයෙන් executable/data image එකේ memory location එකකට map වෙන්න පුළුවන්.

Conceptually:

```text
x
 ↓
symbol
 ↓
address
 ↓
memory
```

---

# 30. Pointer variable එකක්?

මෙන්න තවත් important example එකක්.

```c
int x = 100;
int *p = &x;
```

මෙතන:

```text
x = 100
```

`p`:

> x තියෙන memory address එක store කරන variable එක.

Conceptually:

```text
x
↓
memory address 0x2000

p
↓
0x2000
```

Assembly concept:

```asm
mov eax, 100
mov [somewhere], eax
lea rax, [somewhere]
```

මෙතන `lea` කියන්නේ **Load Effective Address**.

---

# 31. Variable vs Address vs Value

මේ තුන වෙනස්.

```text
Variable
   │
   ├── name
   ├── type
   ├── value
   └── storage/location
```

Example:

```c
int x = 42;
```

Conceptually:

```text
name    = x
type    = int
value   = 42
location = some storage
```

Assembly/CPU level එකට යනකොට:

```text
name
 ↓
symbol/debug/compiler concept

type
 ↓
compiler information / instruction width

value
 ↓
bits

location
 ↓
register or memory
```

---

# 32. මෙන්න Assembly හි Grand Master mental model එක

**Source code variable එක CPU register එකක් නෙවෙයි.**

ඒක හොඳට මතක තියාගන්න.

```text
                    SOURCE LEVEL
                ┌──────────────────┐
                │ int x = 42;      │
                └────────┬─────────┘
                         │
                         │ Compiler
                         ↓
                ┌──────────────────┐
                │ Variable x       │
                │ is mapped to...  │
                └────────┬─────────┘
                         │
              ┌──────────┼──────────┐
              ↓          ↓          ↓
           Register     Stack      Memory
              │          │          │
              ↓          ↓          ↓
            EAX       [RBP-4]     [address]
              │          │          │
              └──────────┼──────────┘
                         ↓
                    MACHINE LEVEL
                         ↓
                       CPU
```

---

# 33. තවත් deeper concept එක — Compiler එක variable එකකට "home" එකක් දෙනවාද?

Optimization නැති situation එකක:

```text
x → [RBP-4]
```

වගේ stable location එකක් තියෙන්න පුළුවන්.

නමුත් optimization තියෙනකොට:

```text
x
│
├── RAX
│
├── RCX
│
├── stack
│
└── optimized away
```

වගේ වෙනස් වෙන්න පුළුවන්.

එකම function එකේ different points වල:

```text
x → RAX

...

x → stack

...

x → RCX
```

වෙන්න පුළුවන්.

මේකට compiler concept එකෙන් **register allocation** සහ **liveness** වගේ concepts සම්බන්ධයි.

---

# 34. Debugger එකේ variable එක පේන්නේ ඇයි?

ඔයා C program එක debug කරනවා:

```c
int x = 10;
```

Debugger එකේ:

```text
x = 10
```

කියලා පෙන්වනවා.

CPU එක ඇත්තටම:

```text
"x"
```

කියලා දන්නේ නැහැ.

Compiler එක debugging information generate කළොත්:

```text
source variable x
       ↓
debug information
       ↓
register / stack / location
```

Debugger එකට programmer-friendly view එකක් හදන්න පුළුවන්.

ඒ නිසා:

> **Debugger variable view ≠ CPU's native understanding of variables.**

---

# 35. Assembly එකේ `x` කියලා variable එකක් declare කළොත්?

NASM වගේ assembler එකේ:

```asm
section .data

x dd 42
```

මෙතන `x` කියන්නේ C-style runtime variable abstraction එකක් නෙවෙයි.

`x` කියන්නේ assembler symbol එකක්.

```text
x
 ↓
address/symbol
 ↓
memory location
```

`dd`:

```text
Define Doubleword
```

එනම් 4-byte data.

ඒ නිසා:

```asm
x dd 42
```

conceptually:

```text
memory:
+---------+
| 42      |
+---------+
   ↑
   x
```

---

# 36. NASM එකේ register vs memory symbol

```asm
mov eax, 42
```

මෙහි:

```text
EAX = 42
```

නමුත්:

```asm
mov eax, [x]
```

මෙහි:

```text
EAX = memory[x]
```

ඒ නිසා brackets:

```asm
[x]
```

ඉතා වැදගත්.

Compare:

```asm
mov eax, x
```

vs

```asm
mov eax, [x]
```

Conceptually:

```text
eax = address/value represented by x
```

vs

```text
eax = contents stored at x
```

NASM syntax/context අනුව exact relocation/address encoding වෙනස් විය හැක, නමුත් beginner mental model එකට මේ distinction එක **address vs contents** ලෙස තියාගන්න.

---

# 37. CPU එක දන්නේ Variable Names ද?

**නැහැ.**

CPU level එකේ:

```text
x
y
total
price
studentName
```

වගේ names meaningful නෑ.

CPU එකට තියෙන්නේ:

```text
instruction encoding
register numbers
immediate values
memory addresses / effective addresses
flags
```

වගේ machine-level information.

---

# 38. ඒකෙන් compiler එකේ job එක තේරෙනවා

Compiler එක basically bridge එකක්:

```text
Human-friendly language
        ↓
C / C++
        ↓
Compiler
        ↓
IR / Optimization
        ↓
Machine instructions
        ↓
Assembly / Object code
        ↓
CPU
```

Variable:

```text
x
```

මේ pipeline එකේ gradually transform වෙනවා:

```text
Source variable
      ↓
Compiler IR value
      ↓
Virtual register / SSA value
      ↓
Physical register OR memory
      ↓
Machine-level bits
```

**මේක compiler architecture එකේ ඉතාම powerful mental model එකක්.**

---

# 39. තවත් deep layer — SSA

Compiler intermediate representation එකේ variable concept එක source variable එකට වඩා වෙනස්.

උදාහරණයක්:

```c
x = 10;
x = x + 5;
```

Compiler IR එකේ conceptually:

```text
x1 = 10
x2 = x1 + 5
```

ඒක **SSA — Static Single Assignment** style representation එකක්.

ඉන්පසු register allocation:

```text
x1 → EAX
x2 → EAX
```

හෝ වෙනත් registers/memory.

ඒ නිසා source variable:

```text
x
```

machine register එකකට **directly equal** කියන idea එක වැරදියි.

---

# 40. Final mental picture 🧠

මේක හොඳට memorize කරන්න:

```text
┌─────────────────────────────────────────────┐
│              PROGRAMMER LEVEL               │
│                                             │
│        int x = 42;                          │
│                                             │
│        "x" = symbolic variable              │
└─────────────────────┬───────────────────────┘
                      │
                      ↓
┌─────────────────────────────────────────────┐
│               COMPILER LEVEL                │
│                                             │
│   Type / Lifetime / Scope / Liveness        │
│   Optimization / Register Allocation        │
└─────────────────────┬───────────────────────┘
                      │
             ┌────────┴─────────┐
             ↓                  ↓
        Register             Memory
             │                  │
          EAX/RAX          [RBP-4], etc.
             │                  │
             └────────┬─────────┘
                      ↓
┌─────────────────────────────────────────────┐
│                CPU LEVEL                    │
│                                             │
│  Registers + RIP + RFLAGS + memory access   │
│                                             │
│       Everything ultimately becomes         │
│              bits + operations              │
└─────────────────────────────────────────────┘
```

## 🔥 එක වාක්‍යයකින්

**1. x86-64 CPU එකේ fundamental architectural state** කියන්නේ ප්‍රධාන වශයෙන් **general-purpose registers, RIP, RFLAGS, SIMD/FPU state, segment/control/debug registers** වැනි CPU execution state එකයි.

**2. Source-level variable එකක් CPU එකේ "variable" ලෙස exist වෙන්නේ නැහැ.** Compiler එක ඒ variable/value එකේ lifetime, usage, optimization සහ ABI requirements අනුව **register එකකට, stack/memory location එකකට map කරන්න හෝ completely eliminate කරන්න** පුළුවන්.

---

### මේ foundation එකෙන් පස්සේ ඊළඟට තේරුම් ගන්න ඕන critical chain එක

```text
Variable
   ↓
Data Type
   ↓
Storage
   ↓
Register / Memory
   ↓
Virtual Address
   ↓
Physical Address
   ↓
Cache
   ↓
RAM
   ↓
CPU Load/Store
   ↓
Microarchitecture
   ↓
Transistors
```

**මෙතනින් "Assembly language" ඇත්තටම CPU එකට සම්බන්ධ වෙන්නේ කොහොමද කියන deepest layer එක පටන් ගන්නවා.**
