package nft

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

func (fact *ApproveFact) unmarshal(enc encoder.Encoder, sd string, bit []byte) error {
	e := util.StringError(utils.ErrStringUnmarshal(*fact))

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e.Wrap(err)
	}
	fact.sender = sender

	hits, err := enc.DecodeSlice(bit)
	if err != nil {
		return e.Wrap(err)
	}

	items := make([]ApproveItem, len(hits))
	for i, hinter := range hits {
		item, ok := hinter.(ApproveItem)
		if !ok {
			return e.Wrap(errors.Errorf(utils.ErrStringTypeCast(ApproveItem{}, hinter)))
		}

		items[i] = item
	}
	fact.items = items

	return nil
}
