package pool

type task struct {
	id     int
	url    string
	result []byte
	err    error
}
