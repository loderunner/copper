package audio

type Buffer []float64
type Channel chan Buffer

type RenderNode interface {
	SetSampleRate(sampleRate float64)
	Render() (uint, error)
}

type InputNode interface {
	NumInputs() uint
	Connect(c Channel, i uint) (bool, Error)
	Disconnect(i uint)
}

type OutputNode interface {
	NumChannels() uint
	Channel(i uint) Channel
}

type Error string

func (err *Error) Error() string {
	return string(err)
}
