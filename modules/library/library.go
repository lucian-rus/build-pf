package library

import (
	"path/filepath"
)

type LibraryType uint8

const (
	TypeProject LibraryType = iota
	TypeInterface
	TypeLibrary
)

type ScopedProperties struct {
	Public  []string `json:"public,omitempty"`
	Private []string `json:"private,omitempty"`
}

type LibraryProperties struct {
	Name    string   `json:"name"`
	Sources []string `json:"sources"`

	// project inheritance
	InheritFlags   bool `json:"inherit_flags"`
	InheritDefines bool `json:"inherit_defines"`

	Includes     ScopedProperties `json:"includes"`
	Defines      ScopedProperties `json:"defines"`
	Flags        ScopedProperties `json:"flags"`
	LinkerFlags  ScopedProperties `json:"linker_flags"`
	Dependencies ScopedProperties `json:"dependencies"`

	// internals -> not meant to be configured via json
	Type       LibraryType
	Root       string
	OutputPath string
}

func (lib *LibraryProperties) SetDefaultValues() {
	(*lib).InheritDefines = true
	(*lib).InheritFlags = true
}

func (lib *LibraryProperties) SpecifyNoMain() {
	(*lib).Flags.Private = append((*lib).Flags.Private, "-c")
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
		// @todo check if this is fine
		lib.Includes.Private = append(lib.Includes.Private, libConfigMap[dependency].Includes.Public...)
	}
}

// treat library dependency as a graph - while not exactly a tree, can somewhat go through it like a tree
// this allows us to explore the build level for each node
// build level reflects the depth at which a certain library can be found

func (lib *LibraryProperties) ResolvePublicDependencies(buildDir string, libConfigMap map[string]LibraryProperties) {
	for _, dependency := range lib.Dependencies.Public {
		lib.Includes.Public = append(lib.Includes.Public, libConfigMap[dependency].Includes.Public...)
	}
}

func (lib *LibraryProperties) ResolveOutputPath(buildDir string) {
	lib.OutputPath = filepath.Join(buildDir, "libs", lib.Name)
}

func (lib *LibraryProperties) InheritProjectFlags(projectConfig LibraryProperties) {
	if !lib.InheritFlags {
		return
	}
}

func (lib *LibraryProperties) InheritProjectDefines(projectConfig LibraryProperties) {
	if !lib.InheritDefines {
		return
	}

	lib.Defines.Public = append(lib.Defines.Public, projectConfig.Defines.Public...)
	lib.Defines.Private = append(lib.Defines.Private, projectConfig.Defines.Private...)
}
