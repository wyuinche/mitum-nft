package state

import (
	"fmt"
	"strings"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	OperatorsBookSuffix         = ":operatorsbook"
	OperatorsBookStateValueHint = hint.MustNewHint("mitum-nft-operators-book-state-value-v0.0.1")
)

type OperatorsBookStateValue struct {
	hint.BaseHinter
	Book types.OperatorsBook
}

func NewOperatorsBookStateValue(book types.OperatorsBook) OperatorsBookStateValue {
	return OperatorsBookStateValue{
		BaseHinter: hint.NewBaseHinter(OperatorsBookStateValueHint),
		Book:       book,
	}
}

func (s OperatorsBookStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s OperatorsBookStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(OperatorsBookStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := s.Book.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s OperatorsBookStateValue) HashBytes() []byte {
	return s.Book.Bytes()
}

func StateOperatorsBookValue(st base.State) (types.OperatorsBook, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return types.OperatorsBook{}, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(OperatorsBookStateValue)
	if !ok {
		return types.OperatorsBook{}, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(OperatorsBookStateValue{}, v)))
	}

	return s.Book, nil
}

func StateKeyOperatorsBook(contract base.Address, collectionID currencytypes.ContractID, owner base.Address) string {
	return fmt.Sprintf("%s:%s%s", StateKeyNFTPrefix(contract, collectionID), owner, OperatorsBookSuffix)
}

func IsStateOperatorsBook(k string) bool {
	return strings.HasPrefix(k, NFTPrefix) && strings.HasSuffix(k, OperatorsBookSuffix)
}
