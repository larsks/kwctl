package types

import (
	"testing"
)

// TestChannelImplementsRadioSettable verifies that Channel implements RadioSettable
func TestChannelImplementsRadioSettable(t *testing.T) {
	var _ RadioSettable = &Channel{}

	c := Channel{}
	c.SetRxFreq(145000000)
	c.SetMode(1)
	c.SetToneFreq(88)
	c.SetReverse(1)

	if c.GetRxFreq() != 145000000 {
		t.Errorf("GetRxFreq() = %d, expected 145000000", c.GetRxFreq())
	}
	if c.GetMode() != 1 {
		t.Errorf("GetMode() = %d, expected 1", c.GetMode())
	}
	if c.GetToneFreq() != 88 {
		t.Errorf("GetToneFreq() = %d, expected 88", c.GetToneFreq())
	}
	if c.GetReverse() != 1 {
		t.Errorf("GetReverse() = %d, expected 1", c.GetReverse())
	}
}

// TestVFOImplementsRadioSettable verifies that VFO implements RadioSettable
func TestVFOImplementsRadioSettable(t *testing.T) {
	var _ RadioSettable = &VFO{}

	v := VFO{}
	v.SetRxFreq(146820000)
	v.SetShift(2)
	v.SetCTCSS(1)

	if v.GetRxFreq() != 146820000 {
		t.Errorf("GetRxFreq() = %d, expected 146820000", v.GetRxFreq())
	}
	if v.GetShift() != 2 {
		t.Errorf("GetShift() = %d, expected 2", v.GetShift())
	}
	if v.GetCTCSS() != 1 {
		t.Errorf("GetCTCSS() = %d, expected 1", v.GetCTCSS())
	}
}
