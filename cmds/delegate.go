package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type DelegateCommand struct {
	OperationCommand
	Operator currencycmds.AddressFlag `arg:"" name:"operator" help:"operator account address"`
	Mode     string                   `name:"mode" help:"delegate mode" optional:""`
	operator base.Address
	mode     nft.DelegateMode
}

func NewDelegateCommand() DelegateCommand {
	cmd := NewOperationCommand()
	return DelegateCommand{OperationCommand: *cmd}
}

func (cmd *DelegateCommand) Run(pctx context.Context) error {
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

func (cmd *DelegateCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	operator, err := cmd.Operator.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid operator format, %q", cmd.Operator.String())
	}
	cmd.operator = operator

	if len(cmd.Mode) < 1 {
		cmd.mode = nft.DelegateAllow
	} else {
		mode := nft.DelegateMode(cmd.Mode)
		if err := mode.IsValid(nil); err != nil {
			return err
		}
		cmd.mode = mode
	}

	return nil

}

func (cmd *DelegateCommand) createOperation() (base.Operation, error) {
	e := util.StringError(utils.ErrStringCreate("delegate operation"))

	items := []nft.DelegateItem{nft.NewDelegateItem(cmd.contract, cmd.collection, cmd.operator, cmd.mode, cmd.Currency.CID)}

	fact := nft.NewDelegateFact([]byte(cmd.Token), cmd.sender, items)

	op, err := nft.NewDelegate(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
