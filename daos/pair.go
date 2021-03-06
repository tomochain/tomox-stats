package daos

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-stats/app"
	"github.com/tomochain/tomox-stats/types"
)

// PairDao contains:
// collectionName: MongoDB collection name
// dbName: name of mongodb to interact with
type PairDao struct {
	collectionName string
	dbName         string
}

type PairDaoOption = func(*PairDao) error

func PairDaoDBOption(dbName string) func(dao *PairDao) error {
	return func(dao *PairDao) error {
		dao.dbName = dbName
		return nil
	}
}

// NewPairDao returns a new instance of AddressDao
func NewPairDao(options ...PairDaoOption) *PairDao {
	dao := &PairDao{}
	dao.collectionName = "pairs"
	dao.dbName = app.Config.DBName

	for _, op := range options {
		err := op(dao)
		if err != nil {
			panic(err)
		}
	}

	index := mgo.Index{
		Key:    []string{"baseTokenAddress", "quoteTokenAddress", "relayerAddress"},
		Unique: true,
	}

	err := db.Session.DB(dao.dbName).C(dao.collectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	return dao
}

// Create function performs the DB insertion task for pair collection
func (dao *PairDao) Create(pair *types.Pair) error {
	pair.ID = bson.NewObjectId()
	pair.CreatedAt = time.Now()
	pair.UpdatedAt = time.Now()

	err := db.Create(dao.dbName, dao.collectionName, pair)
	return err
}

// GetAll function fetches all the pairs in the pair collection of mongodb.
// for GetAll return continous memory
func (dao *PairDao) GetAll() ([]types.Pair, error) {
	var res []types.Pair
	err := db.Get(dao.dbName, dao.collectionName, bson.M{}, 0, 0, &res)
	if err != nil {
		return nil, err
	}

	ret := []types.Pair{}
	keys := make(map[string]bool)

	for _, it := range res {
		code := it.BaseTokenAddress.Hex() + "::" + it.QuoteTokenAddress.Hex()
		if _, value := keys[code]; !value {
			keys[code] = true
			ret = append(ret, it)
		}
	}

	return ret, nil
}

// GetAllByCoinbase get pair by coinbase address
func (dao *PairDao) GetAllByCoinbase(addr common.Address) ([]types.Pair, error) {
	var res []types.Pair
	err := db.Get(dao.dbName, dao.collectionName, bson.M{"relayerAddress": addr.Hex()}, 0, 0, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetActivePairsByCoinbase get active pair by coinbase address
func (dao *PairDao) GetActivePairsByCoinbase(addr common.Address) ([]*types.Pair, error) {
	var res []*types.Pair

	q := bson.M{"active": true, "relayerAddress": addr.Hex()}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	ret := []*types.Pair{}
	keys := make(map[string]bool)

	for _, it := range res {
		code := it.BaseTokenAddress.Hex() + "::" + it.QuoteTokenAddress.Hex()
		if _, value := keys[code]; !value {
			keys[code] = true
			ret = append(ret, it)
		}
	}

	return ret, nil
}

// GetActivePairs get active pair current coinbase
func (dao *PairDao) GetActivePairs() ([]*types.Pair, error) {
	var res []*types.Pair

	q := bson.M{"active": true}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 0, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	ret := []*types.Pair{}
	keys := make(map[string]bool)

	for _, it := range res {
		code := it.BaseTokenAddress.Hex() + "::" + it.QuoteTokenAddress.Hex()
		if _, value := keys[code]; !value {
			keys[code] = true
			ret = append(ret, it)
		}
	}

	return ret, nil
}

// GetByID function fetches details of a pair using pair's mongo ID.
func (dao *PairDao) GetByID(id bson.ObjectId) (*types.Pair, error) {
	var response *types.Pair
	err := db.GetByID(dao.dbName, dao.collectionName, id, &response)
	return response, err
}

// GetByName function fetches details of a pair using pair's name.
// It makes CASE INSENSITIVE search query one pair's name
func (dao *PairDao) GetByName(name string) (*types.Pair, error) {

	tokenSymbols := strings.Split(name, "/")
	return dao.GetByTokenSymbols(tokenSymbols[0], tokenSymbols[1])
}

// GetByTokenSymbols get token by symbol
func (dao *PairDao) GetByTokenSymbols(baseTokenSymbol, quoteTokenSymbol string) (*types.Pair, error) {
	var res []*types.Pair

	q := bson.M{
		"baseTokenSymbol":  baseTokenSymbol,
		"quoteTokenSymbol": quoteTokenSymbol,
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

// GetByTokenAddress function fetches pair based on
// CONTRACT ADDRESS of base token and quote token
func (dao *PairDao) GetByTokenAddress(baseToken, quoteToken common.Address) (*types.Pair, error) {
	var res []*types.Pair

	q := bson.M{
		"baseTokenAddress":  baseToken.Hex(),
		"quoteTokenAddress": quoteToken.Hex(),
	}

	err := db.Get(dao.dbName, dao.collectionName, q, 0, 1, &res)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

// DeleteByToken delete token by contract address
func (dao *PairDao) DeleteByToken(baseAddress common.Address, quoteAddress common.Address) error {
	query := bson.M{"baseTokenAddress": baseAddress.Hex(), "quoteTokenAddress": quoteAddress.Hex()}
	return db.RemoveItem(dao.dbName, dao.collectionName, query)
}

// DeleteByTokenAndCoinbase delete token by coinbase
func (dao *PairDao) DeleteByTokenAndCoinbase(baseAddress common.Address, quoteAddress common.Address, addr common.Address) error {
	query := bson.M{"baseTokenAddress": baseAddress.Hex(), "quoteTokenAddress": quoteAddress.Hex(), "relayerAddress": addr.Hex()}
	return db.RemoveItem(dao.dbName, dao.collectionName, query)
}
