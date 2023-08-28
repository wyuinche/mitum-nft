package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *ApproveItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col,
	ap string,
	idx uint64,
	cid string,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*it))

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.currency = types.CurrencyID(cid)
	it.collection = types.ContractID(col)
	it.idx = idx

	contract, err := base.DecodeAddress(ca, enc)
	if err != nil {
		return e.Wrap(err)
	}
	it.contract = contract

	approved, err := base.DecodeAddress(ap, enc)
	if err != nil {
		return e.Wrap(err)
	}
	it.approved = approved

	return nil
}
