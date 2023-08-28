package types

import (
	"net/url"
	"strings"

	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"

	"github.com/pkg/errors"
)

var MaxPaymentParameter uint = 99

type PaymentParameter uint

func (p PaymentParameter) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(p))

	if uint(p) > MaxPaymentParameter {
		return e.Wrap(errors.Errorf("invalid payment parameter length, %d > max(%d)", p, MaxPaymentParameter))
	}

	return nil
}

func (p PaymentParameter) Bytes() []byte {
	return util.UintToBytes(uint(p))
}

func (p PaymentParameter) Uint() uint {
	return uint(p)
}

var MaxURILength = 1000

type URI string

func (u URI) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(u))

	if _, err := url.Parse(string(u)); err != nil {
		return e.Wrap(err)
	}

	if l := len(u); l > MaxURILength {
		return e.Wrap(errors.Errorf("invalid uri length, %d > max(%d)", l, MaxURILength))
	}

	if u != "" && strings.TrimSpace(string(u)) == "" {
		return e.Wrap(errors.Errorf("empty uri"))
	}

	return nil
}

func (u URI) Bytes() []byte {
	return []byte(u)
}

func (u URI) String() string {
	return string(u)
}

var DesignHint = hint.MustNewHint("mitum-nft-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	parent     base.Address
	creator    base.Address
	collection types.ContractID
	active     bool
	policy     Policy
}

func NewDesign(parent base.Address, creator base.Address, collection types.ContractID, active bool, policy Policy) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		parent:     parent,
		creator:    creator,
		collection: collection,
		active:     active,
		policy:     policy,
	}
}

func (d Design) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(d))

	if err := util.CheckIsValiders(nil, false,
		d.BaseHinter,
		d.parent,
		d.creator,
		d.collection,
		d.policy,
	); err != nil {
		return e.Wrap(err)
	}

	if d.parent.Equal(d.creator) {
		return e.Wrap(errors.Errorf("parent address is same with creator, %s", d.creator))
	}

	return nil
}

func (d Design) Bytes() []byte {
	return util.ConcatBytesSlice(
		d.parent.Bytes(),
		d.creator.Bytes(),
		d.collection.Bytes(),
		utils.BoolToByteSlice(d.active),
		d.policy.Bytes(),
	)
}

func (d Design) Hash() util.Hash {
	return d.GenerateHash()
}

func (d Design) GenerateHash() util.Hash {
	return valuehash.NewSHA256(d.Bytes())
}

func (d Design) Parent() base.Address {
	return d.parent
}

func (d Design) Creator() base.Address {
	return d.creator
}

func (d Design) Collection() types.ContractID {
	return d.collection
}

func (d Design) Active() bool {
	return d.active
}

func (d Design) Policy() Policy {
	return d.policy
}

func (d Design) Equal(c Design) bool {
	if !d.parent.Equal(c.parent) {
		return false
	}

	if !d.creator.Equal(c.creator) {
		return false
	}

	if d.collection != c.collection {
		return false
	}

	if d.active != c.active {
		return false
	}

	if !d.policy.Equal(c.policy) {
		return false
	}

	if d.Hash() != c.Hash() {
		return false
	}

	return true
}
