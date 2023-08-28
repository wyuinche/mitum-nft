package digest

import (
	"context"
	"fmt"
	"sync"
	"time"

	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ProtoconNet/mitum-currency/v3/digest/isaac"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"

	base "github.com/ProtoconNet/mitum2/base"
	mitumutil "github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/fixedtree"
)

var bulkWriteLimit = 500

type BlockSession struct {
	sync.RWMutex
	block                 base.BlockMap
	ops                   []base.Operation
	opstree               fixedtree.Tree
	sts                   []base.State
	st                    *currencydigest.Database
	opsTreeNodes          map[string]base.OperationFixedtreeNode
	blockModels           []mongo.WriteModel
	operationModels       []mongo.WriteModel
	accountModels         []mongo.WriteModel
	balanceModels         []mongo.WriteModel
	currencyModels        []mongo.WriteModel
	contractAccountModels []mongo.WriteModel
	nftCollectionModels   []mongo.WriteModel
	nftModels             []mongo.WriteModel
	nftBoxModels          []mongo.WriteModel
	nftOperatorModels     []mongo.WriteModel
	statesValue           *sync.Map
	balanceAddressList    []string
	nftMap                map[string]struct{}
}

func NewBlockSession(
	st *currencydigest.Database,
	blk base.BlockMap,
	ops []base.Operation,
	opstree fixedtree.Tree,
	sts []base.State,
) (*BlockSession, error) {
	if st.Readonly() {
		return nil, errors.Errorf("readonly mode")
	}

	nst, err := st.New()
	if err != nil {
		return nil, err
	}

	return &BlockSession{
		st:          nst,
		block:       blk,
		ops:         ops,
		opstree:     opstree,
		sts:         sts,
		statesValue: &sync.Map{},
		nftMap:      map[string]struct{}{},
	}, nil
}

func (bs *BlockSession) Prepare() error {
	bs.Lock()
	defer bs.Unlock()
	if err := bs.prepareOperationsTree(); err != nil {
		return err
	}
	if err := bs.prepareBlock(); err != nil {
		return err
	}
	if err := bs.prepareOperations(); err != nil {
		return err
	}
	if err := bs.prepareCurrencies(); err != nil {
		return err
	}
	if err := bs.prepareNFTs(); err != nil {
		return err
	}

	return bs.prepareAccounts()
}

func (bs *BlockSession) Commit(ctx context.Context) error {
	bs.Lock()
	defer bs.Unlock()

	started := time.Now()
	defer func() {
		bs.statesValue.Store("commit", time.Since(started))
		_ = bs.close()
	}()

	if err := bs.writeModels(ctx, defaultColNameBlock, bs.blockModels); err != nil {
		return err
	}

	if len(bs.operationModels) > 0 {
		if err := bs.writeModels(ctx, defaultColNameOperation, bs.operationModels); err != nil {
			return err
		}
	}

	if len(bs.currencyModels) > 0 {
		if err := bs.writeModels(ctx, defaultColNameCurrency, bs.currencyModels); err != nil {
			return err
		}
	}

	if len(bs.accountModels) > 0 {
		if err := bs.writeModels(ctx, defaultColNameAccount, bs.accountModels); err != nil {
			return err
		}
	}

	if len(bs.nftCollectionModels) > 0 {
		fmt.Println("||||||||||||||||||", len(bs.nftCollectionModels))
		fmt.Printf("%s", bs.nftCollectionModels[0])
		if err := bs.writeModels(ctx, defaultColNameNFTCollection, bs.nftCollectionModels); err != nil {
			return err
		}
	}

	if len(bs.nftModels) > 0 {
		for nft := range bs.nftMap {
			err := bs.st.CleanByHeightColName(
				ctx,
				bs.block.Manifest().Height(),
				defaultColNameNFT,
				"nftid",
				nft,
			)
			if err != nil {
				return err
			}
		}

		if err := bs.writeModels(ctx, defaultColNameNFT, bs.nftModels); err != nil {
			return err
		}
	}

	if len(bs.nftOperatorModels) > 0 {
		if err := bs.writeModels(ctx, defaultColNameNFTOperator, bs.nftOperatorModels); err != nil {
			return err
		}
	}

	if len(bs.nftBoxModels) > 0 {
		if err := bs.writeModels(ctx, defaultColNameNFT, bs.nftBoxModels); err != nil {
			return err
		}
	}

	if len(bs.balanceModels) > 0 {
		if err := bs.writeModels(ctx, defaultColNameBalance, bs.balanceModels); err != nil {
			return err
		}
	}

	return nil
}

func (bs *BlockSession) Close() error {
	bs.Lock()
	defer bs.Unlock()
	return bs.close()
}

func (bs *BlockSession) prepareOperationsTree() error {
	nodes := map[string]base.OperationFixedtreeNode{}

	if err := bs.opstree.Traverse(func(_ uint64, no fixedtree.Node) (bool, error) {
		nno := no.(base.OperationFixedtreeNode)
		nodes[nno.Key()] = nno

		return true, nil
	}); err != nil {
		return err
	}

	bs.opsTreeNodes = nodes

	return nil
}

func (bs *BlockSession) prepareBlock() error {
	if bs.block == nil {
		return nil
	}

	bs.blockModels = make([]mongo.WriteModel, 1)

	manifest := isaac.NewManifest(
		bs.block.Manifest().Height(),
		bs.block.Manifest().Previous(),
		bs.block.Manifest().Proposal(),
		bs.block.Manifest().OperationsTree(),
		bs.block.Manifest().StatesTree(),
		bs.block.Manifest().Suffrage(),
		bs.block.Manifest().ProposedAt(),
	)

	doc, err := currencydigest.NewManifestDoc(manifest, bs.st.DatabaseEncoder(), bs.block.Manifest().Height(), bs.ops, bs.block.SignedAt())
	if err != nil {
		return err
	}
	bs.blockModels[0] = mongo.NewInsertOneModel().SetDocument(doc)

	return nil
}

func (bs *BlockSession) prepareOperations() error {
	if len(bs.ops) < 1 {
		return nil
	}

	node := func(h mitumutil.Hash) (bool, bool, base.OperationProcessReasonError) {
		no, found := bs.opsTreeNodes[h.String()]
		if !found {
			return false, false, nil
		}

		return true, no.InState(), no.Reason()
	}

	bs.operationModels = make([]mongo.WriteModel, len(bs.ops))

	for i := range bs.ops {
		op := bs.ops[i]
		found, inState, reason := node(op.Fact().Hash())

		if !found {
			return mitumutil.ErrNotFound.Errorf("operation, %s not found in operations tree", op.Fact().Hash().String())
		}

		doc, err := currencydigest.NewOperationDoc(
			op,
			bs.st.DatabaseEncoder(),
			bs.block.Manifest().Height(),
			bs.block.SignedAt(),
			inState,
			reason,
			uint64(i),
		)
		if err != nil {
			return err
		}
		bs.operationModels[i] = mongo.NewInsertOneModel().SetDocument(doc)
	}

	return nil
}

func (bs *BlockSession) prepareAccounts() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var accountModels []mongo.WriteModel
	var balanceModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]

		switch {
		case statecurrency.IsStateAccountKey(st.Key()):
			j, err := bs.handleAccountState(st)
			if err != nil {
				return err
			}
			accountModels = append(accountModels, j...)
		case statecurrency.IsStateBalanceKey(st.Key()):
			j, address, err := bs.handleBalanceState(st)
			if err != nil {
				return err
			}
			balanceModels = append(balanceModels, j...)
			bs.balanceAddressList = append(bs.balanceAddressList, address)
		default:
			continue
		}
	}

	bs.accountModels = accountModels
	bs.balanceModels = balanceModels

	return nil
}

func (bs *BlockSession) prepareCurrencies() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var currencyModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case statecurrency.IsStateCurrencyDesignKey(st.Key()):
			j, err := bs.handleCurrencyState(st)
			if err != nil {
				return err
			}
			currencyModels = append(currencyModels, j...)
		default:
			continue
		}
	}

	bs.currencyModels = currencyModels

	return nil
}

func (bs *BlockSession) writeModels(ctx context.Context, col string, models []mongo.WriteModel) error {
	started := time.Now()
	defer func() {
		bs.statesValue.Store(fmt.Sprintf("write-models-%s", col), time.Since(started))
	}()

	n := len(models)
	if n < 1 {
		return nil
	} else if n <= bulkWriteLimit {
		return bs.writeModelsChunk(ctx, col, models)
	}

	z := n / bulkWriteLimit
	if n%bulkWriteLimit != 0 {
		z++
	}

	for i := 0; i < z; i++ {
		s := i * bulkWriteLimit
		e := s + bulkWriteLimit
		if e > n {
			e = n
		}

		if err := bs.writeModelsChunk(ctx, col, models[s:e]); err != nil {
			return err
		}
	}

	return nil
}

func (bs *BlockSession) writeModelsChunk(ctx context.Context, col string, models []mongo.WriteModel) error {
	opts := options.BulkWrite().SetOrdered(false)
	if res, err := bs.st.DatabaseClient().Collection(col).BulkWrite(ctx, models, opts); err != nil {
		return err
	} else if res != nil && res.InsertedCount < 1 {
		return errors.Errorf("not inserted to %s", col)
	}

	return nil
}

func (bs *BlockSession) close() error {
	bs.block = nil
	bs.operationModels = nil
	bs.currencyModels = nil
	bs.accountModels = nil
	bs.balanceModels = nil
	bs.contractAccountModels = nil
	bs.nftCollectionModels = nil
	bs.nftModels = nil
	bs.nftOperatorModels = nil
	bs.nftBoxModels = nil

	return bs.st.Close()
}
