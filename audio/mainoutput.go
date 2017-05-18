package audio

import (
	"github.com/gordonklaus/portaudio"
)

type MainOutput struct {
	*portaudio.Stream
	sampleRate float64
	inputs     []Channel
	overflow   []Buffer
}

func NewMainOutput(sampleRate float64) (*MainOutput, error) {
	o := MainOutput{sampleRate: sampleRate}

	device, err := portaudio.DefaultOutputDevice()
	if device == nil {
		return nil, err
	}

	numChannels := device.MaxOutputChannels
	o.inputs = make([]Channel, numChannels)
	for i := 0; i < numChannels; i++ {
		o.inputs[i] = make(Channel)
	}
	ok, err := o.initOutputStream()
	if !ok {
		return nil, err
	}

	minSize := int(device.DefaultHighInputLatency.Seconds() * sampleRate)
	bufferSize := 1
	for bufferSize < minSize {
		bufferSize <<= 1
	}
	o.overflow = make([]Buffer, numChannels)
	for i := 0; i < numChannels; i++ {
		o.overflow[i] = make(Buffer, bufferSize)
	}

	return &o, nil
}

func (o *MainOutput) SetSampleRate(sampleRate float64) {
	o.sampleRate = sampleRate
	o.initOutputStream()
}

func (o *MainOutput) initOutputStream() (bool, error) {
	var err error
	o.Stream, err = portaudio.OpenDefaultStream(0, len(o.inputs), o.sampleRate, 0, o.streamCallback)
	if o.Stream == nil {
		return false, err
	}
	return true, nil
}

func (o *MainOutput) Render()        {}
func (o *MainOutput) NumInputs() int { return len(o.inputs) }

func (o *MainOutput) InputChannel(i int) Channel {
	return o.inputs[i]
}

func (o *MainOutput) streamCallback(outputBuffers []Buffer) {
	for i, outputChannelBuffer := range outputBuffers {
		inputChannel := o.inputs[i]
		overflow := o.overflow[i]

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
				o.overflow[i] = overflow[:remaining]
				copy(o.overflow[i], inputChannelBuffer[remaining:])
			} else {
				o.overflow[i] = overflow[:0]
			}
		}
	}
}
