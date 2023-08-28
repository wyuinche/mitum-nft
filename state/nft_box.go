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
	NFTBoxSuffix         = ":nftbox"
	NFTBoxStateValueHint = hint.MustNewHint("mitum-nft-nft-box-state-value-v0.0.1")
)

type NFTBoxStateValue struct {
	hint.BaseHinter
	Box types.NFTBox
}

func NewNFTBoxStateValue(box types.NFTBox) NFTBoxStateValue {
	return NFTBoxStateValue{
		BaseHinter: hint.NewBaseHinter(NFTBoxStateValueHint),
		Box:        box,
	}
}

func (s NFTBoxStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s NFTBoxStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(NFTBoxStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := s.Box.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s NFTBoxStateValue) HashBytes() []byte {
	return s.Box.Bytes()
}

func StateNFTBoxValue(st base.State) (types.NFTBox, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return types.NFTBox{}, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(NFTBoxStateValue)
	if !ok {
		return types.NFTBox{}, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(NFTBoxStateValue{}, v)))
	}

	return s.Box, nil
}

func StateKeyNFTBox(contract base.Address, collectionID currencytypes.ContractID) string {
	return fmt.Sprintf("%s:%s", StateKeyNFTPrefix(contract, collectionID), NFTBoxSuffix)
}

func IsStateNFTBoxKey(k string) bool {
	return strings.HasPrefix(k, NFTPrefix) && strings.HasSuffix(k, NFTBoxSuffix)
}
