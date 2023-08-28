package state

import (
	"fmt"
	"strings"

	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

var NFTPrefix = "nft:"

func StateKeyNFTPrefix(contract base.Address, collectionID types.ContractID) string {
	return fmt.Sprintf("%s%s:%s", NFTPrefix, contract.String(), collectionID.String())
}

type StateKeyGenerator struct {
	contract     base.Address
	collectionID types.ContractID
}

func NewStateKeyGenerator(contract base.Address, collectionID types.ContractID) StateKeyGenerator {
	return StateKeyGenerator{
		contract,
		collectionID,
	}
}

func (g StateKeyGenerator) Design() string {
	return StateKeyDesign(g.contract, g.collectionID)
}

func (g StateKeyGenerator) NFT(idx uint64) string {
	return StateKeyNFT(g.contract, g.collectionID, idx)
}

func (g StateKeyGenerator) NFTBox() string {
	return StateKeyNFTBox(g.contract, g.collectionID)
}

func (g StateKeyGenerator) LastNFTIndex() string {
	return StateKeyLastNFTIndex(g.contract, g.collectionID)
}

func (g StateKeyGenerator) OperatorsBook(owner base.Address) string {
	return StateKeyOperatorsBook(g.contract, g.collectionID, owner)
}

func ParseStateKey(key string) ([]string, error) {
	parsed := strings.Split(key, ":")

	if parsed[0] != NFTPrefix[:len(NFTPrefix)-1] {
		return nil, errors.Errorf("state key doesn't include NFTPrefix, %s", parsed)
	}
	if len(parsed) < 3 {
		return nil, errors.Errorf("failed to parse state key, %s", parsed)
	}

	return parsed, nil
}
