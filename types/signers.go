package types

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	MaxTotalShare uint = 100
	MaxSigners         = 10
)

var SignersHint = hint.MustNewHint("mitum-nft-signers-v0.0.1")

type Signers struct {
	hint.BaseHinter
	total   uint
	signers []Signer
}

func NewSigners(total uint, signers []Signer) Signers {
	return Signers{
		BaseHinter: hint.NewBaseHinter(SignersHint),
		total:      total,
		signers:    signers,
	}
}

func (s Signers) IsValid([]byte) error {
	e := mitumutil.ErrInvalid.Errorf(utils.ErrStringInvalid(s))

	if err := s.BaseHinter.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	if s.total > MaxTotalShare {
		return e.Wrap(errors.Errorf("invalid total share, %d > max(%d)", s.total, MaxTotalShare))
	}

	if l := len(s.signers); l > MaxSigners {
		return e.Wrap(errors.Errorf("invalid length of signers, %d > max(%d)", l, MaxSigners))
	}

	var total uint = 0

	founds := map[string]struct{}{}
	for _, sn := range s.signers {
		if err := sn.IsValid(nil); err != nil {
			return e.Wrap(err)
		}

		ac := sn.Account().String()
		if _, found := founds[ac]; found {
			return e.Wrap(errors.Errorf(utils.ErrStringDuplicate("signer", ac)))
		}
		founds[ac] = struct{}{}

		total += sn.Share()
	}

	if total != s.total {
		return e.Wrap(errors.Errorf("total share must be equal to the sum of all shares, %d != %d", s.total, total))
	}

	return nil
}

func (s Signers) Bytes() []byte {
	bs := make([][]byte, len(s.signers))

	for i, signer := range s.signers {
		bs[i] = signer.Bytes()
	}

	return mitumutil.ConcatBytesSlice(
		mitumutil.UintToBytes(s.total),
		mitumutil.ConcatBytesSlice(bs...),
	)
}

func (s Signers) Total() uint {
	return s.total
}

func (s Signers) Signers() []Signer {
	return s.signers
}

func (s Signers) Addresses() []base.Address {
	as := make([]base.Address, len(s.signers))
	for i, signer := range s.signers {
		as[i] = signer.Account()
	}
	return as
}

func (s Signers) Index(signer Signer) int {
	return s.IndexByAddress(signer.Account())
}

func (s Signers) IndexByAddress(address base.Address) int {
	for i := range s.signers {
		if address.Equal(s.signers[i].Account()) {
			return i
		}
	}
	return -1
}

func (s Signers) Exists(signer Signer) bool {
	if idx := s.Index(signer); idx >= 0 {
		return true
	}
	return false
}

func (s Signers) Equal(c Signers) bool {
	if s.Total() != c.Total() {
		return false
	}

	if len(s.Signers()) != len(c.Signers()) {
		return false
	}

	ss := s.Signers()
	sort.Slice(ss, func(i, j int) bool {
		return bytes.Compare(ss[j].Bytes(), ss[i].Bytes()) < 0
	})

	cs := c.Signers()
	sort.Slice(cs, func(i, j int) bool {
		return bytes.Compare(cs[j].Bytes(), cs[i].Bytes()) < 0
	})

	for i := range ss {
		if !ss[i].Equal(cs[i]) {
			return false
		}
	}

	return true
}

func (s Signers) IsSigned(signer Signer) bool {
	return s.IsSignedByAddress(signer.Account())
}

func (s Signers) IsSignedByAddress(address base.Address) bool {
	idx := s.IndexByAddress(address)
	if idx < 0 {
		return false
	}
	return s.signers[idx].Signed()
}

func (s *Signers) SetSigner(signer Signer) error {
	idx := s.Index(signer)
	if idx < 0 {
		return errors.Errorf("signer not in signers, %q", signer.Account())
	}
	s.signers[idx] = signer
	return nil
}
