package worker

type Job interface {
	Execute() error
}
