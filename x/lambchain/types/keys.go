package types

const (
	// ModuleName defines the module name
	ModuleName = "lambchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_lambchain"
)

var (
	ParamsKey = []byte("p_lambchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
