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

	OutputPath               string `json:"output_path"`
	BuildMetadataEnable      bool   `json:"build_meta_data_enable"`
	PreprocessorOutputEnable bool   `json:"preprocessor_output_enable"`

	// a project is also a library
	library.LibraryProperties
}

func (proj *ProjectProperties) ResolveSubdirPaths(projectPath string) {
	for index, subdir := range proj.Subdirectories {
		proj.Subdirectories[index] = filepath.Join(projectPath, subdir)
	}
}

func (proj *ProjectProperties) ResolveOutputPath(projectPath string) {
	proj.OutputPath = filepath.Join(projectPath, proj.OutputPath)
}
