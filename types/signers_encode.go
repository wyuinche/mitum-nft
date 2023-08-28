package types

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (s *Signers) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	tt uint,
	bs []byte,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*s))

	s.BaseHinter = hint.NewBaseHinter(ht)
	s.total = tt

	hinters, err := enc.DecodeSlice(bs)
	if err != nil {
		return e.Wrap(err)
	}

	signers := make([]Signer, len(hinters))
	for i, h := range hinters {
		signer, ok := h.(Signer)
		if !ok {
			return e.Wrap(util.ErrInvalid.Errorf("expected %T, not %T", Signer{}, h))
		}

		signers[i] = signer
	}
	s.signers = signers

	return nil
}
