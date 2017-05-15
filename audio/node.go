package audio

type Buffer []float64
type Channel chan Buffer

type Renderer interface {
	Render() (uint, error)
}

type Input interface {
	Connect(c Channel, input uint) (bool, error)
	Disconnect(input uint)
}

type Output interface {
	Channel(i uint) Channel
}

type InputOutput interface {
	Input
	Output
}
