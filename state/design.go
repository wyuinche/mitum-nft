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
	DesignSuffix         = ":design"
	DesignStateValueHint = hint.MustNewHint("mitum-nft-design-state-value-v0.0.1")
)

type DesignStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		Design:     design,
	}
}

func (s DesignStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := s.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s DesignStateValue) HashBytes() []byte {
	return s.Design.Bytes()
}

func StateDesignValue(st base.State) (types.Design, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return types.Design{}, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(DesignStateValue)
	if !ok {
		return types.Design{}, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(DesignStateValue{}, v)))
	}

	return s.Design, nil
}

func StateKeyDesign(contract base.Address, collectionID currencytypes.ContractID) string {
	return fmt.Sprintf("%s:%s", StateKeyNFTPrefix(contract, collectionID), DesignSuffix)
}

func IsStateDesignKey(k string) bool {
	return strings.HasPrefix(k, NFTPrefix) && strings.HasSuffix(k, DesignSuffix)
}
