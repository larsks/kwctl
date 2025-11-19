package types

import (
	flag "github.com/spf13/pflag"
)

// RadioFlagValues holds temporary storage for radio setting flag values.
// This struct is used to capture flag values before they are applied to
// the actual Channel or VFO.
type RadioFlagValues struct {
	RxFreq    int
	RxStep    int
	Mode      int
	Shift     int
	Offset    int
	ToneFreq  int
	CTCSSFreq int
	DCSCode   int
}

// AddRadioSettingFlags adds common radio setting flags to the provided FlagSet.
// The flags are bound to the RadioFlagValues struct for temporary storage.
func AddRadioSettingFlags(flags *flag.FlagSet, values *RadioFlagValues) {
	flags.VarP(NewFrequencyMHz(&values.RxFreq), "rxfreq", "r", "frequency in MHz (e.g., 144.39)")
	flags.VarP(NewStepSize(&values.RxStep), "rxstep", "", "step size in hz (e.g., 5)")
	flags.VarP(NewMode(&values.Mode), "mode", "", "Mode (FM, NFM, AM)")
	flags.VarP(NewShift(&values.Shift), "shift", "s", "Shift (simplex, up, down)")
	flags.Bool("reverse", false, "reverse tx/rx")
	flags.Bool("no-reverse", false, "disable reverse tx/rx")
	flags.VarP(NewFrequencyMHz(&values.Offset), "offset", "o", "offset in MHz (e.g., 0.6)")
	flags.StringP("tone-mode", "t", "none", "select tone mode (none, tone, tsql, dcs)")
	flags.VarP(NewTone(&values.ToneFreq), "txtone", "", "CTCSS tone when sending")
	flags.VarP(NewTone(&values.CTCSSFreq), "rxtone", "", "CTCSS tone when receiving")
	flags.VarP(NewDCS(&values.DCSCode), "dcs", "", "DCS code")
}

// ApplyRadioSettingFlags applies only the flags that were set (visited) to the target.
// This ensures that unset flags don't overwrite existing values in the target.
func ApplyRadioSettingFlags(flags *flag.FlagSet, values *RadioFlagValues, target RadioSettable) {
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "rxfreq":
			target.SetRxFreq(values.RxFreq)
		case "rxstep":
			target.SetRxStep(values.RxStep)
		case "mode":
			target.SetMode(values.Mode)
		case "shift":
			target.SetShift(values.Shift)
		case "reverse":
			target.SetReverse(1)
		case "no-reverse":
			target.SetReverse(0)
		case "offset":
			target.SetOffset(values.Offset)
		case "tone-mode":
			switch f.Value.String() {
			case "none":
				target.SetTone(0)
				target.SetCTCSS(0)
				target.SetDCS(0)
			case "tone":
				target.SetTone(1)
				target.SetCTCSS(0)
				target.SetDCS(0)
			case "tsql":
				target.SetTone(1)
				target.SetCTCSS(1)
				target.SetDCS(0)
			case "dcs":
				target.SetTone(0)
				target.SetCTCSS(0)
				target.SetDCS(1)
			}
		case "txtone":
			target.SetToneFreq(values.ToneFreq)
		case "rxtone":
			target.SetCTCSSFreq(values.CTCSSFreq)
		case "dcs":
			target.SetDCSCode(values.DCSCode)
		}
	})
}
