package types

type (
	VfoStatus struct {
		Vfo            DisplayVFO
		ChannelNumber  string
		ChannelName    string
		TxPower        string
		Mode           string
		SquelchSetting int
		SquelchStatus  int
	}

	Status struct {
		Vfos   [2]VfoStatus
		PttVfo int
		CtlVfo int
	}
)
