package state

import (
	"fmt"
	"strings"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	LastNFTIDXSuffix           = ":lastnftidx"
	LastNFTIndexStateValueHint = hint.MustNewHint("mitum-nft-last-nft-index-state-value-v0.0.1")
)

type LastNFTIndexStateValue struct {
	hint.BaseHinter
	id uint64
}

func NewLastNFTIndexStateValue(id uint64) LastNFTIndexStateValue {
	return LastNFTIndexStateValue{
		BaseHinter: hint.NewBaseHinter(LastNFTIndexStateValueHint),
		id:         id,
	}
}

func (s LastNFTIndexStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s LastNFTIndexStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(LastNFTIndexStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s LastNFTIndexStateValue) HashBytes() []byte {
	return util.Uint64ToBytes(s.id)
}

func StateLastNFTIndexValue(st base.State) (uint64, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return 0, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(LastNFTIndexStateValue)
	if !ok {
		return 0, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(LastNFTIndexStateValue{}, v)))
	}

	return s.id, nil
}

func StateKeyLastNFTIndex(contract base.Address, collectionID currencytypes.ContractID) string {
	return fmt.Sprintf("%s:%s", StateKeyNFTPrefix(contract, collectionID), LastNFTIDXSuffix)
}

func IsStateLastNFTIndexKey(k string) bool {
	return strings.HasPrefix(k, NFTPrefix) && strings.HasSuffix(k, LastNFTIDXSuffix)
}
