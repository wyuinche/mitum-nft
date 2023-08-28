package types

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var MaxOperators = 10

var OperatorsBookHint = hint.MustNewHint("mitum-nft-operator-book-v0.0.1")

type OperatorsBook struct {
	hint.BaseHinter
	collection types.ContractID
	operators  []base.Address
}

func NewOperatorsBook(collection types.ContractID, operators []base.Address) OperatorsBook {
	if operators == nil {
		operators = []base.Address{}
	}
	return OperatorsBook{BaseHinter: hint.NewBaseHinter(OperatorsBookHint), collection: collection, operators: operators}
}

func (o OperatorsBook) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf(utils.ErrStringInvalid(o))

	for _, op := range o.operators {
		if err := op.IsValid(nil); err != nil {
			return e.Wrap(err)
		}
	}

	return nil
}

func (o OperatorsBook) Bytes() []byte {
	bs := make([][]byte, len(o.operators))

	for i, operator := range o.operators {
		bs[i] = operator.Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func (o OperatorsBook) Hash() util.Hash {
	return o.GenerateHash()
}

func (o OperatorsBook) GenerateHash() util.Hash {
	return valuehash.NewSHA256(o.Bytes())
}

func (o OperatorsBook) IsEmpty() bool {
	return len(o.operators) < 1
}

func (o OperatorsBook) Collection() types.ContractID {
	return o.collection
}

func (o OperatorsBook) Equal(c OperatorsBook) bool {
	o.Sort(true)
	c.Sort(true)

	for i := range o.operators {
		if !o.operators[i].Equal(c.operators[i]) {
			return false
		}
	}

	return true
}

func (o *OperatorsBook) Sort(ascending bool) {
	sort.Slice(o.operators, func(i, j int) bool {
		if ascending {
			return bytes.Compare(o.operators[j].Bytes(), o.operators[i].Bytes()) > 0
		}
		return bytes.Compare(o.operators[j].Bytes(), o.operators[i].Bytes()) < 0
	})
}

func (o OperatorsBook) Exists(ag base.Address) bool {
	for _, operator := range o.operators {
		if ag.Equal(operator) {
			return true
		}
	}
	return false
}

func (o OperatorsBook) Get(ag base.Address) (base.Address, error) {
	for _, operator := range o.operators {
		if ag.Equal(operator) {
			return operator, nil
		}
	}
	return types.Address{}, errors.Errorf("account not in operators book, %q", ag)
}

func (o *OperatorsBook) Append(ag base.Address) error {
	if err := ag.IsValid(nil); err != nil {
		return err
	}

	if o.Exists(ag) {
		return errors.Errorf("account already in operators book, %q", ag)
	}

	if len(o.operators) >= MaxOperators {
		return errors.Errorf("max operators, %v", ag)
	}

	o.operators = append(o.operators, ag)

	return nil
}

func (o *OperatorsBook) Remove(ag base.Address) error {
	if !o.Exists(ag) {
		return errors.Errorf("account not in operators book, %q", ag)
	}

	for i := range o.operators {
		if ag.String() == o.operators[i].String() {
			o.operators[i] = o.operators[len(o.operators)-1]
			o.operators[len(o.operators)-1] = nil
			o.operators = o.operators[:len(o.operators)-1]
			break
		}
	}

	return nil
}

func (o OperatorsBook) Operators() []base.Address {
	return o.operators
}
