package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var ApproveItemHint = hint.MustNewHint("mitum-nft-approve-item-v0.0.1")

type ApproveItem struct {
	hint.BaseHinter
	contract   base.Address
	collection types.ContractID
	approved   base.Address
	idx        uint64
	currency   types.CurrencyID
}

func NewApproveItem(contract base.Address, collection types.ContractID, approved base.Address, idx uint64, currency types.CurrencyID) ApproveItem {
	return ApproveItem{
		BaseHinter: hint.NewBaseHinter(ApproveItemHint),
		contract:   contract,
		collection: collection,
		approved:   approved,
		idx:        idx,
		currency:   currency,
	}
}

func (it ApproveItem) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(it))

	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.collection,
		it.approved,
		it.currency,
	); err != nil {
		return e.Wrap(err)
	}

	if it.contract.Equal(it.approved) {
		return e.Wrap(errors.Errorf("contract address is same with approved, %s", it.contract))
	}

	return nil
}

func (it ApproveItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.approved.Bytes(),
		util.Uint64ToBytes(it.idx),
		it.currency.Bytes(),
	)
}

func (it ApproveItem) Contract() base.Address {
	return it.contract
}

func (it ApproveItem) Collection() types.ContractID {
	return it.collection
}

func (it ApproveItem) Approved() base.Address {
	return it.approved
}

func (it ApproveItem) NFTIdx() uint64 {
	return it.idx
}

func (it ApproveItem) Currency() types.CurrencyID {
	return it.currency
}
