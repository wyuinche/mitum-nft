package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (s LastNFTIndexStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": s.Hint().String(),
			"index": s.id,
		},
	)
}

type LastNFTIndexStateValueBSONUnmarshaler struct {
	Hint  string `bson:"_hint"`
	Index uint64 `bson:"index"`
}

func (s *LastNFTIndexStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*s))

	var u LastNFTIndexStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	s.BaseHinter = hint.NewBaseHinter(ht)

	s.id = u.Index

	return nil
}
