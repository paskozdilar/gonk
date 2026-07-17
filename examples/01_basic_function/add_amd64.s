#include "textflag.h"

// func callCAdd(addr uintptr, a, b int32) int32
TEXT ·callCAdd(SB), NOSPLIT, $0
    // Fetch arguments from the stack using the virtual Frame Pointer (FP)
    MOVQ addr+0(FP), AX    // Load the C function address into AX
    MOVL a+8(FP), DI       // Load 'a' (32-bit int) into C's 1st register (DI)
    MOVL b+12(FP), SI      // Load 'b' (32-bit int) into C's 2nd register (SI)
    
    CALL AX                // Safely jump to the C function
    
    // C leaves the 32-bit return value in AX.
    // Go expects its return value at the next stack slot after the arguments.
    MOVL AX, ret+16(FP)    // Move the result to Go's return slot
    RET
