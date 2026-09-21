#include <iostream>
#include "core_engine.h"

int main() {
    std::cout << "[EXE] Main application started." << std::endl;
    
    int result = run_engine_pipeline(21);
    std::cout << "[EXE] Final result: " << result << std::endl;
    std::cout << "[EXE] Application finished successfully." << std::endl;
    return 0;
}