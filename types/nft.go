package types

import (
	"strings"

	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var MaxNFTHashLength = 1024
var MaxNFTIndex uint64 = 10000

type NFTHash string

func (h NFTHash) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(h))

	if l := len(h); l > MaxNFTHashLength {
		return e.Wrap(errors.Errorf("invalid length of nft hash, %d > max(%d)", l, MaxNFTHashLength))
	}

	if h != "" && strings.TrimSpace(string(h)) == "" {
		return e.Wrap(errors.Errorf("empty nft hash"))
	}

	return nil
}

func (hs NFTHash) Bytes() []byte {
	return []byte(hs)
}

func (hs NFTHash) String() string {
	return string(hs)
}

var NFTHint = hint.MustNewHint("mitum-nft-nft-v0.0.1")

var MaxCreators = 10

type NFT struct {
	hint.BaseHinter
	id       uint64
	active   bool
	owner    base.Address
	hash     NFTHash
	uri      URI
	approved base.Address
	creators Signers
}

func NewNFT(
	id uint64,
	active bool,
	owner base.Address,
	hash NFTHash,
	uri URI,
	approved base.Address,
	creators Signers,
) NFT {
	return NFT{
		BaseHinter: hint.NewBaseHinter(NFTHint),
		id:         id,
		active:     active,
		owner:      owner,
		hash:       hash,
		uri:        uri,
		approved:   approved,
		creators:   creators,
	}
}

func (n NFT) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(n))

	if err := util.CheckIsValiders(nil, false,
		n.owner,
		n.hash,
		n.uri,
		n.approved,
		n.creators,
	); err != nil {
		return e.Wrap(err)
	}

	if n.uri == "" {
		return e.Wrap(errors.Errorf("empty uri"))
	}

	return nil
}

func (n NFT) Bytes() []byte {
	return util.ConcatBytesSlice(
		util.Uint64ToBytes(n.id),
		utils.BoolToByteSlice(n.active),
		n.owner.Bytes(),
		n.hash.Bytes(),
		[]byte(n.uri.String()),
		n.approved.Bytes(),
		n.creators.Bytes(),
	)
}

func (n NFT) ID() uint64 {
	return n.id
}

func (n NFT) Active() bool {
	return n.active
}

func (n NFT) Owner() base.Address {
	return n.owner
}

func (n NFT) NFTHash() NFTHash {
	return n.hash
}

func (n NFT) URI() URI {
	return n.uri
}

func (n NFT) Approved() base.Address {
	return n.approved
}

func (n NFT) Creators() Signers {
	return n.creators
}

func (n NFT) Addresses() []base.Address {
	as := n.Creators().Addresses()

	for i, a := range as {
		if n.approved == a {
			break
		}
		if i == len(as)-1 {
			as = append(as, n.approved)
		}
	}

	for i, a := range as {
		if n.owner == a {
			break
		}
		if i == (len(as) - 1) {
			as = append(as, n.owner)
		}
	}

	return as
}

func (n NFT) Equal(c NFT) bool {
	if !(n.ID() == c.ID()) {
		return false
	}

	if n.Active() != c.Active() {
		return false
	}

	if !n.Owner().Equal(c.Owner()) {
		return false
	}

	if n.NFTHash() != c.NFTHash() {
		return false
	}

	if n.URI() != c.URI() {
		return false
	}

	if !n.Approved().Equal(c.Approved()) {
		return false
	}

	if !n.Creators().Equal(c.Creators()) {
		return false
	}

	return n.ID() == c.ID()
}

func (n NFT) ExistsApproved() bool {
	return !n.approved.Equal(n.owner)
}
