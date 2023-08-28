package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/pkg/errors"

	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type ApproveCommand struct {
	OperationCommand
	Approved currencycmds.AddressFlag `arg:"" name:"approved" help:"approved account address" required:"true"`
	NFTidx   uint64                   `arg:"" name:"nft" help:"target nft idx to approve"`
	approved base.Address
}

func NewApproveCommand() ApproveCommand {
	cmd := NewOperationCommand()
	return ApproveCommand{
		OperationCommand: *cmd,
	}
}

func (cmd *ApproveCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	approved, err := cmd.Approved.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved.String())
	}
	cmd.approved = approved

	return nil

}

func (cmd *ApproveCommand) createOperation() (base.Operation, error) {
	e := util.StringError(utils.ErrStringCreate("approve operation"))

	item := nft.NewApproveItem(cmd.contract, cmd.collection, cmd.approved, cmd.NFTidx, cmd.Currency.CID)

	fact := nft.NewApproveFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]nft.ApproveItem{item},
	)

	op, err := nft.NewApprove(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}

	if err := op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID()); err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
