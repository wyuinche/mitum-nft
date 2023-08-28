package nft

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *NFTSignFact) unmarshal(
	enc encoder.Encoder,
	sd string,
	bits []byte,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*fact))

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e.Wrap(err)
	}
	fact.sender = sender

	hits, err := enc.DecodeSlice(bits)
	if err != nil {
		return e.Wrap(err)
	}

	items := make([]NFTSignItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(NFTSignItem)
		if !ok {
			return e.Wrap(errors.Errorf(utils.ErrStringTypeCast(NFTSignItem{}, hinter)))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
