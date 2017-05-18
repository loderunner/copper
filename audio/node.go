package audio

type Buffer []float32
type Channel chan Buffer

type Renderer interface {
	SetSampleRate(sampleRate float64)
	Render() (int, error)
}

type Input interface {
	NumInputs() int
	InputChannel(i int) Channel
}

type Output interface {
	NumOutputs() int
	Connect(c Channel, i int) (bool, error)
	Disconnect(i int)
}

type Error string

func (err Error) Error() string {
	return string(err)
}
