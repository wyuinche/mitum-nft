package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: types.SignerHint, Instance: types.Signer{}},
	{Hint: types.SignersHint, Instance: types.Signers{}},
	{Hint: types.NFTHint, Instance: types.NFT{}},
	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.OperatorsBookHint, Instance: types.OperatorsBook{}},
	{Hint: types.PolicyHint, Instance: types.Policy{}},
	{Hint: types.NFTBoxHint, Instance: types.NFTBox{}},

	{Hint: nft.CollectionRegisterHint, Instance: nft.CollectionRegister{}},
	{Hint: nft.CollectionPolicyUpdaterHint, Instance: nft.CollectionPolicyUpdater{}},
	{Hint: nft.MintItemHint, Instance: nft.MintItem{}},
	{Hint: nft.MintHint, Instance: nft.Mint{}},
	{Hint: nft.NFTTransferItemHint, Instance: nft.NFTTransferItem{}},
	{Hint: nft.NFTTransferHint, Instance: nft.NFTTransfer{}},
	{Hint: nft.DelegateItemHint, Instance: nft.DelegateItem{}},
	{Hint: nft.DelegateHint, Instance: nft.Delegate{}},
	{Hint: nft.ApproveItemHint, Instance: nft.ApproveItem{}},
	{Hint: nft.ApproveHint, Instance: nft.Approve{}},
	{Hint: nft.NFTSignItemHint, Instance: nft.NFTSignItem{}},
	{Hint: nft.NFTSignHint, Instance: nft.NFTSign{}},

	{Hint: state.LastNFTIndexStateValueHint, Instance: state.LastNFTIndexStateValue{}},
	{Hint: state.NFTStateValueHint, Instance: state.NFTStateValue{}},
	{Hint: state.NFTBoxStateValueHint, Instance: state.NFTBoxStateValue{}},
	{Hint: state.OperatorsBookStateValueHint, Instance: state.OperatorsBookStateValue{}},
	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: nft.CollectionRegisterFactHint, Instance: nft.CollectionRegisterFact{}},
	{Hint: nft.CollectionPolicyUpdaterFactHint, Instance: nft.CollectionPolicyUpdaterFact{}},
	{Hint: nft.MintFactHint, Instance: nft.MintFact{}},
	{Hint: nft.NFTTransferFactHint, Instance: nft.NFTTransferFact{}},
	{Hint: nft.DelegateFactHint, Instance: nft.DelegateFact{}},
	{Hint: nft.ApproveFactHint, Instance: nft.ApproveFact{}},
	{Hint: nft.NFTSignFactHint, Instance: nft.NFTSignFact{}},
}

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	allExtendedLen := currencyExtendedLen + len(AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:], AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	allSupportedExtendedLen := currencySupportedExtendedLen + len(AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:], AddedSupportedHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
