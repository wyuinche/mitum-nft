package types

import (
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	util "github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var SignerHint = hint.MustNewHint("mitum-nft-signer-v0.0.1")

var MaxSignerShare uint = 100

type Signer struct {
	hint.BaseHinter
	account base.Address
	share   uint
	signed  bool
}

func NewSigner(account base.Address, share uint, signed bool) Signer {
	return Signer{
		BaseHinter: hint.NewBaseHinter(SignerHint),
		account:    account,
		share:      share,
		signed:     signed,
	}
}

func (s Signer) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := util.CheckIsValiders(nil, false,
		s.BaseHinter,
		s.account,
	); err != nil {
		return e.Wrap(err)
	}

	if s.share > MaxSignerShare {
		return e.Wrap(errors.Errorf("invalid signer share, %d > max(%d)", s.share, MaxSignerShare))
	}

	return nil
}

func (s Signer) Bytes() []byte {
	return util.ConcatBytesSlice(
		s.account.Bytes(),
		util.UintToBytes(s.share),
		utils.BoolToByteSlice(s.signed),
	)
}

func (s Signer) Account() base.Address {
	return s.account
}

func (s Signer) Share() uint {
	return s.share
}

func (s Signer) Signed() bool {
	return s.signed
}

func (s Signer) Equal(c Signer) bool {
	if s.Share() != c.Share() {
		return false
	}

	if !s.Account().Equal(c.Account()) {
		return false
	}

	if s.Signed() != c.Signed() {
		return false
	}

	return true
}
