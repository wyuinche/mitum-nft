package state

import (
	"fmt"
	"strconv"
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
	NFTSuffix         = ":nft"
	NFTStateValueHint = hint.MustNewHint("nft-state-value-v0.0.1")
)

type NFTStateValue struct {
	hint.BaseHinter
	NFT types.NFT
}

func NewNFTStateValue(n types.NFT) NFTStateValue {
	return NFTStateValue{
		BaseHinter: hint.NewBaseHinter(NFTStateValueHint),
		NFT:        n,
	}
}

func (s NFTStateValue) Hint() hint.Hint {
	return s.BaseHinter.Hint()
}

func (s NFTStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(NFTStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := s.NFT.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (s NFTStateValue) HashBytes() []byte {
	return s.NFT.Bytes()
}

func StateNFTValue(st base.State) (types.NFT, error) {
	e := util.ErrNotFound.Errorf(ErrStringStateNotFound(st.Key()))

	v := st.Value()
	if v == nil {
		return types.NFT{}, e.Wrap(errors.Errorf("nil value"))
	}

	s, ok := v.(NFTStateValue)
	if !ok {
		return types.NFT{}, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(NFTStateValue{}, v)))
	}

	return s.NFT, nil
}

func StateKeyNFT(contract base.Address, collectionID currencytypes.ContractID, id uint64) string {
	return fmt.Sprintf("%s:%s%s", StateKeyNFTPrefix(contract, collectionID), strconv.FormatUint(id, 10), NFTSuffix)
}

func IsStateNFTKey(k string) bool {
	return strings.HasPrefix(k, NFTPrefix) && strings.HasSuffix(k, NFTSuffix)
}
