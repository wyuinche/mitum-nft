package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

func (d *Design) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	pr, cr, sb string,
	ac bool,
	bpo []byte,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*d))

	d.BaseHinter = hint.NewBaseHinter(ht)
	d.collection = types.ContractID(sb)
	d.active = ac

	parent, err := base.DecodeAddress(pr, enc)
	if err != nil {
		return e.Wrap(err)
	}
	d.parent = parent

	creator, err := base.DecodeAddress(cr, enc)
	if err != nil {
		return e.Wrap(err)
	}
	d.creator = creator

	if hinter, err := enc.Decode(bpo); err != nil {
		return e.Wrap(err)
	} else if po, ok := hinter.(Policy); !ok {
		return e.Wrap(errors.Errorf("expected %T, not %T", Policy{}, hinter))
	} else {
		d.policy = po
	}

	return nil
}
