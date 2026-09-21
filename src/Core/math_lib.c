#include "math_lib.h"

int multiply_by_two(int x) {
    // NASM මඟින් සාදන ලද asm_add එක භාවිතා කරයි
    return asm_add(x, x);
}