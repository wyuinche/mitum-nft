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
	ApproveFactHint = hint.MustNewHint("mitum-nft-approve-operation-fact-v0.0.1")
	ApproveHint     = hint.MustNewHint("mitum-nft-approve-operation-v0.0.1")
)

type ApproveFact struct {
	base.BaseFact
	sender base.Address
	items  []ApproveItem
}

func NewApproveFact(token []byte, sender base.Address, items []ApproveItem) ApproveFact {
	fact := ApproveFact{
		BaseFact: base.NewBaseFact(ApproveFactHint, token),
		sender:   sender,
		items:    items,
	}

	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact ApproveFact) IsValid(b []byte) error {
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

		n := strconv.FormatUint(item.NFTIdx(), 10)

		if _, found := founds[n]; found {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("nft", n)))
		}

		founds[n] = struct{}{}
	}

	return nil
}

func (fact ApproveFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact ApproveFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact ApproveFact) Bytes() []byte {
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

func (fact ApproveFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact ApproveFact) Sender() base.Address {
	return fact.sender
}

func (fact ApproveFact) Items() []ApproveItem {
	return fact.items
}

type Approve struct {
	common.BaseOperation
}

func NewApprove(fact ApproveFact) (Approve, error) {
	return Approve{BaseOperation: common.NewBaseOperation(ApproveHint, fact)}, nil
}

func (op *Approve) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
