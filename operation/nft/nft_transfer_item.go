package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var NFTTransferItemHint = hint.MustNewHint("mitum-nft-transfer-item-v0.0.1")

type NFTTransferItem struct {
	hint.BaseHinter
	contract   base.Address
	collection types.ContractID
	receiver   base.Address
	idx        uint64
	currency   types.CurrencyID
}

func NewNFTTransferItem(contract base.Address, collection types.ContractID, receiver base.Address, idx uint64, currency types.CurrencyID) NFTTransferItem {
	return NFTTransferItem{
		BaseHinter: hint.NewBaseHinter(NFTTransferItemHint),
		contract:   contract,
		collection: collection,
		receiver:   receiver,
		idx:        idx,
		currency:   currency,
	}
}

func (it NFTTransferItem) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(it))

	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.collection,
		it.receiver,
		it.currency,
	); err != nil {
		return e.Wrap(err)
	}

	if it.contract.Equal(it.receiver) {
		return e.Wrap(errors.Errorf("contract address is same with receiver, %s", it.receiver))
	}

	return nil
}

func (it NFTTransferItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.receiver.Bytes(),
		util.Uint64ToBytes(it.idx),
		it.currency.Bytes(),
	)
}

func (it NFTTransferItem) Contract() base.Address {
	return it.contract
}

func (it NFTTransferItem) Collection() types.ContractID {
	return it.collection
}

func (it NFTTransferItem) Receiver() base.Address {
	return it.receiver
}

func (it NFTTransferItem) NFTIdx() uint64 {
	return it.idx
}

func (it NFTTransferItem) Currency() types.CurrencyID {
	return it.currency
}
