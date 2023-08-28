package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type PolicyJSONMarshaler struct {
	hint.BaseHinter
	Name      CollectionName   `json:"name"`
	Royalty   PaymentParameter `json:"royalty"`
	URI       URI              `json:"uri"`
	Whitelist []base.Address   `json:"whitelist"`
}

func (p Policy) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(PolicyJSONMarshaler{
		BaseHinter: p.BaseHinter,
		Name:       p.name,
		Royalty:    p.royalty,
		URI:        p.uri,
		Whitelist:  p.whitelist,
	})
}

type PolicyJSONUnmarshaler struct {
	Hint      hint.Hint `json:"_hint"`
	Name      string    `json:"name"`
	Royalty   uint      `json:"royalty"`
	URI       string    `json:"uri"`
	Whitelist []string  `json:"whitelist"`
}

func (p *Policy) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError("failed to decode json of Policy")

	var u PolicyJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return p.unmarshal(enc, u.Hint, u.Name, u.Royalty, u.URI, u.Whitelist)
}
