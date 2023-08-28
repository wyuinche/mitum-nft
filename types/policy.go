package types

import (
	"bytes"
	"regexp"
	"sort"

	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var MaxWhitelist = 10

var (
	MinLengthCollectionName = 3
	MaxLengthCollectionName = 30
	ReValidCollectionName   = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\s]+$`)
)

type CollectionName string

func (n CollectionName) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(n))

	l := len(n)
	if l < MinLengthCollectionName {
		return e.Wrap(errors.Errorf("invalid length of collection name, %d < min(%d)", l, MinLengthCollectionName))
	}
	if l > MaxLengthCollectionName {
		return e.Wrap(errors.Errorf("invalid length of collection name, %d > max(%d)", l, MaxLengthCollectionName))
	}

	if !ReValidCollectionName.Match([]byte(n)) {
		return e.Wrap(errors.New(utils.ErrStringFormat("collection name", n.String())))
	}

	return nil
}

func (n CollectionName) Bytes() []byte {
	return []byte(n)
}

func (n CollectionName) String() string {
	return string(n)
}

var PolicyHint = hint.MustNewHint("mitum-nft-collection-policy-v0.0.1")

type Policy struct {
	hint.BaseHinter
	name      CollectionName
	royalty   PaymentParameter
	uri       URI
	whitelist []base.Address
}

func NewPolicy(name CollectionName, royalty PaymentParameter, uri URI, whitelist []base.Address) Policy {
	return Policy{
		BaseHinter: hint.NewBaseHinter(PolicyHint),
		name:       name,
		royalty:    royalty,
		uri:        uri,
		whitelist:  whitelist,
	}
}

func (p Policy) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(p))

	if err := util.CheckIsValiders(nil, false,
		p.name,
		p.royalty,
		p.uri,
	); err != nil {
		return e.Wrap(err)
	}

	if l := len(p.whitelist); l > MaxWhitelist {
		return e.Wrap(errors.Errorf("invalid length of whitelist, %d > max(%d)", l, MaxWhitelist))
	}

	founds := map[string]struct{}{}
	for _, a := range p.whitelist {
		if err := a.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
		if _, found := founds[a.String()]; found {
			return e.Wrap(errors.New(utils.ErrStringDuplicate("whitelist account", a.String())))
		}
		founds[a.String()] = struct{}{}
	}

	return nil
}

func (p Policy) Bytes() []byte {
	bs := make([][]byte, len(p.whitelist))
	for i, white := range p.whitelist {
		bs[i] = white.Bytes()
	}

	return util.ConcatBytesSlice(
		p.name.Bytes(),
		p.royalty.Bytes(),
		p.uri.Bytes(),
		util.ConcatBytesSlice(bs...),
	)
}

func (p Policy) Name() CollectionName {
	return p.name
}

func (p Policy) Royalty() PaymentParameter {
	return p.royalty
}

func (p Policy) URI() URI {
	return p.uri
}

func (p Policy) Whitelist() []base.Address {
	return p.whitelist
}

func (p Policy) Equal(c Policy) bool {
	if p.name != c.name {
		return false
	}

	if p.royalty != c.royalty {
		return false
	}

	if p.uri != c.uri {
		return false
	}

	if len(p.whitelist) != len(c.whitelist) {
		return false
	}

	wl := p.Whitelist()
	cwl := c.Whitelist()
	sort.Slice(wl, func(i, j int) bool {
		return bytes.Compare(wl[j].Bytes(), wl[i].Bytes()) < 0
	})
	sort.Slice(cwl, func(i, j int) bool {
		return bytes.Compare(wl[j].Bytes(), wl[i].Bytes()) < 0
	})

	for i := range wl {
		if !wl[i].Equal(cwl[i]) {
			return false
		}
	}

	return true
}
