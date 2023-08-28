package types

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type SignersJSONMarshaler struct {
	hint.BaseHinter
	Total   uint     `json:"total"`
	Signers []Signer `json:"signers"`
}

func (s Signers) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SignersJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Total:      s.total,
		Signers:    s.signers,
	})
}

type SignersJSONUnmarshaler struct {
	Hint    hint.Hint       `json:"_hint"`
	Total   uint            `json:"total"`
	Signers json.RawMessage `json:"signers"`
}

func (s *Signers) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u SignersJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return s.unmarshal(enc, u.Hint, u.Total, u.Signers)
}
