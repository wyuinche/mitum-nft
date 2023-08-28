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

var MaxMintItems = 10

var (
	MintFactHint = hint.MustNewHint("mitum-nft-mint-operation-fact-v0.0.1")
	MintHint     = hint.MustNewHint("mitum-nft-mint-operation-v0.0.1")
)

type MintFact struct {
	base.BaseFact
	sender base.Address
	items  []MintItem
}

func NewMintFact(token []byte, sender base.Address, items []MintItem) MintFact {
	fact := MintFact{
		BaseFact: base.NewBaseFact(MintFactHint, token),
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MintFact) IsValid(b []byte) error {
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

	for _, item := range fact.items {
		if err := item.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		if item.contract.Equal(fact.sender) {
			return e.Wrap(errors.Errorf("contract address is same with sender, %s", fact.sender))
		}
	}

	return nil
}

func (fact MintFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact MintFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MintFact) Bytes() []byte {
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

func (fact MintFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact MintFact) Sender() base.Address {
	return fact.sender
}

func (fact MintFact) Items() []MintItem {
	return fact.items
}

type Mint struct {
	common.BaseOperation
}

func NewMint(fact MintFact) (Mint, error) {
	return Mint{BaseOperation: common.NewBaseOperation(MintHint, fact)}, nil
}

func (op *Mint) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
