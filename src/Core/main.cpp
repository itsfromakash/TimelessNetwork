#include <iostream>
#include "p2p_engine.h" // auto-generated header file from Go build

int main() {
    std::cout << "========================================" << std::endl;
    std::cout << "[C++ Core Engine] Starting Application..." << std::endl;
    std::cout << "========================================" << std::endl;

    // Go .a Static library එකේ export කරපු function එක call කිරීම
    StartP2PEngine();

    std::cout << "\n========================================" << std::endl;
    std::cout << "[C++ Core Engine] Application Terminated." << std::endl;
    std::cout << "========================================" << std::endl;

    return 0;
}