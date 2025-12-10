package parser

// @todo this package reads headers/source files and determines the includes.
// the includes shall be searched by the crawler and written as configurations
// this can help with both the generator and automation of the build process

// @todo long-term
// should add support for custom include support (e.g includes are done via macro)
// should parse macros and check whether or not conditionally compiled includes will be resolved
// should get access to build internals, e.g what each library has so it can automatically append a specific lib as dependency
