package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type OperatorsBookJSONMarshaler struct {
	hint.BaseHinter
	Collection types.ContractID `json:"collection"`
	Operators  []base.Address   `json:"operators"`
}

func (o OperatorsBook) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperatorsBookJSONMarshaler{
		BaseHinter: o.BaseHinter,
		Collection: o.collection,
		Operators:  o.operators,
	})
}

type OperatorsBookJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Collection string    `json:"collection"`
	Operators  []string  `json:"operators"`
}

func (o *OperatorsBook) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*o))

	var u OperatorsBookJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return o.unmarshal(enc, u.Hint, u.Collection, u.Operators)
}
