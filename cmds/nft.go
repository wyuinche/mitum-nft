package cmds

type NFTCommand struct {
	CollectionRegister      CollectionRegisterCommand      `cmd:"" name:"collection-register" help:"register new collection design"`
	CollectionPolicyUpdater CollectionPolicyUpdaterCommand `cmd:"" name:"collection-policy-updater" help:"update collection design"`
	Mint                    MintCommand                    `cmd:"" name:"mint" help:"mint new nft to collection"`
	NFTTransfer             NFTTransferCommand             `cmd:"" name:"nft-transfer" help:"transfer nfts to receiver"`
	Delegate                DelegateCommand                `cmd:"" name:"delegate" help:"delegate operator or cancel operator delegation"`
	Approve                 ApproveCommand                 `cmd:"" name:"approve" help:"approve account for nft"`
	NFTSign                 NFTSignCommand                 `cmd:"" name:"nft-sign" help:"sign nft as creator | copyrighter"`
}
