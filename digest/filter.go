package digest

import (
	"net/http"
	"strings"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func parseRequest(w http.ResponseWriter, r *http.Request, v string) (string, error, int) {
	s, found := mux.Vars(r)[v]
	if !found {
		return "", errors.Errorf("empty collection id"), http.StatusNotFound
	}

	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return "", errors.Errorf("empty collection id"), http.StatusBadRequest
	}
	return s, nil, http.StatusOK
}

func buildNFTsFilterByAddress(address base.Address, offset string, reverse bool, collection string) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching address
	filterAddress := bson.D{{"owner", bson.D{{"$in", []string{address.String()}}}}}
	filterA = append(filterA, filterAddress)

	// if collection query exist, find by collection first
	if len(collection) > 0 {
		filterCollection := bson.D{
			{"collection", bson.D{{"$eq", collection}}},
		}
		filterA = append(filterA, filterCollection)
	}

	// if offset exist, apply offset
	if len(offset) > 0 {
		if !reverse {
			filterOffset := bson.D{
				{"nftid", bson.D{{"$gt", offset}}},
			}
			filterA = append(filterA, filterOffset)
			// if reverse true, lesser then offset height
		} else {
			filterHeight := bson.D{
				{"nftid", bson.D{{"$lt", offset}}},
			}
			filterA = append(filterA, filterHeight)
		}
	}

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	return filter, nil
}

func buildNFTsFilterByCollection(contract, col string, offset string, reverse bool) (bson.D, error) {
	filterA := bson.A{}

	// filter fot matching collection
	filterSymbol := bson.D{{"collection", bson.D{{"$in", []string{col}}}}}
	filterToken := bson.D{{"istoken", true}}
	filterA = append(filterA, filterToken)
	filterA = append(filterA, filterSymbol)

	// if offset exist, apply offset
	if len(offset) > 0 {
		if !reverse {
			filterOffset := bson.D{
				{"nftid", bson.D{{"$gt", offset}}},
			}
			filterA = append(filterA, filterOffset)
			// if reverse true, lesser then offset height
		} else {
			filterHeight := bson.D{
				{"nftid", bson.D{{"$lt", offset}}},
			}
			filterA = append(filterA, filterHeight)
		}
	}

	filter := bson.D{}
	if len(filterA) > 0 {
		filter = bson.D{
			{"$and", filterA},
		}
	}

	return filter, nil
}
