package nft

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/pkg/errors"

	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

var MintItemHint = hint.MustNewHint("mitum-nft-mint-item-v0.0.1")

type MintItem struct {
	hint.BaseHinter
	contract   base.Address
	collection currencytypes.ContractID
	hash       types.NFTHash
	uri        types.URI
	creators   types.Signers
	currency   currencytypes.CurrencyID
}

func NewMintItem(
	contract base.Address,
	collection currencytypes.ContractID,
	hash types.NFTHash,
	uri types.URI,
	creators types.Signers,
	currency currencytypes.CurrencyID,
) MintItem {
	return MintItem{
		BaseHinter: hint.NewBaseHinter(MintItemHint),
		contract:   contract,
		collection: collection,
		hash:       hash,
		uri:        uri,
		creators:   creators,
		currency:   currency,
	}
}

func (it MintItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.hash.Bytes(),
		it.uri.Bytes(),
		it.creators.Bytes(),
		it.currency.Bytes(),
	)
}

func (it MintItem) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(it))

	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.collection,
		it.hash,
		it.uri,
		it.creators,
		it.currency,
	); err != nil {
		return e.Wrap(err)
	}

	as := it.creators.Addresses()
	for _, a := range as {
		if a.Equal(it.contract) {
			return e.Wrap(errors.Errorf("contract address is same with creator, %s", it.contract))
		}
	}

	return nil
}

func (it MintItem) Contract() base.Address {
	return it.contract
}

func (it MintItem) Collection() currencytypes.ContractID {
	return it.collection
}

func (it MintItem) NFTHash() types.NFTHash {
	return it.hash
}

func (it MintItem) URI() types.URI {
	return it.uri
}

func (it MintItem) Creators() types.Signers {
	return it.creators
}

func (it MintItem) Currency() currencytypes.CurrencyID {
	return it.currency
}
