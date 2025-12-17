package library

import (
	"path/filepath"
	"slices"
)

type LibraryProperties struct {
	Name    string   `json:"name"`
	Defines []string `json:"defines"`
	Sources []string `json:"sources"`
	Flags   []string `json:"flags"`

	// project inheritance
	InheritFlags   bool `json:"inherit_flags"`
	InheritDefines bool `json:"inherit_defines"`

	Includes struct {
		Public  []string `json:"public"`  // these includes can be accessed by libs that depend on this lib
		Private []string `json:"private"` // these includes are not visible
	}

	Dependencies struct {
		Public  []string `json:"public"`  // allows dependencies to be inherited by other libs
		Private []string `json:"private"` // dependencies will NOT be inherited
	}

	// internals -> not meant to be configured via json
	Root          string
	LinkedObjects []string
	Headers       []string
}

func (lib *LibraryProperties) SpecifyNoMain() {
	(*lib).Flags = append((*lib).Flags, "-c")
}

func (lib *LibraryProperties) ResolvePrivateIncludesGlobalPaths() {
	for index, include := range lib.Includes.Private {
		lib.Includes.Private[index] = filepath.Join(lib.Root, include)
	}
}

func (lib *LibraryProperties) ResolvePublicIncludesGlobalPaths() {
	for index, include := range lib.Includes.Public {
		lib.Includes.Public[index] = filepath.Join(lib.Root, include)
	}
}

func (lib *LibraryProperties) ResolveSourcesGlobalPaths() {
	for index, source := range lib.Sources {
		lib.Sources[index] = filepath.Join(lib.Root, source)
	}
}

// @todo the dependencies shall be checked

func (lib *LibraryProperties) ResolvePrivateDependencies(buildDir string, libConfigMap map[string]LibraryProperties) {
	for _, dependency := range lib.Dependencies.Private {
		libPath := filepath.Join(buildDir, libConfigMap[dependency].Name)

		lib.Includes.Private = append(lib.Includes.Private, libConfigMap[dependency].Includes.Public...)
		// extremely dumb way of doing this. @todo remove it
		if !slices.Contains(lib.Flags, "-c") {
			lib.LinkedObjects = append(lib.LinkedObjects, libPath)
		}
	}
}

// treat library dependency as a graph - while not exactly a tree, can somewhat go through it like a tree
// this allows us to explore the build level for each node
// build level reflects the depth at which a certain library can be found

func (lib *LibraryProperties) ResolvePublicDependencies(buildDir string, libConfigMap map[string]LibraryProperties) {
	for _, dependency := range lib.Dependencies.Public {
		libPath := filepath.Join(buildDir, libConfigMap[dependency].Name)

		lib.Includes.Public = append(lib.Includes.Public, libConfigMap[dependency].Includes.Public...)

		// extremely dumb way of doing this. @todo remove it
		if !slices.Contains(lib.Flags, "-c") {
			lib.LinkedObjects = append(lib.LinkedObjects, libPath)
		}
	}
}
