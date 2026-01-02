package cache

// this is the base cache structure. rest are inherited and can be extended
type FileCache struct {
	Path      string
	Timestamp int
}

type SourceCache struct {
	FileCache
}

type HeaderCache struct {
	FileCache
}

type BuildCache struct {
	FileCache
}
