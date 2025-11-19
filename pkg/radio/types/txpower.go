package types

import (
	"fmt"
)

type (
	TxPower int
)

const (
	TX_POWER_LOW    TxPower = 2
	TX_POWER_MEDIUM TxPower = 1
	TX_POWER_HIGH   TxPower = 0
)

var txpowerNames map[string]TxPower = map[string]TxPower{
	"high":   TX_POWER_HIGH,
	"medium": TX_POWER_MEDIUM,
	"low":    TX_POWER_LOW,
}

func (t TxPower) String() string {
	switch t {
	case TX_POWER_LOW:
		return "low"
	case TX_POWER_MEDIUM:
		return "medium"
	case TX_POWER_HIGH:
		return "high"
	default:
		return "<invalid>"
	}
}

func ParseTxPower(s string) (TxPower, error) {
	if val, exists := txpowerNames[s]; exists {
		return val, nil
	}

	return 0, fmt.Errorf("invalid tx power: %s", s)
}
