package ipdb

// ip db interface
type IPDB interface {
	Find(ip string) string
}
