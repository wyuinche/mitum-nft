package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s NFTBoxStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  s.Hint().String(),
			"nftbox": s.Box,
		},
	)
}

type NFTBoxStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	Box  bson.Raw `bson:"nftbox"`
}

func (s *NFTBoxStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*s))

	var u NFTBoxStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var box types.NFTBox
	if err := box.DecodeBSON(u.Box, enc); err != nil {
		return e.Wrap(err)
	}
	s.Box = box

	return nil
}
