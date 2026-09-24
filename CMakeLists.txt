cmake_minimum_required(VERSION 3.15)
project(TimelessNetwork CXX)

set(CMAKE_CXX_STANDARD 17)

# Include path and library setup
include_directories(${CMAKE_SOURCE_DIR}/include)
set(GO_P2P_LIB ${CMAKE_SOURCE_DIR}/lib/p2p_engine.a)

# Executable target
add_executable(TimelessNetworkApp src/Core/main.cpp)

# Link Go Static Library
target_link_libraries(TimelessNetworkApp PRIVATE ${GO_P2P_LIB})

# Windows System Libraries for WebRTC Networking & Go Runtime
if(WIN32)
    target_link_libraries(TimelessNetworkApp PRIVATE ws2_32 winmm iphlpapi userenv ntdll pthread)
endif()