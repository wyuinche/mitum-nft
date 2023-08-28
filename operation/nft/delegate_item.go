package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	DelegateAllow  = DelegateMode("allow")
	DelegateCancel = DelegateMode("cancel")
)

type DelegateMode string

func (mode DelegateMode) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(mode))

	if !(mode == DelegateAllow || mode == DelegateCancel) {
		return e.Wrap(errors.Errorf("wrong delegate mode, %s", mode))
	}

	return nil
}

func (mode DelegateMode) Bytes() []byte {
	return []byte(mode)
}

func (mode DelegateMode) String() string {
	return string(mode)
}

func (mode DelegateMode) Equal(cmode DelegateMode) bool {
	return string(mode) == string(cmode)
}

var DelegateItemHint = hint.MustNewHint("mitum-nft-delegate-item-v0.0.1")

type DelegateItem struct {
	hint.BaseHinter
	contract   base.Address
	collection types.ContractID
	operator   base.Address
	mode       DelegateMode
	currency   types.CurrencyID
}

func NewDelegateItem(contract base.Address, collection types.ContractID, operator base.Address, mode DelegateMode, currency types.CurrencyID) DelegateItem {
	return DelegateItem{
		BaseHinter: hint.NewBaseHinter(DelegateItemHint),
		contract:   contract,
		collection: collection,
		operator:   operator,
		mode:       mode,
		currency:   currency,
	}
}

func (it DelegateItem) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(it))

	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
		it.collection,
		it.operator,
		it.mode,
		it.currency,
	); err != nil {
		return e.Wrap(err)
	}

	if it.contract.Equal(it.operator) {
		return e.Wrap(errors.Errorf("contract address is same with operator, %s", it.contract))
	}

	return nil
}

func (it DelegateItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		it.collection.Bytes(),
		it.operator.Bytes(),
		it.mode.Bytes(),
		it.currency.Bytes(),
	)
}

func (it DelegateItem) Contract() base.Address {
	return it.contract
}

func (it DelegateItem) Collection() types.ContractID {
	return it.collection
}

func (it DelegateItem) Operator() base.Address {
	return it.operator
}

func (it DelegateItem) Mode() DelegateMode {
	return it.mode
}

func (it DelegateItem) Currency() types.CurrencyID {
	return it.currency
}
