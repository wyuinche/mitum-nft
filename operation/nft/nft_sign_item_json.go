package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTSignItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address     `json:"contract"`
	Collection types.ContractID `json:"collection"`
	IDX        uint64           `json:"nftidx"`
	Currency   types.CurrencyID `json:"currency"`
}

func (it NFTSignItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTSignItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Collection: it.collection,
		IDX:        it.idx,
		Currency:   it.currency,
	})
}

type NFTSignItemJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Collection string    `json:"collection"`
	NFT        uint64    `json:"nftidx"`
	Currency   string    `json:"currency"`
}

func (it *NFTSignItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*it))

	var u NFTSignItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return it.unmarshal(enc, u.Hint, u.Contract, u.Collection, u.NFT, u.Currency)
}
