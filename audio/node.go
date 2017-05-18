package audio

type Buffer []float32
type Channel chan Buffer

type RenderNode interface {
	SetSampleRate(sampleRate float64)
	Render() (int, error)
}

type InputNode interface {
	NumInputs() int
	InputChannel(i int) Channel
}

type OutputNode interface {
	NumOutputs() int
	Connect(c Channel, i int) (bool, error)
	Disconnect(i int)
}

type Error string

func (err Error) Error() string {
	return string(err)
}
