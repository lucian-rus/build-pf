package project

import (
	"gobi/modules/library"
	"path/filepath"
)

type ProjectProperties struct {
	// @todo add support for pacgo packages
	Version        int      `json:"version"`
	Subdirectories []string `json:"subdirectories"`

	BuildToolchainPath string `json:"build_toolchain_path"`
	// @maybe move compiler to library level
	Compiler  string `json:"compiler"`
	Linker    string `json:"linker"`
	Assembler string `json:"assembler"`

	OutputDirPath string `json:"output_dir"`
	// BuildMetadataEnable      bool   `json:"build_meta_data_enable"`
	// PreprocessorOutputEnable bool   `json:"preprocessor_output_enable"`

	// @todo enable user to support multiple kinds of header and source extensions
	// e.g hpp, cc, hh etc

	// @todo provide user with option to do variant handling
	// common for platform-type projects

	// @todo provide user with option to exclude specific directories ->
	// maybe make this mutually exlusive with the subdirectories field
	// this will enable the user to add example, vector etc directories in the repo that will not
	// be scanned by the crawler

	// a project is also a library
	library.LibraryProperties
}

func (proj *ProjectProperties) ResolveSubdirPaths(projectPath string) {
	for index, subdir := range proj.Subdirectories {
		proj.Subdirectories[index] = filepath.Join(projectPath, subdir)
	}
}

func (proj *ProjectProperties) ResolveOutputPath(projectPath string) {
	proj.OutputDirPath = filepath.Join(projectPath, proj.OutputDirPath)
}
