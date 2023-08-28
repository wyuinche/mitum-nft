package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s OperatorsBookStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":         s.Hint().String(),
			"operatorsbook": s.Book,
		},
	)
}

type OperatorsBookStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	Book bson.Raw `bson:"operatorsbook"`
}

func (s *OperatorsBookStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*s))

	var u OperatorsBookStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var book types.OperatorsBook
	if err := book.DecodeBSON(u.Book, enc); err != nil {
		return e.Wrap(err)
	}
	s.Book = book

	return nil
}
