bits 64
default rel

section .text
global asm_add

; int asm_add(int a, int b)
asm_add:
    mov rax, rcx     ;
    add rax, rdx     ;
    ret              ;