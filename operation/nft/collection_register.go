package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	CollectionRegisterFactHint = hint.MustNewHint("mitum-nft-collection-register-operation-fact-v0.0.1")
	CollectionRegisterHint     = hint.MustNewHint("mitum-nft-collection-register-operation-v0.0.1")
)

type CollectionRegisterFact struct {
	base.BaseFact
	sender     base.Address
	contract   base.Address
	collection currencytypes.ContractID
	name       types.CollectionName
	royalty    types.PaymentParameter
	uri        types.URI
	whitelist  []base.Address
	currency   currencytypes.CurrencyID
}

func NewCollectionRegisterFact(
	token []byte,
	sender base.Address,
	contract base.Address,
	collection currencytypes.ContractID,
	name types.CollectionName,
	royalty types.PaymentParameter,
	uri types.URI,
	whitelist []base.Address,
	currency currencytypes.CurrencyID,
) CollectionRegisterFact {
	fact := CollectionRegisterFact{
		BaseFact:   base.NewBaseFact(CollectionRegisterFactHint, token),
		sender:     sender,
		contract:   contract,
		collection: collection,
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whitelist:  whitelist,
		currency:   currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact CollectionRegisterFact) IsValid(b []byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.collection,
		fact.name,
		fact.royalty,
		fact.uri,
		fact.currency,
	); err != nil {
		return e.Wrap(err)
	}

	if fact.sender.Equal(fact.contract) {
		return e.Wrap(errors.Errorf("contract address is same with sender, %s", fact.sender))
	}

	return nil
}

func (fact CollectionRegisterFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CollectionRegisterFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CollectionRegisterFact) Bytes() []byte {
	bs := make([][]byte, len(fact.whitelist))
	for i, a := range fact.whitelist {
		bs[i] = a.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.collection.Bytes(),
		fact.name.Bytes(),
		fact.royalty.Bytes(),
		fact.uri.Bytes(),
		fact.currency.Bytes(),
		util.ConcatBytesSlice(bs...),
	)
}

func (fact CollectionRegisterFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CollectionRegisterFact) Sender() base.Address {
	return fact.sender
}

func (fact CollectionRegisterFact) Contract() base.Address {
	return fact.contract
}

func (fact CollectionRegisterFact) Collection() currencytypes.ContractID {
	return fact.collection
}

func (fact CollectionRegisterFact) Name() types.CollectionName {
	return fact.name
}

func (fact CollectionRegisterFact) Royalty() types.PaymentParameter {
	return fact.royalty
}

func (fact CollectionRegisterFact) URI() types.URI {
	return fact.uri
}

func (fact CollectionRegisterFact) Whitelist() []base.Address {
	return fact.whitelist
}

func (fact CollectionRegisterFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

type CollectionRegister struct {
	common.BaseOperation
}

func NewCollectionRegister(fact CollectionRegisterFact) (CollectionRegister, error) {
	return CollectionRegister{BaseOperation: common.NewBaseOperation(CollectionRegisterHint, fact)}, nil
}

func (op *CollectionRegister) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
