package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTStateValueJSONMarshaler struct {
	hint.BaseHinter
	NFT types.NFT `json:"nft"`
}

func (s NFTStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		NFTStateValueJSONMarshaler(s),
	)
}

type NFTStateValueJSONUnmarshaler struct {
	Hint hint.Hint       `json:"_hint"`
	NFT  json.RawMessage `json:"nft"`
}

func (s *NFTStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u NFTStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var n types.NFT
	if err := n.DecodeJSON(u.NFT, enc); err != nil {
		return e.Wrap(err)
	}
	s.NFT = n

	return nil
}
