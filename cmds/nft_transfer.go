package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type NFTTransferCommand struct {
	OperationCommand
	Receiver currencycmds.AddressFlag `arg:"" name:"receiver" help:"nft owner" required:"true"`
	NFT      uint64                   `arg:"" name:"nft" help:"target nft"`
	receiver base.Address
}

func NewNFTTranfserCommand() NFTTransferCommand {
	cmd := NewOperationCommand()
	return NFTTransferCommand{OperationCommand: *cmd}
}

func (cmd *NFTTransferCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.Encoders
	enc = cmd.Encoder

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *NFTTransferCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	if a, err := cmd.Receiver.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	} else {
		cmd.receiver = a
	}

	col := types.ContractID(cmd.Collection)
	if err := col.IsValid(nil); err != nil {
		return err
	} else {
		cmd.collection = col
	}

	return nil

}

func (cmd *NFTTransferCommand) createOperation() (base.Operation, error) {
	e := util.StringError(utils.ErrStringCreate("nft-transfer operation"))

	item := nft.NewNFTTransferItem(cmd.contract, cmd.collection, cmd.receiver, cmd.NFT, cmd.Currency.CID)
	fact := nft.NewNFTTransferFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]nft.NFTTransferItem{item},
	)

	op, err := nft.NewNFTTransfer(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}

	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
