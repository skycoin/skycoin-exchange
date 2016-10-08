package mobile

type noder interface {
	GetBalance(addr string) (uint64, error)
	ValidateAddr(addr string) error
}
