package digest

import (
	"context"
	"strconv"

	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"

	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var maxLimit int64 = 50

var (
	defaultColNameAccount       = "digest_ac"
	defaultColNameBalance       = "digest_bl"
	defaultColNameCurrency      = "digest_cr"
	defaultColNameOperation     = "digest_op"
	defaultColNameBlock         = "digest_bm"
	defaultColNameNFTCollection = "digest_nftcollection"
	defaultColNameNFT           = "digest_nft"
	defaultColNameNFTOperator   = "digest_nftoperator"
)

func NFTCollection(st *currencydigest.Database, contract, col string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("collection", col)

	var design *types.Design
	var sta base.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameNFTCollection,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			d, err := state.StateDesignValue(sta)
			if err != nil {
				return err
			}
			design = &d

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return design, nil
}

func NFT(st *currencydigest.Database, contract, col, idx string) (*types.NFT, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("collection", col)
	filter = filter.Add("nftid", idx)

	var nft *types.NFT
	var sta base.State
	var err error
	if err = st.DatabaseClient().GetByFilter(
		defaultColNameNFT,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}
			n, err := state.StateNFTValue(sta)
			if err != nil {
				return err
			}
			nft = &n

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return nft, nil
}

func NFTsByAddress(
	st *currencydigest.Database,
	address base.Address,
	reverse bool,
	offset string,
	limit int64,
	collectionid string,
	callback func(string /* nft id */, types.NFT) (bool, error),
) error {
	filter, err := buildNFTsFilterByAddress(address, offset, reverse, collectionid)
	if err != nil {
		return err
	}

	sr := 1
	if reverse {
		sr = -1
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("height", sr).D(),
	)

	switch {
	case limit <= 0: // no limit
	case limit > maxLimit:
		opt = opt.SetLimit(maxLimit)
	default:
		opt = opt.SetLimit(limit)
	}

	return st.DatabaseClient().Find(
		context.Background(),
		defaultColNameNFT,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.DatabaseEncoders())
			if err != nil {
				return false, err
			}
			nft, err := state.StateNFTValue(st)
			if err != nil {
				return false, err
			}

			return callback(strconv.FormatUint(nft.ID(), 10), nft)
		},
		opt,
	)
}

func NFTsByCollection(
	st *currencydigest.Database,
	contract,
	col string,
	reverse bool,
	offset string,
	limit int64,
	callback func(nft types.NFT, st base.State) (bool, error),
) error {
	filter, err := buildNFTsFilterByCollection(contract, col, offset, reverse)
	if err != nil {
		return err
	}

	sr := 1
	if reverse {
		sr = -1
	}

	opt := options.Find().SetSort(
		util.NewBSONFilter("height", sr).D(),
	)

	switch {
	case limit <= 0: // no limit
	case limit > maxLimit:
		opt = opt.SetLimit(maxLimit)
	default:
		opt = opt.SetLimit(limit)
	}

	return st.DatabaseClient().Find(
		context.Background(),
		defaultColNameNFT,
		filter,
		func(cursor *mongo.Cursor) (bool, error) {
			st, err := currencydigest.LoadState(cursor.Decode, st.DatabaseEncoders())
			if err != nil {
				return false, err
			}
			nft, err := state.StateNFTValue(st)
			if err != nil {
				return false, err
			}
			return callback(nft, st)
		},
		opt,
	)
}

func NFTOperators(
	st *currencydigest.Database,
	contract, col, account string,
) (*types.OperatorsBook, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("collection", col)
	filter = filter.Add("account", account)

	var operators *types.OperatorsBook
	var sta base.State
	var err error
	if err := st.DatabaseClient().GetByFilter(
		defaultColNameNFTOperator,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = currencydigest.LoadState(res.Decode, st.DatabaseEncoders())
			if err != nil {
				return err
			}

			ops, err := state.StateOperatorsBookValue(sta)
			if err != nil {
				return err
			}
			operators = &ops

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, err
	}

	return operators, nil
}
