package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/operation/processor"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-nft-operation-processors-map")

type processorInfo struct {
	hint      hint.Hint
	processor types.GetNewProcessor
}

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database
	var opr *currencyprocessor.OperationProcessor
	var set *hint.CompatibleSet[isaac.NewOperationProcessorInternalFunc]

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
		currencycmds.OperationProcessorContextKey, &opr,
		launch.OperationProcessorsMapContextKey, &set,
	); err != nil {
		return pctx, err
	}

	err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	if err != nil {
		return pctx, err
	}
	err = opr.SetGetNewProcessorFunc(processor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}

	ps := []processorInfo{
		{nft.CollectionRegisterHint, nft.NewCollectionRegisterProcessor()},
		{nft.CollectionPolicyUpdaterHint, nft.NewCollectionPolicyUpdaterProcessor()},
		{nft.MintHint, nft.NewMintProcessor()},
		{nft.NFTTransferHint, nft.NewNFTTransferProcessor()},
		{nft.DelegateHint, nft.NewDelegateProcessor()},
		{nft.ApproveHint, nft.NewApproveProcessor()},
		{nft.NFTSignHint, nft.NewNFTSignProcessor()},
	}

	for _, p := range ps {
		if err := opr.SetProcessor(p.hint, p.processor); err != nil {
			return pctx, err
		}

		if err := set.Add(p.hint, func(height base.Height) (base.OperationProcessor, error) {
			return opr.New(
				height,
				db.State,
				nil,
				nil,
			)
		}); err != nil {
			return pctx, err
		}
	}

	var f currencycmds.ProposalOperationFactHintFunc = IsSupportedProposalOperationFactHintFunc

	pctx = context.WithValue(pctx, currencycmds.OperationProcessorContextKey, opr)
	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, set) //revive:disable-line:modifies-parameter
	pctx = context.WithValue(pctx, currencycmds.ProposalOperationFactHintContextKey, f)

	return pctx, nil
}
