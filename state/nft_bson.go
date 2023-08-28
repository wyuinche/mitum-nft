package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s NFTStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": s.Hint().String(),
			"nft":   s.NFT,
		},
	)
}

type NFTStateValueBSONUnmarshaler struct {
	Hint string   `bson:"_hint"`
	NFT  bson.Raw `bson:"nft"`
}

func (s *NFTStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*s))

	var u NFTStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	var n types.NFT
	if err := n.DecodeBSON(u.NFT, enc); err != nil {
		return e.Wrap(err)
	}
	s.NFT = n

	return nil
}
