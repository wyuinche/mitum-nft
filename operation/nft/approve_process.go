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

var approveItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveItemProcessor)
	},
}

var approveProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(ApproveProcessor)
	},
}

func (Approve) Process(
	ctx context.Context, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type ApproveItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   ApproveItem
}

func (ipp *ApproveItemProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError(ErrStringPreProcess(*ipp))

	it := ipp.item
	g := state.NewStateKeyGenerator(it.Contract(), it.Collection())

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(it.Approved()), getStateFunc); err != nil {
		return e.Wrap(ErrStateNotFound("approved", it.Approved().String(), err))
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return e.Wrap(ErrStateNotFound("currency", it.Currency().String(), err))
	}

	k := utils.StringerChain(it.Contract(), it.Collection())

	st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
	if err != nil {
		return e.Wrap(ErrStateNotFound("design", k, err))
	}

	design, err := state.StateDesignValue(st)
	if err != nil {
		return e.Wrap(ErrStateNotFound("design value", k, err))
	}

	if !design.Active() {
		return e.Wrap(errors.Errorf("deactivated collection, %s", it.Collection()))
	}

	st, err = currencystate.ExistsState(stateextension.StateKeyContractAccount(design.Parent()), "contract account", getStateFunc)
	if err != nil {
		return e.Wrap(ErrStateNotFound("parent", design.Parent().String(), err))
	}

	ca, err := stateextension.StateContractAccountValue(st)
	if err != nil {
		return e.Wrap(ErrStateNotFound("parent value", design.Parent().String(), err))
	}

	if !ca.IsActive() {
		return e.Wrap(errors.Errorf("deactivated contract account, %s", design.Parent()))
	}

	st, err = currencystate.ExistsState(g.NFT(it.NFTIdx()), "key of nft", getStateFunc)
	if err != nil {
		return e.Wrap(ErrStateNotFound("nft", utils.StringChain(k, fmt.Sprintf("%d", it.NFTIdx())), err))
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return e.Wrap(ErrStateNotFound("nft value", utils.StringChain(k, fmt.Sprintf("%d", it.NFTIdx())), err))
	}

	if !nv.Active() {
		return e.Wrap(errors.Errorf("burned nft, %d", it.NFTIdx()))
	}

	if it.Approved().Equal(nv.Approved()) {
		return e.Wrap(errors.Errorf("already approved, %s", it.Approved()))
	}

	if it.Approved().Equal(nv.Owner()) {
		return e.Wrap(errors.Errorf("approved account is same with owner, %s", it.Approved()))
	}

	if !nv.Owner().Equal(ipp.sender) {
		if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(nv.Owner()), getStateFunc); err != nil {
			return e.Wrap(ErrStateNotFound("nft owner", nv.Owner().String(), err))
		}

		st, err = currencystate.ExistsState(g.OperatorsBook(nv.Owner()), "key of operators", getStateFunc)
		if err != nil {
			return e.Wrap(errors.Errorf("unauthorized sender, %s: %v", ipp.sender, err))
		}

		operators, err := state.StateOperatorsBookValue(st)
		if err != nil {
			return e.Wrap(ErrStateNotFound("operators", utils.StringChain(k, nv.Owner().String()), err))
		}

		if !operators.Exists(ipp.sender) {
			return e.Wrap(errors.Errorf("unauthorized sender, %s", ipp.sender))
		}
	}

	return nil
}

func (ipp *ApproveItemProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError(ErrStringProcess(*ipp))

	it := ipp.item
	g := state.NewStateKeyGenerator(it.Contract(), it.Collection())

	st, err := currencystate.ExistsState(g.NFT(it.NFTIdx()), "key of nft", getStateFunc)
	if err != nil {
		return nil, e.Wrap(ErrStateNotFound("nft", utils.StringChain(utils.StringerChain(it.Contract(), it.Collection()), fmt.Sprintf("%d", it.NFTIdx())), err))
	}

	nv, err := state.StateNFTValue(st)
	if err != nil {
		return nil, e.Wrap(ErrStateNotFound("nft value", utils.StringChain(utils.StringerChain(it.Contract(), it.Collection()), fmt.Sprintf("%d", it.NFTIdx())), err))
	}

	n := types.NewNFT(nv.ID(), nv.Active(), nv.Owner(), nv.NFTHash(), nv.URI(), it.Approved(), nv.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, e.Wrap(err)
	}

	sts := []base.StateMergeValue{currencystate.NewStateMergeValue(st.Key(), state.NewNFTStateValue(n))}

	return sts, nil
}

func (ipp *ApproveItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = ApproveItem{}

	approveItemProcessorPool.Put(ipp)

	return nil
}

type ApproveProcessor struct {
	*base.BaseOperationProcessor
}

func NewApproveProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := ApproveProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := approveProcessorPool.Get()
		opp, ok := nopp.(*ApproveProcessor)
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

func (opp *ApproveProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(ApproveFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, BaseErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot run approve-operation", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&ApproveItemProcessor{}, ip)))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, ErrBaseOperationProcess("", "", err), nil
		}

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *ApproveProcessor) Process(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(ApproveFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(ApproveFact{}, op.Fact())))
	}

	var sts []base.StateMergeValue // nolint:prealloc
	for _, item := range fact.Items() {
		ip := approveItemProcessorPool.Get()
		ipc, ok := ip.(*ApproveItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&ApproveItemProcessor{}, ip)))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, ErrBaseOperationProcess("", "", err), nil
		}
		sts = append(sts, s...)

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

func (opp *ApproveProcessor) Close() error {
	approveProcessorPool.Put(opp)
	return nil
}
