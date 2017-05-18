package audio

type Buffer []float32
type Channel chan Buffer

type RenderNode interface {
	SetSampleRate(sampleRate float64)
	Render() (int, error)
}

type InputNode interface {
	NumInputs() int
	Connect(c Channel, i int) (bool, error)
	Disconnect(i int)
}

type OutputNode interface {
	NumChannels() int
	Channel(i int) Channel
}

type Error string

func (err Error) Error() string {
	return string(err)
}
