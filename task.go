package pool

type task struct {
	id      int
	address string
	result  []byte
	err     error
}
