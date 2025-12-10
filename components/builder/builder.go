// Package builder implements the build functionality
package builder

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

		// this is internal
		Libraries []string // consider if this is public or not
	}
}

type ProjectProperties struct {
	// @todo add support for pacgo packages
	Version        int      `json:"version"`
	Subdirectories []string `json:"subdirectories"`

	BuildToolchainPath string `json:"build_toolchain_path"`
	// @maybe move compiler to library level
	Compiler  string `json:"compiler"`
	Linker    string `json:"linker"`
	Assembler string `json:"assembler"`

	OutputPath               string `json:"output_path"`
	BuildMetadataEnable      bool   `json:"build_meta_data_enable"`
	PreprocessorOutputEnable bool   `json:"preprocessor_output_enable"`

	// a project is also a library with other libraries linked
	LibraryProperties
}

// @todo shall separate builder and project
var (
	ProjectConfiguration ProjectProperties // this HAS to be unique

	LibConfigurations = make(
		map[string]LibraryProperties,
	) // meant to hold the configurations instead of the array
)

func (proj *ProjectProperties) ResolveSubdirPaths(projectPath string) {
	for index, subdir := range proj.Subdirectories {
		proj.Subdirectories[index] = filepath.Join(projectPath, subdir)
	}
}

func (proj *ProjectProperties) ResolveOutputPath(projectPath string) {
	proj.OutputPath = filepath.Join(projectPath, proj.OutputPath)
}

func (lib *LibraryProperties) SpecifyNoMain() {
	(*lib).Flags = append((*lib).Flags, "-c")
}

// @todo maybe extract a common method, as this repeats the code a bit

func (lib *LibraryProperties) ResolvePrivateIncludesGlobalPaths(libraryPath string) {
	for index, include := range lib.Includes.Private {
		lib.Includes.Private[index] = filepath.Join(libraryPath, include)
	}
}

func (lib *LibraryProperties) ResolvePublicIncludesGlobalPaths(libraryPath string) {
	for index, include := range lib.Includes.Public {
		lib.Includes.Public[index] = filepath.Join(libraryPath, include)
	}
}

func (lib *LibraryProperties) ResolveSourcesGlobalPaths(libraryPath string) {
	for index, source := range lib.Sources {
		lib.Sources[index] = filepath.Join(libraryPath, source)
	}
}

// @todo the dependencies shall be checked

func (lib *LibraryProperties) ResolvePrivateDependencies() {
	for _, dependency := range lib.Dependencies.Private {
		libPath := filepath.Join(ProjectConfiguration.OutputPath, LibConfigurations[dependency].Name)

		lib.Includes.Private = append(lib.Includes.Private, LibConfigurations[dependency].Includes.Public...)
		// extremely dumb way of doing this. @todo remove it
		if !slices.Contains(lib.Flags, "-c") {
			lib.Dependencies.Libraries = append(lib.Dependencies.Libraries, libPath)
		}
	}
}

// treat library dependency as a graph - while not exactly a tree, can somewhat go through it like a tree
// this allows us to explore the build level for each node
// build level reflects the depth at which a certain library can be found

func (lib *LibraryProperties) ResolvePublicDependencies() {
	for _, dependency := range lib.Dependencies.Public {
		libPath := filepath.Join(ProjectConfiguration.OutputPath, LibConfigurations[dependency].Name)

		lib.Includes.Public = append(lib.Includes.Public, LibConfigurations[dependency].Includes.Public...)

		// extremely dumb way of doing this. @todo remove it
		if !slices.Contains(lib.Flags, "-c") {
			lib.Dependencies.Libraries = append(lib.Dependencies.Libraries, libPath)
		}
	}
}

func (lib *LibraryProperties) Build() error {
	setCompiler(*lib)

	setFlags(*lib)
	setDefines(*lib)
	setIncludes(*lib)

	setOutputProperties(*lib)
	setInputProperties(*lib)

	return runBuilder()
}
