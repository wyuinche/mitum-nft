package types

import (
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (p Policy) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":     p.Hint().String(),
		"name":      p.name,
		"royalty":   p.royalty,
		"uri":       p.uri,
		"whitelist": p.whitelist,
	})
}

type PolicyBSONUnmarshaler struct {
	Hint    string   `bson:"_hint"`
	Name    string   `bson:"name"`
	Royalty uint     `bson:"royalty"`
	URI     string   `bson:"uri"`
	Whites  []string `bson:"whitelist"`
}

func (p *Policy) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeBSON(*p))

	var u PolicyBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return p.unmarshal(enc, ht, u.Name, u.Royalty, u.URI, u.Whites)
}
