package nft

import (
	"strconv"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	NFTSignFactHint = hint.MustNewHint("mitum-nft-sign-operation-fact-v0.0.1")
	NFTSignHint     = hint.MustNewHint("mitum-nft-sign-operation-v0.0.1")
)

type NFTSignFact struct {
	base.BaseFact
	sender base.Address
	items  []NFTSignItem
}

func NewNFTSignFact(token []byte, sender base.Address, items []NFTSignItem) NFTSignFact {
	fact := NFTSignFact{
		BaseFact: base.NewBaseFact(NFTSignFactHint, token),
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact NFTSignFact) IsValid(b []byte) error {
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

		nid := strconv.FormatUint(item.NFTIdx(), 10)
		if _, found := founds[nid]; found {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("nft idx", nid)))
		}

		founds[nid] = struct{}{}
	}

	return nil
}

func (fact NFTSignFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact NFTSignFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact NFTSignFact) Bytes() []byte {
	bs := make([][]byte, len(fact.items))
	for i := range fact.items {
		bs[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(bs...),
	)
}

func (fact NFTSignFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact NFTSignFact) Sender() base.Address {
	return fact.sender
}

func (fact NFTSignFact) Items() []NFTSignItem {
	return fact.items
}

type NFTSign struct {
	common.BaseOperation
}

func NewNFTSign(fact NFTSignFact) (NFTSign, error) {
	return NFTSign{BaseOperation: common.NewBaseOperation(NFTSignHint, fact)}, nil
}

func (op *NFTSign) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
