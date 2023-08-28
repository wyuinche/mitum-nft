package nft

import (
	"context"
	"fmt"
	"sync"

	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"

	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var delegateItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateItemProcessor)
	},
}

var delegateProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DelegateProcessor)
	},
}

func (Delegate) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type DelegateItemProcessor struct {
	h      util.Hash
	sender base.Address
	box    *types.OperatorsBook
	item   DelegateItem
}

func (ipp *DelegateItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError(ErrStringPreProcess(*ipp))

	it := ipp.item

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(it.Operator()), getStateFunc); err != nil {
		return e.Wrap(ErrStateNotFound("operator", it.Operator().String(), err))
	}

	if ipp.sender.Equal(it.Operator()) {
		return e.Wrap(errors.Errorf("sender cannot be operator itself, %s", it.Operator()))
	}

	return nil
}

func (ipp *DelegateItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError(ErrStringProcess(*ipp))

	it := ipp.item

	if ipp.box == nil {
		return nil, e.Wrap(ErrStateNotFound("nft box", utils.StringerChain(it.Contract(), it.Collection()), nil))
	}

	switch ipp.item.Mode() {
	case DelegateAllow:
		if err := ipp.box.Append(it.Operator()); err != nil {
			return nil, e.Wrap(errors.Errorf("failed to append operator, %s: %v", it.Operator(), err))
		}
	case DelegateCancel:
		if err := ipp.box.Remove(ipp.item.Operator()); err != nil {
			return nil, e.Wrap(errors.Errorf("failed to remove operator, %s: %v", it.Operator(), err))
		}
	default:
		return nil, e.Wrap(errors.Errorf("wrong mode for delegate item, %s: \"allow\" | \"cancel\"", ipp.item.Mode()))
	}

	ipp.box.Sort(true)

	return nil, nil
}

func (ipp *DelegateItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = DelegateItem{}
	ipp.box = nil

	delegateItemProcessorPool.Put(ipp)

	return nil
}

type DelegateProcessor struct {
	*base.BaseOperationProcessor
}

func NewDelegateProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := DelegateProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := delegateProcessorPool.Get()
		opp, ok := nopp.(*DelegateProcessor)
		if !ok {
			return nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(t, nopp)))
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *DelegateProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(DelegateFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(DelegateFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, BaseErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot have operators", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	for _, item := range fact.Items() {
		g := state.NewStateKeyGenerator(item.Contract(), item.Collection())
		k := utils.StringerChain(item.Contract(), item.Collection())

		st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
		if err != nil {
			return nil, BaseErrStateNotFound("design", k, err), nil
		}

		design, err := state.StateDesignValue(st)
		if err != nil {
			return nil, BaseErrStateNotFound("design value", k, err), nil
		}

		if !design.Active() {
			return nil, ErrBaseOperationProcess("deactivated collection", k, nil), nil
		}

		st, err = currencystate.ExistsState(stateextension.StateKeyContractAccount(design.Parent()), "key of contract", getStateFunc)
		if err != nil {
			return nil, BaseErrStateNotFound("parent", design.Parent().String(), err), nil
		}

		ca, err := stateextension.StateContractAccountValue(st)
		if err != nil {
			return nil, BaseErrStateNotFound("parent value", design.Parent().String(), err), nil
		}

		if !ca.IsActive() {
			return nil, ErrBaseOperationProcess("deactivated contract account", design.Parent().String(), nil), nil
		}
	}

	for _, item := range fact.Items() {
		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&DelegateItemProcessor{}, ip)))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, ErrBaseOperationProcess("", "", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *DelegateProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(DelegateFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(DelegateFact{}, op.Fact())))
	}

	boxes := map[string]*types.OperatorsBook{}
	for _, item := range fact.Items() {
		g := state.NewStateKeyGenerator(item.Contract(), item.Collection())
		k := g.OperatorsBook(fact.Sender())

		var operators types.OperatorsBook
		switch st, found, err := getStateFunc(k); {
		case err != nil:
			return nil, BaseErrStateNotFound("operators", k, err), nil
		case !found:
			operators = types.NewOperatorsBook(item.Collection(), nil)
		default:
			o, err := state.StateOperatorsBookValue(st)
			if err != nil {
				return nil, BaseErrStateNotFound("operators value", k, err), nil
			} else {
				operators = o
			}
		}
		boxes[k] = &operators
	}

	var sts []base.StateMergeValue // nolint:prealloc

	ipcs := make([]*DelegateItemProcessor, len(fact.items))
	for i, item := range fact.Items() {
		g := state.NewStateKeyGenerator(item.Contract(), item.Collection())

		ip := delegateItemProcessorPool.Get()
		ipc, ok := ip.(*DelegateItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&DelegateItemProcessor{}, ip)))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.box = boxes[g.OperatorsBook(fact.Sender())]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, ErrBaseOperationProcess("", "", err), nil
		}
		sts = append(sts, s...)

		ipcs[i] = ipc
	}

	for ak, box := range boxes {
		bv := currencystate.NewStateMergeValue(ak, state.NewOperatorsBookStateValue(*box))
		sts = append(sts, bv)
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	required, err := CalculateItemsFee(getStateFunc, fact.items)
	if err != nil {
		return nil, ErrBaseOperationProcess("failed to calculate fee", "", err), nil
	}

	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, ErrBaseOperationProcess("failed to check enough balance", fact.sender.String(), err), nil
	}

	for i, b := range sb {
		v, ok := b.Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(statecurrency.BalanceStateValue{}, b.Value())))
		}
		stv := statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(required[i][0])))
		sts = append(sts, currencystate.NewStateMergeValue(b.Key(), stv))
	}

	return sts, nil, nil
}

func (opp *DelegateProcessor) Close() error {
	delegateProcessorPool.Put(opp)
	return nil
}
