package tranpack

import (
	"public/tlv"
	"riskmanagement/tranpack/position"
)

func ClosePosition(p position.Position) ([]byte, error) {
	tlvpk := new(tlv.TlvPacker)
	tlvpk.
}
