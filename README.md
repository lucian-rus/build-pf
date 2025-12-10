## c/c++ builder

GoBI - go build infra

### building and running
* to build, run `go build` and copy the `gobi` executable to a global path
* in the c/c++ project directory run `gobi <command>` 

### in work

to be expected:
* build c/c++ files based on configuration files -> cmake like. uses json. microservices like apps
* build c/c++ by crawling the directory similar to block app
* define structures within json config files -> should allow a user to define where to put specific components in a env json config file.
* should allow a user to generate new components based on the paths
* should allow the user most of the functionality provided by cmake
* should work with pacgo
* should be very easy to configure
* should have a main json in which subfolders cand be added
* should support pre-build and post-build execution of custom scripts or apps
* should support custom commands 
* consider adding support for custom entry point -> for now, entry point is root dir and main should be here

* add support for cmake/make built files, objects and libs -> allow skip of compilation
* generator for .vscode files and linter support
* search and index all include files and resolve them dynamically before build-time -> this will work properly after include checks are implemented

### generator example
```
{
    "generator": {
        "components_config_path": "path_1"
    }
}
```