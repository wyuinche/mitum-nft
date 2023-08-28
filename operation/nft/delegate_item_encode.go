package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *DelegateItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col, op, md, cid string,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*it))

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.collection = types.ContractID(col)
	it.mode = DelegateMode(md)
	it.currency = types.CurrencyID(cid)

	contract, err := base.DecodeAddress(ca, enc)
	if err != nil {
		return e.Wrap(err)
	}
	it.contract = contract

	operator, err := base.DecodeAddress(op, enc)
	if err != nil {
		return e.Wrap(err)
	}
	it.operator = operator

	return nil
}
