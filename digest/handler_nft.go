package digest

import (
	"net/http"
	"strconv"
	"time"

	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-nft/v2/types"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

func (hd *Handlers) handleNFT(w http.ResponseWriter, r *http.Request) {
	cachekey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	id, err, status := parseRequest(w, r, "id")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTInGroup(contract, collection, id)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTInGroup(contract, collection, id string) (interface{}, error) {
	switch nft, err := NFT(hd.database, contract, collection, id); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTHal(contract, collection, *nft)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTHal(contract, collection string, nft types.NFT) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathNFT, "contract", contract, "collection", collection, "id", strconv.FormatUint(nft.ID(), 10))
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(nft, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleNFTCollection(w http.ResponseWriter, r *http.Request) {
	cachekey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTCollectionInGroup(contract, collection)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTCollectionInGroup(contract, collection string) (interface{}, error) {
	switch design, err := NFTCollection(hd.database, contract, collection); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTCollectionHal(contract, collection, *design)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTCollectionHal(contract, collection string, design types.Design) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathNFTCollection, "contract", contract, "collection", collection)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(design, currencydigest.NewHalLink(h, nil))

	return hal, nil
}

func (hd *Handlers) handleNFTs(w http.ResponseWriter, r *http.Request) {
	limit := currencydigest.ParseLimitQuery(r.URL.Query().Get("limit"))
	offset := currencydigest.ParseStringQuery(r.URL.Query().Get("offset"))
	reverse := currencydigest.ParseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := currencydigest.CacheKey(
		r.URL.Path, currencydigest.StringOffsetQuery(offset),
		currencydigest.StringBoolQuery("reverse", reverse),
	)

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleNFTsInGroup(contract, collection, offset, reverse, limit)

		return []interface{}{i, filled}, err
	})

	if err != nil {
		hd.Log().Err(err).Str("collection", collection).Msg("failed to get nfts")
		currencydigest.HTTP2HandleError(w, err)

		return
	}

	var b []byte
	var filled bool
	{
		l := v.([]interface{})
		b = l[0].([]byte)
		filled = l[1].(bool)
	}

	currencydigest.HTTP2WriteHalBytes(hd.encoder, w, b, http.StatusOK)

	if !shared {
		expire := hd.expireNotFilled
		if len(offset) > 0 && filled {
			expire = time.Minute
		}

		currencydigest.HTTP2WriteCache(w, cachekey, expire)
	}
}

func (hd *Handlers) handleNFTsInGroup(
	contract, collection string,
	offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("collection-nfts")
	} else {
		limit = l
	}

	var vas []currencydigest.Hal
	if err := NFTsByCollection(
		hd.database, contract, collection, reverse, offset, limit,
		func(nft types.NFT, st base.State) (bool, error) {
			hal, err := hd.buildNFTHal(contract, collection, nft)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, false, err
	} else if len(vas) < 1 {
		return nil, false, errors.Errorf("nfts not found")
	}

	i, err := hd.buildNFTsHal(contract, collection, vas, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	b, err := hd.encoder.Marshal(i)
	return b, int64(len(vas)) == limit, err
}

func (hd *Handlers) buildNFTsHal(
	contract, col string,
	vas []currencydigest.Hal,
	offset string,
	reverse bool,
) (currencydigest.Hal, error) {
	baseSelf, err := hd.combineURL(HandlerPathNFTs, "contract", contract, "collection", col)
	if err != nil {
		return nil, err
	}

	self := baseSelf
	if len(offset) > 0 {
		self = currencydigest.AddQueryValue(baseSelf, currencydigest.StringOffsetQuery(offset))
	}
	if reverse {
		self = currencydigest.AddQueryValue(baseSelf, currencydigest.StringBoolQuery("reverse", reverse))
	}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(vas, currencydigest.NewHalLink(self, nil))

	h, err := hd.combineURL(HandlerPathNFTCollection, "contract", contract, "collection", col)
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("collection", currencydigest.NewHalLink(h, nil))

	var nextoffset string

	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(types.NFT)
		nextoffset = strconv.FormatUint(va.ID(), 10)
	}

	if len(nextoffset) > 0 {
		next := baseSelf
		next = currencydigest.AddQueryValue(next, currencydigest.StringOffsetQuery(nextoffset))

		if reverse {
			next = currencydigest.AddQueryValue(next, currencydigest.StringBoolQuery("reverse", reverse))
		}

		hal = hal.AddLink("next", currencydigest.NewHalLink(next, nil))
	}

	hal = hal.AddLink("reverse", currencydigest.NewHalLink(currencydigest.AddQueryValue(baseSelf, currencydigest.StringBoolQuery("reverse", !reverse)), nil))

	return hal, nil
}

func (hd *Handlers) handleNFTOperators(w http.ResponseWriter, r *http.Request) {
	cachekey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := parseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	collection, err, status := parseRequest(w, r, "collection")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	account, err, status := parseRequest(w, r, "account")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleNFTOperatorsInGroup(contract, collection, account)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)
		if !shared {
			currencydigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleNFTOperatorsInGroup(contract, collection, account string) (interface{}, error) {
	switch operators, err := NFTOperators(hd.database, contract, collection, account); {
	case err != nil:
		return nil, err
	default:
		hal, err := hd.buildNFTOperatorsHal(contract, collection, account, *operators)
		if err != nil {
			return nil, err
		}
		return hd.encoder.Marshal(hal)
	}
}

func (hd *Handlers) buildNFTOperatorsHal(contract, collection, account string, operators types.OperatorsBook) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathNFTOperators, "contract", contract, "collection", collection, "account", account)
	if err != nil {
		return nil, err
	}

	hal := currencydigest.NewBaseHal(operators, currencydigest.NewHalLink(h, nil))

	return hal, nil
}
