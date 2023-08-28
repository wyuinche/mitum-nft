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
	CollectionPolicyUpdaterFactHint = hint.MustNewHint("mitum-nft-collection-policy-updater-operation-fact-v0.0.1")
	CollectionPolicyUpdaterHint     = hint.MustNewHint("mitum-nft-collection-policy-updater-operation-v0.0.1")
)

type CollectionPolicyUpdaterFact struct {
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

func NewCollectionPolicyUpdaterFact(
	token []byte,
	sender, contract base.Address,
	collection currencytypes.ContractID,
	name types.CollectionName,
	royalty types.PaymentParameter,
	uri types.URI,
	whitelist []base.Address,
	currency currencytypes.CurrencyID,
) CollectionPolicyUpdaterFact {
	fact := CollectionPolicyUpdaterFact{
		BaseFact:   base.NewBaseFact(CollectionPolicyUpdaterFactHint, token),
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

func (fact CollectionPolicyUpdaterFact) IsValid(b []byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return e.Wrap(err)
	}

	if err := util.CheckIsValiders(
		nil, false,
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

	if l := len(fact.whitelist); l > types.MaxWhitelist {
		return e.Wrap(errors.Errorf("invalid length of whitelist, %d > max(%d)", l, types.MaxWhitelist))
	}

	if fact.sender.Equal(fact.contract) {
		return e.Wrap(errors.Errorf("contract address is same with sender, %s", fact.sender))
	}

	founds := map[string]struct{}{}
	for _, a := range fact.whitelist {
		if err := a.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if fact.contract.Equal(a) {
			return e.Wrap(errors.Errorf("contract address is same with whitelist account, %s", fact.contract))
		}

		if _, found := founds[a.String()]; found {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("whitelist account", a.String())))
		}

		founds[a.String()] = struct{}{}
	}

	return nil
}

func (fact CollectionPolicyUpdaterFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CollectionPolicyUpdaterFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CollectionPolicyUpdaterFact) Bytes() []byte {
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

func (fact CollectionPolicyUpdaterFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact CollectionPolicyUpdaterFact) Sender() base.Address {
	return fact.sender
}

func (fact CollectionPolicyUpdaterFact) Contract() base.Address {
	return fact.contract
}

func (fact CollectionPolicyUpdaterFact) Collection() currencytypes.ContractID {
	return fact.collection
}

func (fact CollectionPolicyUpdaterFact) Name() types.CollectionName {
	return fact.name
}

func (fact CollectionPolicyUpdaterFact) Royalty() types.PaymentParameter {
	return fact.royalty
}

func (fact CollectionPolicyUpdaterFact) URI() types.URI {
	return fact.uri
}

func (fact CollectionPolicyUpdaterFact) Whitelist() []base.Address {
	return fact.whitelist
}

func (fact CollectionPolicyUpdaterFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

type CollectionPolicyUpdater struct {
	common.BaseOperation
}

func NewCollectionPolicyUpdater(fact CollectionPolicyUpdaterFact) (CollectionPolicyUpdater, error) {
	return CollectionPolicyUpdater{BaseOperation: common.NewBaseOperation(CollectionPolicyUpdaterHint, fact)}, nil
}

func (op *CollectionPolicyUpdater) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
