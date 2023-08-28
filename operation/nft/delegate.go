package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	DelegateFactHint = hint.MustNewHint("mitum-nft-delegate-operation-fact-v0.0.1")
	DelegateHint     = hint.MustNewHint("mitum-nft-delegate-operation-v0.0.1")
)

var MaxDelegateItems = 10

type DelegateFact struct {
	base.BaseFact
	sender base.Address
	items  []DelegateItem
}

func NewDelegateFact(token []byte, sender base.Address, items []DelegateItem) DelegateFact {
	fact := DelegateFact{
		BaseFact: base.NewBaseFact(DelegateFactHint, token),
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact DelegateFact) IsValid(b []byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
	); err != nil {
		return e.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return e.Wrap(err)
	}

	if l := len(fact.items); l < 1 {
		return e.Wrap(errors.Errorf("empty items, %T", fact))
	} else if l > int(MaxItems) {
		return e.Wrap(errors.Errorf("invalid length of items, %d > max(%d)", l, MaxItems))
	}

	founds := map[string]struct{}{}
	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if item.contract.Equal(fact.sender) {
			return e.Wrap(errors.Errorf("contract address is same with sender, %s", fact.sender))
		}

		k := item.Operator().String() + ":" + item.Collection().String()

		if _, found := founds[k]; found {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("collection-operator", k)))
		}

		founds[k] = struct{}{}
	}

	return nil
}

func (fact DelegateFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact DelegateFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact DelegateFact) Bytes() []byte {
	bs := make([][]byte, len(fact.items))
	for i, item := range fact.items {
		bs[i] = item.Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(bs...),
	)
}

func (fact DelegateFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact DelegateFact) Sender() base.Address {
	return fact.sender
}

func (fact DelegateFact) Items() []DelegateItem {
	return fact.items
}

type Delegate struct {
	common.BaseOperation
}

func NewDelegate(fact DelegateFact) (Delegate, error) {
	return Delegate{BaseOperation: common.NewBaseOperation(DelegateHint, fact)}, nil
}

func (op *Delegate) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
