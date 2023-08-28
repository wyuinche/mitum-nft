package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	base "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type CollectionRegisterCommand struct {
	OperationCommand
	Name      string                   `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty   uint                     `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	URI       string                   `name:"uri" help:"collection uri" optional:""`
	White     currencycmds.AddressFlag `name:"white" help:"whitelisted address" optional:""`
	name      types.CollectionName
	royalty   types.PaymentParameter
	uri       types.URI
	whitelist []base.Address
}

func NewCollectionRegisterCommand() CollectionRegisterCommand {
	cmd := NewOperationCommand()
	return CollectionRegisterCommand{OperationCommand: *cmd}
}

func (cmd *CollectionRegisterCommand) Run(pctx context.Context) error {
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

func (cmd *CollectionRegisterCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	var white base.Address = nil
	if cmd.White.String() != "" {
		if a, err := cmd.White.Encode(enc); err != nil {
			return errors.Wrapf(err, "invalid whitelist address format, %q", cmd.White)
		} else {
			white = a
		}
	}

	name := types.CollectionName(cmd.Name)
	if err := name.IsValid(nil); err != nil {
		return err
	} else {
		cmd.name = name
	}

	royalty := types.PaymentParameter(cmd.Royalty)
	if err := royalty.IsValid(nil); err != nil {
		return err
	} else {
		cmd.royalty = royalty
	}

	uri := types.URI(cmd.URI)
	if err := uri.IsValid(nil); err != nil {
		return err
	} else {
		cmd.uri = uri
	}

	whitelist := []base.Address{}
	if white != nil {
		whitelist = append(whitelist, white)
	}
	cmd.whitelist = whitelist

	return nil
}

func (cmd *CollectionRegisterCommand) createOperation() (base.Operation, error) {
	e := util.StringError(utils.ErrStringCreate("collection-register operation"))

	fact := nft.NewCollectionRegisterFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.collection,
		cmd.name,
		cmd.royalty,
		cmd.uri,
		cmd.whitelist,
		cmd.Currency.CID,
	)

	op, err := nft.NewCollectionRegister(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}

	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
