package domain

type Server interface {
	ListenAndServe(addr string) error
}
