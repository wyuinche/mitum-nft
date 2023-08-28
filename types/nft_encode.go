package types

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (n *NFT) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	id uint64,
	ac bool,
	ow string,
	hs string,
	uri string,
	ap string,
	bcs []byte,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*n))

	n.BaseHinter = hint.NewBaseHinter(ht)
	n.active = ac
	n.hash = NFTHash(hs)
	n.uri = URI(uri)
	n.id = id

	owner, err := base.DecodeAddress(ow, enc)
	if err != nil {
		return e.Wrap(err)
	}
	n.owner = owner

	approved, err := base.DecodeAddress(ap, enc)
	if err != nil {
		return e.Wrap(err)
	}
	n.approved = approved

	if hinter, err := enc.Decode(bcs); err != nil {
		return e.Wrap(err)
	} else if s, ok := hinter.(Signers); !ok {
		return e.Wrap(errors.Errorf("expected %T, not %T", Signers{}, hinter))
	} else {
		n.creators = s
	}

	return nil
}
