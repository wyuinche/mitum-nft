package nft

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"

	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var collectionCollectionPolicyUpdaterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CollectionPolicyUpdaterProcessor)
	},
}

func (CollectionPolicyUpdater) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CollectionPolicyUpdaterProcessor struct {
	*base.BaseOperationProcessor
}

func NewCollectionPolicyUpdaterProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := CollectionPolicyUpdaterProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := collectionCollectionPolicyUpdaterProcessorPool.Get()
		opp, ok := nopp.(*CollectionPolicyUpdaterProcessor)
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

func (opp *CollectionPolicyUpdaterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(CollectionPolicyUpdaterFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(CollectionPolicyUpdaterFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, BaseErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot update policy", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, BaseErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract(), fact.Collection())
	k := utils.StringerChain(fact.Contract(), fact.Collection())

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

	if !design.Creator().Equal(fact.Sender()) {
		return nil, ErrBaseOperationProcess("not creator of collection design", k, nil), nil
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

	return ctx, nil, nil
}

func (opp *CollectionPolicyUpdaterProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(CollectionPolicyUpdaterFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(CollectionPolicyUpdaterFact{}, op.Fact())))
	}

	g := state.NewStateKeyGenerator(fact.contract, fact.collection)
	k := utils.StringerChain(fact.contract, fact.collection)

	st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
	if err != nil {
		return nil, BaseErrStateNotFound("design", k, err), nil
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return nil, BaseErrStateNotFound("design value", k, err), nil
	}

	sts := make([]base.StateMergeValue, 2)

	d := types.NewDesign(
		design.Parent(),
		design.Creator(),
		design.Collection(),
		design.Active(),
		types.NewPolicy(fact.name, fact.royalty, fact.uri, fact.whitelist),
	)
	if err := d.IsValid(nil); err != nil {
		return nil, BaseErrInvalid(d, err), nil
	}

	sts[0] = currencystate.NewStateMergeValue(g.Design(), state.NewDesignStateValue(d))

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, BaseErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, ErrBaseOperationProcess("failed to check currency fee", fact.Currency().String(), err), nil
	}

	st, err = currencystate.ExistsState(statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of balance", getStateFunc)
	if err != nil {
		return nil, BaseErrStateNotFound("balance", utils.StringerChain(fact.Sender(), fact.Currency()), err), nil
	}
	sb := currencystate.NewStateMergeValue(st.Key(), st.Value())

	switch b, err := statecurrency.StateBalanceValue(st); {
	case err != nil:
		return nil, BaseErrStateNotFound("balance value", utils.StringerChain(fact.Sender(), fact.Currency()), err), nil
	case b.Big().Compare(fee) < 0:
		return nil, ErrBaseOperationProcess("not enough balance of sender", fact.Sender().String(), err), nil
	}

	v, ok := sb.Value().(statecurrency.BalanceStateValue)
	if !ok {
		return nil, ErrBaseOperationProcess(utils.ErrStringTypeCast(statecurrency.BalanceStateValue{}, sb.Value()), utils.StringerChain(fact.Sender(), fact.Currency()), nil), nil
	}

	sts[1] = currencystate.NewStateMergeValue(
		sb.Key(),
		statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)

	return sts, nil, nil
}

func (opp *CollectionPolicyUpdaterProcessor) Close() error {
	collectionCollectionPolicyUpdaterProcessorPool.Put(opp)
	return nil
}
