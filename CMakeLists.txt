cmake_minimum_required(VERSION 3.0)
project(template)

### set output directory
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY output)

### generic defines 

### add component subdirectory to the build process
add_subdirectory(components)
### define target - main entry point for application
add_executable(${PROJECT_NAME}  app/main.c
)


### define includes for application

### main application defines

### set properties here
set_target_properties(${PROJECT_NAME} PROPERTIES LINKER_LANGUAGE C)

# link additional libraries needed
target_link_libraries(${PROJECT_NAME}
    lib.TestComponent
)
