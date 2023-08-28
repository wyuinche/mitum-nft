package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTBoxStateValueJSONMarshaler struct {
	hint.BaseHinter
	Box types.NFTBox `json:"nftbox"`
}

func (s NFTBoxStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		NFTBoxStateValueJSONMarshaler(s),
	)
}

type NFTBoxStateValueJSONUnmarshaler struct {
	Hint hint.Hint       `json:"_hint"`
	Box  json.RawMessage `json:"nftbox"`
}

func (s *NFTBoxStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u NFTBoxStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var box types.NFTBox
	if err := box.DecodeJSON(u.Box, enc); err != nil {
		return e.Wrap(err)
	}
	s.Box = box

	return nil
}
