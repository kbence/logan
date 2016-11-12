package pipeline

type OutputPipeline interface {
	Start() chan bool
}
