cmake_minimum_required(VERSION 3.10)
project(tjc-serial-display)

# 设置默认构建类型为 Release
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release)
endif()

# 设置输出目录
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/release)

# 查找 Go 编译器
find_program(GO_EXECUTABLE go REQUIRED)

# 设置 Go 源文件
set(GO_SOURCE ${CMAKE_SOURCE_DIR}/cmd/tjc.go)

# 设置输出文件
set(OUTPUT_BINARY ${CMAKE_RUNTIME_OUTPUT_DIRECTORY}/tjc)

# 添加自定义目标来构建 Go 程序
add_custom_target(tjc ALL
    COMMAND ${CMAKE_COMMAND} -E make_directory ${CMAKE_RUNTIME_OUTPUT_DIRECTORY}
    COMMAND ${GO_EXECUTABLE} build -o ${OUTPUT_BINARY} ${GO_SOURCE}
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Building Go binary: tjc"
    VERBATIM
)

# 添加清理目标
add_custom_target(clean-all
    COMMAND ${CMAKE_COMMAND} -E remove_directory ${CMAKE_RUNTIME_OUTPUT_DIRECTORY}
    COMMAND ${GO_EXECUTABLE} clean
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Cleaning build artifacts"
)

# 添加安装规则
install(PROGRAMS ${OUTPUT_BINARY}
    DESTINATION bin
    COMPONENT runtime
)

# 显示构建信息
message(STATUS "Go executable: ${GO_EXECUTABLE}")
message(STATUS "Build type: ${CMAKE_BUILD_TYPE}")
message(STATUS "Output directory: ${CMAKE_RUNTIME_OUTPUT_DIRECTORY}")
message(STATUS "Output binary: ${OUTPUT_BINARY}")
