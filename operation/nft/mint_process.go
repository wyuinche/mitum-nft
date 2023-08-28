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

var mintItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MintItemProcessor)
	},
}

var mintProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(MintProcessor)
	},
}

func (Mint) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type MintItemProcessor struct {
	h      util.Hash
	sender base.Address
	item   MintItem
	idx    uint64
	box    *types.NFTBox
}

func (ipp *MintItemProcessor) PreProcess(
	_ context.Context, _ base.Operation, getStateFunc base.GetStateFunc,
) error {
	e := util.StringError(ErrStringPreProcess(*ipp))

	it := ipp.item
	g := state.NewStateKeyGenerator(it.Contract(), it.Collection())

	k := utils.StringerChain(it.Contract(), it.Collection())

	if err := currencystate.CheckNotExistsState(g.NFT(ipp.idx), getStateFunc); err != nil {
		return e.Wrap(ErrStateAlreadyExists("nft", utils.StringChain(k, fmt.Sprintf("%d", ipp.idx)), err))
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyCurrencyDesign(it.Currency()), getStateFunc); err != nil {
		return e.Wrap(ErrStateNotFound("currency", it.Currency().String(), err))
	}

	if it.Creators().Total() != 0 {
		creators := it.Creators().Signers()
		for _, creator := range creators {
			ac := creator.Account()
			if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(ac), getStateFunc); err != nil {
				return e.Wrap(ErrStateNotFound("creator", ac.String(), err))
			}
			if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(ac), getStateFunc); err != nil {
				return e.Wrap(errors.Errorf("contract account cannot be a creator, %s: %v", ac, err))
			}
			if creator.Signed() {
				return e.Wrap(errors.Errorf("cannot sign at the same time as minting, %s", ac))
			}
		}
	}

	return nil
}

func (ipp *MintItemProcessor) Process(
	_ context.Context, _ base.Operation, _ base.GetStateFunc,
) ([]base.StateMergeValue, error) {
	e := util.StringError(ErrStringProcess(*ipp))
	it := ipp.item

	sts := make([]base.StateMergeValue, 1)

	n := types.NewNFT(ipp.idx, true, ipp.sender, it.NFTHash(), it.URI(), ipp.sender, it.Creators())
	if err := n.IsValid(nil); err != nil {
		return nil, e.Wrap(ErrInvalid(n, err))
	}

	sts[0] = currencystate.NewStateMergeValue(state.StateKeyNFT(it.Contract(), it.Collection(), ipp.idx), state.NewNFTStateValue(n))

	if err := ipp.box.Append(n.ID()); err != nil {
		return nil, e.Wrap(errors.Errorf("failed to append nft id to nft box, %d: %v", n.ID(), err))
	}

	return sts, nil
}

func (ipp *MintItemProcessor) Close() error {
	ipp.h = nil
	ipp.sender = nil
	ipp.item = MintItem{}
	ipp.idx = 0
	ipp.box = nil

	mintItemProcessorPool.Put(ipp)

	return nil
}

type MintProcessor struct {
	*base.BaseOperationProcessor
}

func NewMintProcessor() currencytypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		t := MintProcessor{}
		e := util.StringError(utils.ErrStringCreate(fmt.Sprintf("new %T", t)))

		nopp := mintProcessorPool.Get()
		opp, ok := nopp.(*MintProcessor)
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

func (opp *MintProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	e := util.StringError(ErrStringPreProcess(*opp))

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return ctx, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(MintFact{}, op.Fact())))
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, nil, e.Wrap(err)
	}

	if err := currencystate.CheckExistsState(statecurrency.StateKeyAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, BaseErrStateNotFound("sender", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckNotExistsState(stateextension.StateKeyContractAccount(fact.Sender()), getStateFunc); err != nil {
		return nil, ErrBaseOperationProcess("contract account cannot mint nfts", fact.Sender().String(), err), nil
	}

	if err := currencystate.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, ErrBaseOperationProcess("invalid signing", "", err), nil
	}

	idxes := map[currencytypes.ContractID]uint64{}
	for _, item := range fact.Items() {
		g := state.NewStateKeyGenerator(item.Contract(), item.Collection())
		k := utils.StringerChain(item.Contract(), item.Collection())

		collection := item.Collection()

		if _, found := idxes[collection]; !found {
			st, err := currencystate.ExistsState(g.Design(), "key of design", getStateFunc)
			if err != nil {
				return nil, BaseErrStateNotFound("design", k, err), nil
			}

			design, err := state.StateDesignValue(st)
			if err != nil {
				return nil, BaseErrStateNotFound("design value", k, err), nil
			}

			if !design.Active() {
				return nil, base.NewBaseOperationProcessReasonError("deactivated collection, %q", collection), nil
			}

			st, err = currencystate.ExistsState(stateextension.StateKeyContractAccount(design.Parent()), "key of contract", getStateFunc)
			if err != nil {
				return nil, BaseErrStateNotFound("parent", design.Parent().String(), err), nil
			}

			parent, err := stateextension.StateContractAccountValue(st)
			if err != nil {
				return nil, BaseErrStateNotFound("parent value", design.Parent().String(), err), nil
			}

			if !parent.Owner().Equal(fact.Sender()) {
				return nil, ErrBaseOperationProcess("sender is not owner of contract", design.Parent().String(), nil), nil
			}

			if !parent.IsActive() {
				return nil, ErrBaseOperationProcess("deactivated contract account", design.Parent().String(), nil), nil
			}

			st, err = currencystate.ExistsState(g.LastNFTIndex(), "key of collection last index", getStateFunc)
			if err != nil {
				return nil, BaseErrStateNotFound("collection last index", k, err), nil
			}

			nftID, err := state.StateLastNFTIndexValue(st)
			if err != nil {
				return nil, BaseErrStateNotFound("collection last index value", k, err), nil
			}

			idxes[collection] = nftID
		}
	}

	for _, item := range fact.Items() {
		ip := mintItemProcessorPool.Get()
		ipc, ok := ip.(*MintItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&MintItemProcessor{}, ip)))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.idx = idxes[item.Collection()]
		ipc.box = nil

		if err := ipc.PreProcess(ctx, op, getStateFunc); err != nil {
			return nil, ErrBaseOperationProcess("", "", err), nil
		}
		idxes[item.Collection()] += 1

		ipc.Close()
	}

	return ctx, nil, nil
}

func (opp *MintProcessor) Process( // nolint:dupl
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	e := util.StringError(ErrStringProcess(*opp))

	fact, ok := op.Fact().(MintFact)
	if !ok {
		return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(MintFact{}, op.Fact())))
	}

	idxes := map[string]uint64{}
	boxes := map[string]*types.NFTBox{}

	for _, item := range fact.items {
		g := state.NewStateKeyGenerator(item.Contract(), item.Collection())
		k := g.OperatorsBook(fact.Sender())

		ik := g.LastNFTIndex()
		if _, found := idxes[ik]; !found {
			st, err := currencystate.ExistsState(ik, "key of collection last index", getStateFunc)
			if err != nil {
				return nil, BaseErrStateNotFound("collection last index", k, err), nil
			}

			nftID, err := state.StateLastNFTIndexValue(st)
			if err != nil {
				return nil, BaseErrStateNotFound("collection last index value", k, err), nil
			}

			idxes[ik] = nftID
		}

		bk := g.NFTBox()
		if _, found := boxes[bk]; !found {
			var box types.NFTBox

			switch st, found, err := getStateFunc(bk); {
			case err != nil:
				return nil, BaseErrStateNotFound("nft box", k, err), nil
			case !found:
				box = types.NewNFTBox(nil)
			default:
				b, err := state.StateNFTBoxValue(st)
				if err != nil {
					return nil, BaseErrStateNotFound("nft box value", k, err), nil
				}
				box = b
			}

			boxes[bk] = &box
		}
	}

	var sts []base.StateMergeValue // nolint:prealloc

	ipcs := make([]*MintItemProcessor, len(fact.Items()))
	for i, item := range fact.Items() {
		g := state.NewStateKeyGenerator(item.Contract(), item.Collection())

		ik := g.LastNFTIndex()
		bk := g.NFTBox()

		ip := mintItemProcessorPool.Get()
		ipc, ok := ip.(*MintItemProcessor)
		if !ok {
			return nil, nil, e.Wrap(errors.Errorf(utils.ErrStringTypeCast(&MintItemProcessor{}, ip)))
		}

		ipc.h = op.Hash()
		ipc.sender = fact.Sender()
		ipc.item = item
		ipc.idx = idxes[ik]
		ipc.box = boxes[bk]

		s, err := ipc.Process(ctx, op, getStateFunc)
		if err != nil {
			return nil, ErrBaseOperationProcess("", "", err), nil
		}
		sts = append(sts, s...)

		idxes[ik] += 1
		ipcs[i] = ipc
	}

	for key, idx := range idxes {
		iv := currencystate.NewStateMergeValue(key, state.NewLastNFTIndexStateValue(idx))
		sts = append(sts, iv)
	}

	for key, box := range boxes {
		bv := currencystate.NewStateMergeValue(key, state.NewNFTBoxStateValue(*box))
		sts = append(sts, bv)
	}

	for _, ipc := range ipcs {
		ipc.Close()
	}

	idxes = nil
	boxes = nil

	required, err := CalculateItemsFee(getStateFunc, fact.items)
	if err != nil {
		return nil, ErrBaseOperationProcess("failed to calculate fee", "", err), nil
	}

	sb, err := currency.CheckEnoughBalance(fact.sender, required, getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("failed to check enough balance; %w", err), nil
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

func (opp *MintProcessor) Close() error {
	mintProcessorPool.Put(opp)

	return nil
}
