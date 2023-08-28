package types

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (p *Policy) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	nm string,
	ry uint,
	uri string,
	bws []string,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*p))

	p.BaseHinter = hint.NewBaseHinter(ht)
	p.name = CollectionName(nm)
	p.royalty = PaymentParameter(ry)
	p.uri = URI(uri)

	whitelist := make([]base.Address, len(bws))
	for i, w := range bws {
		a, err := base.DecodeAddress(w, enc)
		if err != nil {
			return e.Wrap(err)
		}
		whitelist[i] = a
	}
	p.whitelist = whitelist

	return nil
}
