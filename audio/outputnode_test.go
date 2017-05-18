package audio

import (
	"github.com/gordonklaus/portaudio"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func setUp() {
	portaudio.Initialize()
}

func tearDown() {
	portaudio.Terminate()
}

func TestNewOutput(t *testing.T) {
	o, err := NewOutput(44100.0)
	if o == nil {
		t.Error("Failed to instantiate output node.", err)
	}
	if o.Stream == nil {
		t.Error("Failed to instantiate PortAudio output stream.")
	}
	if o.sampleRate != 44100.0 {
		t.Errorf("Invalid sample rate. Expected 44100.0, got %f", o.sampleRate)
	}
	if len(o.inputs) == 0 {
		t.Error("Output node has no inputs.")
	}
	if len(o.overflow) == 0 {
		t.Error("Output node has no overflow buffers.")
	} else {
		for i, ov := range o.overflow {
			if len(ov) == 0 {
				t.Errorf("Overflow buffer %d has 0 size.", i)
			}
		}
	}

	o, err = NewOutput(0)
	if o != nil {
		t.Error("Instantiated output mode with invalid sample rate.")
	}
}
