package nft

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencystate "github.com/ProtoconNet/mitum-currency/v3/state"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

var collectionRegisterProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CollectionRegisterProcessor)
	},
}

func (CollectionRegister) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CollectionRegisterProcessor struct {
	*base.BaseOperationProcessor
}

func NewCollectionRegisterProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := CollectionRegisterProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := collectionRegisterProcessorPool.Get()
		opp, ok := nopp.(*CollectionRegisterProcessor)
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

func (opp *CollectionRegisterProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(CollectionRegisterFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(CollectionRegisterFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, BaseErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot register collection", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(fact.Currency()), getStateFunc); err != nil {
		return nil, BaseErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	st, err := currencystate.ExistsState(stateextension.StateKeyContractAccount(fact.Contract()), "key of contract", getStateFunc)
	if err != nil {
		return nil, BaseErrStateNotFound("target contract", fact.Contract().String(), err), nil
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return nil, BaseErrStateNotFound("target contract value", fact.Contract().String(), err), nil
	}

	if !ca.Owner().Equal(fact.Sender()) {
		return nil, ErrBaseOperationProcess("sender is not owner of contract", fact.Contract().String(), nil), nil
	}

	if !ca.IsActive() {
		return nil, ErrBaseOperationProcess("deactivated contract account", fact.Contract().String(), nil), nil
	}

	g := state.NewStateKeyGenerator(fact.Contract(), fact.Collection())
	k := utils.StringerChain(fact.Contract(), fact.Collection())

	if err := currencystate.CheckNotExistsState(g.Design(), getStateFunc); err != nil {
		return nil, BaseErrStateAlreadyExists("design", k, err), nil
	}

	if err := currencystate.CheckNotExistsState(g.LastNFTIndex(), getStateFunc); err != nil {
		return nil, BaseErrStateAlreadyExists("last index of collection design", k, err), nil
	}

	whitelist := fact.Whitelist()
	for _, w := range whitelist {
		if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(w), getStateFunc); err != nil {
			return nil, BaseErrStateNotFound("whitelist account", w.String(), err), nil
		} else if err = currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(w), getStateFunc); err != nil {
			return nil, ErrBaseOperationProcess("whitelist account is contract account", w.String(), err), nil
		}
	}

	return ctx, nil, nil
}

func (opp *CollectionRegisterProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(CollectionRegisterFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(CollectionRegisterFact{}, op.Fact())))
	}

	sts := make([]base.StateMergeValue, 3)

	g := state.NewStateKeyGenerator(fact.Contract(), fact.Collection())

	policy := types.NewPolicy(fact.Name(), fact.Royalty(), fact.URI(), fact.Whitelist())
	design := types.NewDesign(fact.Contract(), fact.Sender(), fact.Collection(), true, policy)
	if err := design.IsValid(nil); err != nil {
		return nil, BaseErrInvalid(design, err), nil
	}

	sts[0] = currencystate.NewStateMergeValue(
		g.Design(),
		state.NewDesignStateValue(design),
	)
	sts[1] = currencystate.NewStateMergeValue(
		g.LastNFTIndex(),
		state.NewLastNFTIndexStateValue(0),
	)

	currencyPolicy, err := currencystate.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, BaseErrStateNotFound("currency", fact.Currency().String(), err), nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, ErrBaseOperationProcess("failed to check currency fee", fact.Currency().String(), err), nil
	}

	st, err := currencystate.ExistsState(statecurrency.StateKeyBalance(fact.Sender(), fact.Currency()), "key of balance", getStateFunc)
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

	sts[2] = currencystate.NewStateMergeValue(
		sb.Key(),
		statecurrency.NewBalanceStateValue(v.Amount.WithBig(v.Amount.Big().Sub(fee))),
	)
	return sts, nil, nil
}

func (opp *CollectionRegisterProcessor) Close() error {
	collectionRegisterProcessorPool.Put(opp)

	return nil
}
