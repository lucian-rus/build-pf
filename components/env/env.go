package env

type LibraryProperties struct {
	Name    string   `json:"name"`
	Defines []string `json:"defines"`
	Sources []string `json:"sources"`
	Flags   []string `json:"flags"`

	Includes struct {
		Public  []string `json:"public"`  // these includes can be accessed by libs that depend on this lib
		Private []string `json:"private"` // these includes are not visible
	}

	Dependencies struct {
		Public    []string `json:"public"`  // allows dependencies to be inherited by other libs
		Private   []string `json:"private"` // dependencies will NOT be inherited
		Libraries []string // consider if this is public or not
	}
}

type ProjectProperties struct {
	// @todo add support for pacgo packages
	Version        int      `json:"version"`
	Subdirectories []string `json:"subdirectories"`

	BuildToolchainPath string `json:"build_toolchain_path"`
	Compiler           string `json:"compiler"`
	Linker             string `json:"linker"`
	Assembler          string `json:"assembler"`

	OutputPath         string `json:"output_path"`
	BuildMetadata      bool   `json:"build_meta_data"`
	PreprocessorOutput bool   `json:"preprocessor_output"`

	// a project is also a library with other libraries linked
	LibraryProperties
}

var (
	ProjectConfiguration   ProjectProperties // this HAS to be unique
	LibrariesConfiguration []LibraryProperties

	LibrariesMap = make(map[string]int) // maps the name of a library to its index in the `LibrariesConfiguration` array for fast access

	EnableDebugData = false
)
