package types

import (
	"fmt"

	"github.com/larsks/kwctl/internal/helpers"
)

type DCS struct {
	valuePtr *int // Pointer to the dcs value being configured
}

var dcsForward map[int]string = map[int]string{
	0:   "023",
	1:   "025",
	2:   "026",
	3:   "031",
	4:   "032",
	5:   "036",
	6:   "043",
	7:   "047",
	8:   "051",
	9:   "053",
	10:  "054",
	11:  "065",
	12:  "071",
	13:  "072",
	14:  "073",
	15:  "074",
	16:  "114",
	17:  "115",
	18:  "116",
	19:  "122",
	20:  "125",
	21:  "131",
	22:  "132",
	23:  "134",
	24:  "143",
	25:  "145",
	26:  "152",
	27:  "155",
	28:  "156",
	29:  "162",
	30:  "165",
	31:  "172",
	32:  "174",
	33:  "205",
	34:  "212",
	35:  "223",
	36:  "225",
	37:  "226",
	38:  "243",
	39:  "244",
	40:  "245",
	41:  "246",
	42:  "251",
	43:  "252",
	44:  "255",
	45:  "261",
	46:  "263",
	47:  "265",
	48:  "266",
	49:  "271",
	50:  "274",
	51:  "306",
	52:  "311",
	53:  "315",
	54:  "325",
	55:  "331",
	56:  "332",
	57:  "343",
	58:  "346",
	59:  "351",
	60:  "356",
	61:  "364",
	62:  "365",
	63:  "371",
	64:  "411",
	65:  "412",
	66:  "413",
	67:  "423",
	68:  "431",
	69:  "432",
	70:  "445",
	71:  "446",
	72:  "452",
	73:  "454",
	74:  "455",
	75:  "462",
	76:  "464",
	77:  "465",
	78:  "466",
	79:  "503",
	80:  "506",
	81:  "516",
	82:  "523",
	83:  "565",
	84:  "532",
	85:  "546",
	86:  "565",
	87:  "606",
	88:  "612",
	89:  "624",
	90:  "627",
	91:  "631",
	92:  "632",
	93:  "654",
	94:  "662",
	95:  "664",
	96:  "703",
	97:  "712",
	98:  "723",
	99:  "731",
	100: "732",
	101: "734",
	102: "743",
	103: "754",
}

var dcsReverse map[string]int = helpers.ReverseMap(dcsForward)

func NewDCS(dcsPtr *int) *DCS {
	return &DCS{valuePtr: dcsPtr}
}

func (t *DCS) String() string {
	if t.valuePtr == nil {
		return ""
	}
	dcs, exists := dcsForward[*t.valuePtr]
	if !exists {
		return ""
	}
	return dcs
}

// Set parses a dcs name and stores it as an integer code
func (t *DCS) Set(value string) error {
	val, exists := dcsReverse[value]
	if !exists {
		return fmt.Errorf("invalid dcs: %s", value)
	}
	*t.valuePtr = val
	return nil
}

// Type returns the type name for help text
func (t *DCS) Type() string {
	return "dcs"
}
