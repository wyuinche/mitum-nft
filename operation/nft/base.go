package nft

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

var MaxItems = 10

type Item interface {
	Currency() types.CurrencyID
}

func CalculateItemsFee(getStateFunc base.GetStateFunc, items ...any) (map[types.CurrencyID][2]common.Big, error) {
	required := map[types.CurrencyID][2]common.Big{}

	for _, item := range items {
		it, ok := item.(Item)
		if !ok {
			return nil, errors.Errorf("expected Item, not %T", item)
		}

		rq := [2]common.Big{common.ZeroBig, common.ZeroBig}

		if k, found := required[it.Currency()]; found {
			rq = k
		}

		policy, err := currencystate.ExistsCurrencyPolicy(it.Currency(), getStateFunc)
		if err != nil {
			return nil, err
		}

		switch k, err := policy.Feeer().Fee(common.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[it.Currency()] = [2]common.Big{rq[0], rq[1]}
		default:
			required[it.Currency()] = [2]common.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil

}
