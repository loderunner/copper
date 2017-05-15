package audio

type Buffer []float64
type Channel chan Buffer

type interface Renderer {
	func Render() (uint, error)
}

type interface Input {
    func Connect(c Channel, input uint) (bool, error)
    func Disconnect(input uint)
}

type interface Output {
	func Channel(i uint) (Channel)
}

type interface InputOutput {
	Input,
	Output
}
