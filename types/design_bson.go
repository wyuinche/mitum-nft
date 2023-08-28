package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (d Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      d.Hint().String(),
			"parent":     d.parent,
			"creator":    d.creator,
			"collection": d.collection,
			"active":     d.active,
			"policy":     d.policy,
		})
}

type DesignBSONUnmarshaler struct {
	Hint       string   `bson:"_hint"`
	Parent     string   `bson:"parent"`
	Creator    string   `bson:"creator"`
	Collection string   `bson:"collection"`
	Active     bool     `bson:"active"`
	Policy     bson.Raw `bson:"policy"`
}

func (d *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*d))

	var u DesignBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(enc, ht, u.Parent, u.Creator, u.Collection, u.Active, u.Policy)
}
