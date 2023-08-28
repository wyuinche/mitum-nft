package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type NFTTransferItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address     `json:"contract"`
	Collection types.ContractID `json:"collection"`
	Receiver   base.Address     `json:"receiver"`
	IDX        uint64           `json:"nftidx"`
	Currency   types.CurrencyID `json:"currency"`
}

func (it NFTTransferItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(NFTTransferItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Collection: it.collection,
		Receiver:   it.receiver,
		IDX:        it.idx,
		Currency:   it.currency,
	})
}

type NFTTransferItemJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Collection string    `json:"collection"`
	Receiver   string    `json:"receiver"`
	IDX        uint64    `json:"nftidx"`
	Currency   string    `json:"currency"`
}

func (it *NFTTransferItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*it))

	var u NFTTransferItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return it.unmarshal(enc, u.Hint, u.Contract, u.Collection, u.Receiver, u.IDX, u.Currency)
}
