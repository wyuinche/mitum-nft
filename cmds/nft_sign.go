package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/utils"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type NFTSignCommand struct {
	OperationCommand
	NFT uint64 `arg:"" name:"nft" help:"target nft; \"<collection>,<idx>\""`
}

func NewNFTSignCommand() NFTSignCommand {
	cmd := NewOperationCommand()
	return NFTSignCommand{OperationCommand: *cmd}
}

func (cmd *NFTSignCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *NFTSignCommand) parseFlags() error {
	if err := cmd.OperationCommand.parseFlags(); err != nil {
		return err
	}

	col := types.ContractID(cmd.Collection)
	if err := col.IsValid(nil); err != nil {
		return err
	} else {
		cmd.collection = col
	}

	return nil

}

func (cmd *NFTSignCommand) createOperation() (base.Operation, error) {
	e := util.StringError(utils.ErrStringCreate("nft-sign operation"))

	item := nft.NewNFTSignItem(cmd.contract, cmd.collection, cmd.NFT, cmd.Currency.CID)
	fact := nft.NewNFTSignFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]nft.NFTSignItem{item},
	)

	op, err := nft.NewNFTSign(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}

	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
