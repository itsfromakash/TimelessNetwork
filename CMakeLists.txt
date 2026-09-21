cmake_minimum_required(VERSION 3.20)
project(NativePipelineProject C CXX ASM_NASM)

set(CMAKE_C_STANDARD 11)
set(CMAKE_CXX_STANDARD 17)

# NASM Assembly Language Enable කිරීම (Windows MSVC සඳහා win64 Target කිරීම)
enable_language(ASM_NASM)
set(CMAKE_ASM_NASM_FLAGS "-f win64")

# 1. Static Library (.lib) සැකසීම (NASM + C)
add_library(MathStaticLib STATIC
    src/asm_func.asm
    src/math_lib.c
)
target_include_directories(MathStaticLib PUBLIC include)

# 2. Dynamic Library (.dll) සැකසීම (C++ -> Static Library වෙත Link වේ)
# add_library(CoreEngineDLL SHARED
#     src/core_engine.cpp
# )

# target_compile_definitions(CoreEngineDLL PRIVATE BUILDING_DLL)
# target_include_directories(CoreEngineDLL PUBLIC include)
# target_link_libraries(CoreEngineDLL PRIVATE MathStaticLib)

# 3. Executable (.exe) සැකසීම (Main App -> Dynamic Library වෙත Link වේ)
add_executable(MainApp
    src/main.cpp
    src/core_engine.cpp

)
target_include_directories(MainApp PRIVATE include)
target_link_libraries(MainApp PRIVATE MathStaticLib)