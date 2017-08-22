cmake_minimum_required(VERSION 3.8)
project(fgo_cgroup)

set(CMAKE_C_STANDARD 99)
link_libraries(libcgroup.so)
set(SOURCE_FILES main.c test_mount.c test_mount.h)
add_executable(fgo_cgroup ${SOURCE_FILES})