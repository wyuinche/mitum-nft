package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignJSONMarshaler struct {
	hint.BaseHinter
	Parent     base.Address     `json:"parent"`
	Creator    base.Address     `json:"creator"`
	Collection types.ContractID `json:"collection"`
	Active     bool             `json:"active"`
	Policy     Policy           `json:"policy"`
}

func (d Design) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DesignJSONMarshaler{
		BaseHinter: d.BaseHinter,
		Parent:     d.parent,
		Creator:    d.creator,
		Collection: d.collection,
		Active:     d.active,
		Policy:     d.policy,
	})
}

type DesignJSONUnmarshaler struct {
	Hint       hint.Hint       `json:"_hint"`
	Parent     string          `json:"parent"`
	Creator    string          `json:"creator"`
	Collection string          `json:"collection"`
	Active     bool            `json:"active"`
	Policy     json.RawMessage `json:"policy"`
}

func (d *Design) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*d))

	var u DesignJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(enc, u.Hint, u.Parent, u.Creator, u.Collection, u.Active, u.Policy)
}
