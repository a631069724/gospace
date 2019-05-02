package tranpack

import (
	"public/tlv"
	"riskmanagement/tranpack/position"
)

func ClosePosition(p position.Position, channel string, flag string) ([]byte, error) {
	tlvpk := new(tlv.TlvPacker)
	tlvpk.Set("000", []byte("21"))
	tlvpk.Set("001", []byte(p.GetUserId()))
	tlvpk.Set("003", []byte(p.GetContract()))
	tlvpk.Set("006", []byte(p.GetPositionId()))
	tlvpk.Set("020", []byte(p.GetCloseBuyOrSell()))
	tlvpk.Set("021", []byte(p.GetPositionNum()))
	tlvpk.Set("060", []byte(channel))
	tlvpk.Set("061", []byte(flag))
	return tlvpk.Pack(), nil
}
