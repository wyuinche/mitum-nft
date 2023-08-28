package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var NFTSignItemHint = hint.MustNewHint("mitum-nft-sign-item-v0.0.1")

type NFTSignItem struct {
	hint.BaseHinter
	contract   base.Address
	collection types.ContractID
	idx        uint64
	currency   types.CurrencyID
}

func NewNFTSignItem(contract base.Address, collection types.ContractID, idx uint64, currency types.CurrencyID) NFTSignItem {
	return NFTSignItem{
		BaseHinter: hint.NewBaseHinter(NFTSignItemHint),
		contract:   contract,
		collection: collection,
		idx:        idx,
		currency:   currency,
	}
}

func (it NFTSignItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		util.Uint64ToBytes(it.idx),
		it.currency.Bytes(),
	)
}

func (it NFTSignItem) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(it))

	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.collection,
		it.currency,
	); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (it NFTSignItem) NFTIdx() uint64 {
	return it.idx
}

func (it NFTSignItem) Contract() base.Address {
	return it.contract
}

func (it NFTSignItem) Collection() types.ContractID {
	return it.collection
}

func (it NFTSignItem) Currency() types.CurrencyID {
	return it.currency
}
