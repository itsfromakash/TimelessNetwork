#include "core_engine.h"
#include "math_lib.h"
#include <iostream>

int run_engine_pipeline(int input) {
        std::cout << "Core Engine initialized successfully." << std::endl;
        int result = multiply_by_two(input);
        return result;
}