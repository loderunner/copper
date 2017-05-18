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
}

func TestNewOutputZeroSampleRate(t *testing.T) {
	o, _ := NewOutput(0)
	if o != nil {
		t.Error("Instantiated output node with invalid sample rate: 0")
	}
}

func TestSetSampleRate(t *testing.T) {
	o, _ := NewOutput(44100.0)

	o.SetSampleRate(48000.0)
	validateOutput(t, o, 48000.0)
}

func TestSetSampleRateZero(t *testing.T) {
	o, _ := NewOutput(44100.0)

	o.SetSampleRate(0)
	if o.Stream != nil {
		t.Error("Set invalid sample rate: 0")
	}
}

func TestRender(t *testing.T) {}

func TestNumImputs(t *testing.T) {
	o, _ := NewOutput(44100.0)
	got, want := o.NumInputs(), len(o.inputs)
	if got != want {
		t.Error("NumInputs() returned %d, expected %d", got, want)
	}
}

func TestConnect(t *testing.T) {
	o, _ := NewOutput(44100.0)

	c := make(Channel)
	ok, err := o.Connect(c, 0)
	if !ok {
		t.Error("Failed to channel to input 0.", err)
	}
}

func TestConnectInvalidChannel(t *testing.T) {
	o, _ := NewOutput(44100.0)

	c := make(Channel)
	ok, _ := o.Connect(c, 256)
	if ok {
		t.Error("Unexpected connection to input 256.")
	}
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
