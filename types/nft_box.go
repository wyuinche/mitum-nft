package types

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var NFTBoxHint = hint.MustNewHint("mitum-nft-nft-box-v0.0.1")

type NFTBox struct {
	hint.BaseHinter
	nfts []uint64
}

func NewNFTBox(nfts []uint64) NFTBox {
	if nfts == nil {
		nfts = []uint64{}
	}
	return NFTBox{BaseHinter: hint.NewBaseHinter(NFTBoxHint), nfts: nfts}
}

func (nbx NFTBox) Bytes() []byte {
	bs := make([][]byte, len(nbx.nfts))
	for i, n := range nbx.nfts {
		bs[i] = util.Uint64ToBytes(n)
	}
	return util.ConcatBytesSlice(bs...)
}

func (nbx NFTBox) Hint() hint.Hint {
	return NFTBoxHint
}

func (nbx NFTBox) Hash() util.Hash {
	return nbx.GenerateHash()
}

func (nbx NFTBox) GenerateHash() util.Hash {
	return valuehash.NewSHA256(nbx.Bytes())
}

func (nbx NFTBox) IsEmpty() bool {
	return len(nbx.nfts) == 0
}

func (nbx NFTBox) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(nbx))

	if err := nbx.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (nbx NFTBox) Equal(c NFTBox) bool {
	nbx.Sort(true)
	c.Sort(true)
	for i := range nbx.nfts {
		if !(nbx.nfts[i] == c.nfts[i]) {
			return false
		}
	}
	return true
}

func (nbx *NFTBox) Sort(ascending bool) {
	sort.Slice(nbx.nfts, func(i, j int) bool {
		if ascending {
			return bytes.Compare(util.Uint64ToBytes(nbx.nfts[j]), util.Uint64ToBytes(nbx.nfts[i])) > 0
		}
		return bytes.Compare(util.Uint64ToBytes(nbx.nfts[j]), util.Uint64ToBytes(nbx.nfts[i])) < 0
	})
}

func (nbx NFTBox) Exists(id uint64) bool {
	for _, n := range nbx.nfts {
		if id == n {
			return true
		}
	}
	return false
}

func (nbx *NFTBox) Append(n uint64) error {
	if nbx.Exists(n) {
		return errors.Errorf("nft already exists in NFTBox, %q", n)
	}
	if uint64(len(nbx.nfts)) >= MaxNFTIndex {
		return errors.Errorf("max nfts in collection, %q", n)
	}
	nbx.nfts = append(nbx.nfts, n)
	return nil
}

func (nbx *NFTBox) Remove(n uint64) error {
	if !nbx.Exists(n) {
		return errors.Errorf("nft not found in NFTBox, %q", n)
	}
	for i := range nbx.nfts {
		if n == nbx.nfts[i] {
			nbx.nfts[i] = nbx.nfts[len(nbx.nfts)-1]
			nbx.nfts = nbx.nfts[:len(nbx.nfts)-1]
			return nil
		}
	}
	return nil
}

func (nbx NFTBox) NFTs() []uint64 {
	return nbx.nfts
}
