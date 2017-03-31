package pool

type task struct {
	id      int
	address string
	proxy   string
	result  []byte
	err     error
}
