package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	jsonenc "github.com/ProtoconNet/mitum2/util/encoder/json"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type ApproveItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address     `json:"contract"`
	Collection types.ContractID `json:"collection"`
	Approved   base.Address     `json:"approved"`
	NFTidx     uint64           `json:"nftidx"`
	Currency   types.CurrencyID `json:"currency"`
}

func (it ApproveItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ApproveItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		Collection: it.collection,
		Approved:   it.approved,
		NFTidx:     it.idx,
		Currency:   it.currency,
	})
}

type ApproveItemJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	Collection string    `json:"collection"`
	Approved   string    `json:"approved"`
	NFTidx     uint64    `json:"nftidx"`
	Currency   string    `json:"currency"`
}

func (it *ApproveItem) DecodeJSON(b []byte, enc *jsonenc.Encoder) error {
	e := util.StringError(utils.ErrStringDecodeJSON(*it))

	var u ApproveItemJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return it.unmarshal(enc, u.Hint, u.Contract, u.Collection, u.Approved, u.NFTidx, u.Currency)
}
