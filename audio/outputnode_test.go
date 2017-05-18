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
		t.Fatal("Failed to instantiate output node.", err)
	}
	validateOutput(t, o, 44100.0)

	o, err = NewOutput(0)
	if o != nil {
		t.Fatal("Instantiated output mode with invalid sample rate.")
	}
}

func TestSetSampleRate(t *testing.T) {
	o, err := NewOutput(44100.0)
	if o == nil {
		t.Fatal("Failed to instantiate output node.", err)
	}
	validateOutput(t, o, 44100.0)

	o.SetSampleRate(48000.0)
	validateOutput(t, o, 48000.0)
}

func validateOutput(t *testing.T, o *Output, sampleRate float64) {
	if o.Stream == nil {
		t.Error("Failed to instantiate PortAudio output stream.")
	}
	if o.sampleRate != sampleRate {
		t.Errorf("Invalid sample rate. Expected %f, got %f", sampleRate, o.sampleRate)
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
}
