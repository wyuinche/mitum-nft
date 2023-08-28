package types

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (s *Signer) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ac string,
	sh uint,
	sg bool,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*s))

	s.BaseHinter = hint.NewBaseHinter(ht)
	s.share = sh
	s.signed = sg

	a, err := base.DecodeAddress(ac, enc)
	if err != nil {
		return e.Wrap(err)
	}
	s.account = a

	return nil
}
