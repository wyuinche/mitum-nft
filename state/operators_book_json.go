package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"

	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type OperatorsBookStateValueJSONMarshaler struct {
	hint.BaseHinter
	Book types.OperatorsBook `json:"operatorsbook"`
}

func (s OperatorsBookStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		OperatorsBookStateValueJSONMarshaler(s),
	)
}

type OperatorsBookStateValueJSONUnmarshaler struct {
	Hint hint.Hint       `json:"_hint"`
	Book json.RawMessage `json:"operatorsbook"`
}

func (s *OperatorsBookStateValue) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u OperatorsBookStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.BaseHinter = hint.NewBaseHinter(u.Hint)

	var book types.OperatorsBook
	if err := book.DecodeJSON(u.Book, enc); err != nil {
		return e.Wrap(err)
	}
	s.Book = book

	return nil
}
