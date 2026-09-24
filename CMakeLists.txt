cmake_minimum_required(VERSION 3.15)
project(TimelessNetwork CXX)

set(CMAKE_CXX_STANDARD 17)

# Headers search path
include_directories(${CMAKE_SOURCE_DIR}/include)

# Go Static Library Path
set(GO_P2P_LIB ${CMAKE_SOURCE_DIR}/lib/p2p_engine.a)

# C++ Executable Target (src/Core/main.cpp)
add_executable(TimelessNetworkApp src/Core/main.cpp)

# Link Go Library to C++ Executable
target_link_libraries(TimelessNetworkApp PRIVATE ${GO_P2P_LIB})

# Windows System Networking Dependencies for WebRTC
if(WIN32)
    target_link_libraries(TimelessNetworkApp PRIVATE ws2_32 winmm iphlpapi userenv ntdll pthread)
endif()