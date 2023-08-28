package nft

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *CollectionRegisterFact) unmarshal(
	enc encoder.Encoder,
	sd, ca, sb, nm string,
	ry uint,
	uri string,
	bws []string,
	cid string,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*fact))

	fact.currency = currencytypes.CurrencyID(cid)
	fact.collection = currencytypes.ContractID(sb)
	fact.name = types.CollectionName(nm)
	fact.royalty = types.PaymentParameter(ry)
	fact.uri = types.URI(uri)

	sender, err := base.DecodeAddress(sd, enc)
	if err != nil {
		return e.Wrap(err)
	}
	fact.sender = sender

	contract, err := base.DecodeAddress(ca, enc)
	if err != nil {
		return e.Wrap(err)
	}
	fact.contract = contract

	whitelist := make([]base.Address, len(bws))
	for i, bw := range bws {
		w, err := base.DecodeAddress(bw, enc)
		if err != nil {
			return e.Wrap(err)
		}
		whitelist[i] = w
	}
	fact.whitelist = whitelist

	return nil
}
