package audio

import (
	"github.com/gordonklaus/portaudio"
)

type Output struct {
	*portaudio.Stream
	sampleRate float64
	inputs     []Channel
	overflow   []Buffer
}

func NewOutput(sampleRate float64) (*Output, error) {
	output := Output{sampleRate: sampleRate}

	device, err := portaudio.DefaultOutputDevice()
	if device == nil {
		return nil, err
	}

	numChannels := device.MaxOutputChannels
	output.inputs = make([]Channel, numChannels)
	ok, err := output.initOutputStream()
	if !ok {
		return nil, err
	}

	minSize := int(device.DefaultHighInputLatency.Seconds() * sampleRate)
	bufferSize := 1
	for bufferSize < minSize {
		bufferSize <<= 1
	}
	output.overflow = make([]Buffer, numChannels)
	for i := 0; i < numChannels; i++ {
		output.overflow[i] = make(Buffer, bufferSize)
	}

	return &output, nil
}

func (output *Output) SetSampleRate(sampleRate float64) {
	output.sampleRate = sampleRate
	output.initOutputStream()
}

func (output *Output) initOutputStream() (bool, error) {
	var err error
	output.Stream, err = portaudio.OpenDefaultStream(0, len(output.inputs), output.sampleRate, 0, output.streamCallback)
	if output.Stream == nil {
		return false, err
	}
	return true, nil
}

func (output *Output) Render()        {}
func (output *Output) NumInputs() int { return len(output.inputs) }

func (output *Output) Connect(c Channel, i int) (bool, error) {
	if i >= len(output.inputs) {
		return false, Error("Input index out of bounds.")
	}

	output.inputs[i] = c
	return true, nil
}

func (output *Output) Disconnect(i int) {
	if i < len(output.inputs) {
		output.inputs[i] = nil
	}
}

func (output *Output) streamCallback(outputBuffers []Buffer) {
	for i, outputChannelBuffer := range outputBuffers {
		inputChannel := output.inputs[i]
		overflow := output.overflow[i]

		j := 0
		// If remaining overflow, copy to output
		if len(overflow) > 0 {
			j += copy(outputChannelBuffer, overflow)
		}

		// Receive buffers from input and copy to output
		for j < len(outputChannelBuffer) {
			inputChannelBuffer := <-inputChannel
			j += copy(outputChannelBuffer[j:], inputChannelBuffer) // Store overflowing buffer for next callback
			if j > len(outputChannelBuffer) {
				remaining := j - len(outputChannelBuffer)
				output.overflow[i] = overflow[:remaining]
				copy(output.overflow[i], inputChannelBuffer[remaining:])
			} else {
				output.overflow[i] = overflow[:0]
			}
		}
	}
}
