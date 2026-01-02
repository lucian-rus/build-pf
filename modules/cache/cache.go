package cache

// this is the base cache structure. rest are inherited and can be extended
type FileCache struct {
	Path      string
	Timestamp int
}

type SourceCache struct {
	Library string

	FileCache
}

type HeaderCache struct {
	Library string

	FileCache
}

type BuildCache struct {
	FileCache
}
