package ipdb

type IPDB interface {
	Find(ip string) string
}
