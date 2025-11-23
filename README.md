## c/c++ builder

GoBI - go build infra

to be expected:
* build c/c++ files based on configuration files -> cmake like. uses json. microservices like apps
* build c/c++ by crawling the directory similar to block app
* define structures within json config files -> should allow a user to define where to put specific components in a env json config file.
e.g : 
```
{
    "generator": {
        "components_config_path": "path_1"
    }
}
```
* should allow a user to generate new components based on the paths
* should allow the user most of the functionality provided by cmake
* should work with pacgo
* should be very easy to configure
* should have a main json in which subfolders cand be added
* should support pre-build and post-build execution of custom scripts or apps
* should support custom commands 
* consider adding support for custom entry point -> for now, entry point is root dir and main should be here
