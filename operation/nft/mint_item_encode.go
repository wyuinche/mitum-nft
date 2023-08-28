package nft

import (
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/pkg/errors"

	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *MintItem) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	ca, col, hs, uri string,
	bcr []byte,
	cid string,
) error {
	e := util.StringError("failed to unmarshal MintItem")

	it.BaseHinter = hint.NewBaseHinter(ht)
	it.collection = currencytypes.ContractID(col)
	it.hash = types.NFTHash(hs)
	it.uri = types.URI(uri)
	it.currency = currencytypes.CurrencyID(cid)

	contract, err := base.DecodeAddress(ca, enc)
	if err != nil {
		return e.Wrap(err)
	}
	it.contract = contract

	if hinter, err := enc.Decode(bcr); err != nil {
		return e.Wrap(err)
	} else if creators, ok := hinter.(types.Signers); !ok {
		return e.Wrap(errors.Errorf(utils.ErrStringTypeCast(types.Signers{}, hinter)))
	} else {
		it.creators = creators
	}

	return nil
}
