package bucket

type Worker struct {
	name string
	running chan bool
}
