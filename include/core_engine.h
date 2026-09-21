#pragma once

#ifdef BUILDING_DLL
    #define ENGINE_API __declspec(dllexport)
#else
    #define ENGINE_API __declspec(dllimport)
#endif

int run_engine_pipeline(int input);
