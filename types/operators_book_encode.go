package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (o *OperatorsBook) unmarshal(
	enc encoder.Encoder,
	ht hint.Hint,
	col string,
	bags []string,
) error {
	e := util.StringError(utils.ErrStringUnmarshal(*o))

	o.BaseHinter = hint.NewBaseHinter(ht)
	o.collection = types.ContractID(col)

	operators := make([]base.Address, len(bags))
	for i, b := range bags {
		a, err := base.DecodeAddress(b, enc)
		if err != nil {
			return e.Wrap(err)
		}
		operators[i] = a
	}
	o.operators = operators

	return nil
}
