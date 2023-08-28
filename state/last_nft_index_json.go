package state

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type LastNFTIndexStateValueJSONMarshaler struct {
	hint.BaseHinter
	Index uint64 `json:"index"`
}

func (s LastNFTIndexStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		LastNFTIndexStateValueJSONMarshaler{
			BaseHinter: s.BaseHinter,
			Index:      s.id,
		},
	)
}

type LastNFTIndexStateValueJSONUnmarshaler struct {
	Hint  hint.Hint `json:"_hint"`
	Index uint64    `json:"index"`
}

func (s *LastNFTIndexStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u LastNFTIndexStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}
	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	s.id = u.Index

	return nil
}
