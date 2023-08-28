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
	NFTTransferFactHint = hint.MustNewHint("mitum-nft-transfer-operation-fact-v0.0.1")
	NFTTransferHint     = hint.MustNewHint("mitum-nft-transfer-operation-v0.0.1")
)

type NFTTransferFact struct {
	base.BaseFact
	sender base.Address
	items  []NFTTransferItem
}

func NewNFTTransferFact(token []byte, sender base.Address, items []NFTTransferItem) NFTTransferFact {
	fact := NFTTransferFact{
		BaseFact: base.NewBaseFact(NFTTransferFactHint, token),
		sender:   sender,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact NFTTransferFact) IsValid(b []byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(fact))

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
	); err != nil {
		return e.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return err
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

func (fact NFTTransferFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact NFTTransferFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact NFTTransferFact) Bytes() []byte {
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

func (fact NFTTransferFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact NFTTransferFact) Sender() base.Address {
	return fact.sender
}

func (fact NFTTransferFact) Items() []NFTTransferItem {
	return fact.items
}

type NFTTransfer struct {
	common.BaseOperation
}

func NewNFTTransfer(fact NFTTransferFact) (NFTTransfer, error) {
	return NFTTransfer{BaseOperation: common.NewBaseOperation(NFTTransferHint, fact)}, nil
}

func (op *NFTTransfer) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}
	return nil
}
