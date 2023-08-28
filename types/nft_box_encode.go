package types

import (
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (nbx *NFTBox) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ns []uint64,
) error {
	nbx.BaseHinter = hint.NewBaseHinter(ht)
	nbx.nfts = ns
	return nil
}
