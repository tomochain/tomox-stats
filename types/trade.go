package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/tomochain/tomox-stats/utils/math"

	"github.com/globalsign/mgo/bson"
)

const (
	TradeStatusPending = "PENDING"
	TradeStatusSuccess = "SUCCESS"
	TradeStatusError   = "ERROR"
)

// Trade struct holds arguments corresponding to a "Taker Order"
// To be valid an accept by the matching engine (and ultimately the exchange smart-contract),
// the trade signature must be made from the trader Maker account
type Trade struct {
	ID             bson.ObjectId  `json:"id,omitempty" bson:"_id"`
	Taker          common.Address `json:"taker" bson:"taker"`
	Maker          common.Address `json:"maker" bson:"maker"`
	BaseToken      common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken     common.Address `json:"quoteToken" bson:"quoteToken"`
	MakerOrderHash common.Hash    `json:"makerOrderHash" bson:"makerOrderHash"`
	TakerOrderHash common.Hash    `json:"takerOrderHash" bson:"takerOrderHash"`
	Hash           common.Hash    `json:"hash" bson:"hash"`
	TxHash         common.Hash    `json:"txHash" bson:"txHash"`
	PairName       string         `json:"pairName" bson:"pairName"`
	PricePoint     *big.Int       `json:"pricepoint" bson:"pricepoint"`
	Amount         *big.Int       `json:"amount" bson:"amount"`
	MakeFee        *big.Int       `json:"makeFee" bson:"makeFee"`
	TakeFee        *big.Int       `json:"takeFee" bson:"takeFee"`
	Status         string         `json:"status" bson:"status"`
	CreatedAt      time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt" bson:"updatedAt"`
	TakerOrderSide string         `json:"takerOrderSide" bson:"takerOrderSide"`
	TakerOrderType string         `json:"takerOrderType" bson:"takerOrderType"`
	MakerOrderType string         `json:"makerOrderType" bson:"makerOrderType"`
	MakerExchange  common.Address `json:"makerExchange" bson:"makerExchange"`
	TakerExchange  common.Address `json:"takerExchange" bson:"takerExchange"`
}

// TradeSpec for query
type TradeSpec struct {
	BaseToken      string
	QuoteToken     string
	RelayerAddress common.Address
	DateFrom       int64
	DateTo         int64
}

// TradeRes response api
type TradeRes struct {
	Total  int      `json:"total" bson:"total"`
	Trades []*Trade `json:"trades" bson:"orders"`
}
type TradeRecord struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Taker          string        `json:"taker" bson:"taker"`
	Maker          string        `json:"maker" bson:"maker"`
	BaseToken      string        `json:"baseToken" bson:"baseToken"`
	QuoteToken     string        `json:"quoteToken" bson:"quoteToken"`
	MakerOrderHash string        `json:"makerOrderHash" bson:"makerOrderHash"`
	TakerOrderHash string        `json:"takerOrderHash" bson:"takerOrderHash"`
	Hash           string        `json:"hash" bson:"hash"`
	TxHash         string        `json:"txHash" bson:"txHash"`
	PairName       string        `json:"pairName" bson:"pairName"`
	Amount         string        `json:"amount" bson:"amount"`
	MakeFee        string        `json:"makeFee" bson:"makeFee"`
	TakeFee        string        `json:"takeFee" bson:"takeFee"`
	PricePoint     string        `json:"pricepoint" bson:"pricepoint"`
	Status         string        `json:"status" bson:"status"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	TakerOrderSide string        `json:"takerOrderSide" bson:"takerOrderSide"`
	TakerOrderType string        `json:"takerOrderType" bson:"takerOrderType"`
	MakerOrderType string        `json:"makerOrderType" bson:"makerOrderType"`
	MakerExchange  string        `json:"makerExchange" bson:"makerExchange"`
	TakerExchange  string        `json:"takerExchange" bson:"takerExchange"`
}

// MarshalJSON returns the json encoded byte array representing the trade struct
func (t *Trade) MarshalJSON() ([]byte, error) {
	trade := map[string]interface{}{
		"taker":          t.Taker,
		"maker":          t.Maker,
		"status":         t.Status,
		"hash":           t.Hash,
		"pairName":       t.PairName,
		"pricepoint":     t.PricePoint.String(),
		"amount":         t.Amount.String(),
		"makeFee":        t.MakeFee.String(),
		"takeFee":        t.TakeFee.String(),
		"createdAt":      t.CreatedAt.Format(time.RFC3339Nano),
		"takerOrderSide": t.TakerOrderSide,
		"takerOrderType": t.TakerOrderType,
		"makerOrderType": t.MakerOrderType,
		"makerExchange":  t.MakerExchange,
		"takerExchange":  t.TakerExchange,
	}

	if (t.BaseToken != common.Address{}) {
		trade["baseToken"] = t.BaseToken.Hex()
	}

	if (t.QuoteToken != common.Address{}) {
		trade["quoteToken"] = t.QuoteToken.Hex()
	}

	if (t.TxHash != common.Hash{}) {
		trade["txHash"] = t.TxHash.Hex()
	}

	if (t.TakerOrderHash != common.Hash{}) {
		trade["takerOrderHash"] = t.TakerOrderHash.Hex()
	}

	if (t.MakerOrderHash != common.Hash{}) {
		trade["makerOrderHash"] = t.MakerOrderHash.Hex()
	}

	return json.Marshal(trade)
}

// UnmarshalJSON creates a trade object from a json byte string
func (t *Trade) UnmarshalJSON(b []byte) error {
	trade := map[string]interface{}{}

	err := json.Unmarshal(b, &trade)
	if err != nil {
		return err
	}

	if trade["makerOrderHash"] == nil {
		return errors.New("Order Hash is not set")
	} else {
		t.MakerOrderHash = common.HexToHash(trade["makerOrderHash"].(string))
	}

	if trade["takerOrderHash"] != nil {
		t.TakerOrderHash = common.HexToHash(trade["takerOrderHash"].(string))
	}

	if trade["hash"] == nil {
		return errors.New("Hash is not set")
	} else {
		t.Hash = common.HexToHash(trade["hash"].(string))
	}

	if trade["quoteToken"] == nil {
		return errors.New("Quote token is not set")
	} else {
		t.QuoteToken = common.HexToAddress(trade["quoteToken"].(string))
	}

	if trade["baseToken"] == nil {
		return errors.New("Base token is not set")
	} else {
		t.BaseToken = common.HexToAddress(trade["baseToken"].(string))
	}

	if trade["maker"] == nil {
		return errors.New("Maker is not set")
	} else {
		t.Taker = common.HexToAddress(trade["taker"].(string))
	}

	if trade["taker"] == nil {
		return errors.New("Taker is not set")
	} else {
		t.Maker = common.HexToAddress(trade["maker"].(string))
	}

	if trade["id"] != nil && bson.IsObjectIdHex(trade["id"].(string)) {
		t.ID = bson.ObjectIdHex(trade["id"].(string))
	}

	if trade["txHash"] != nil {
		t.TxHash = common.HexToHash(trade["txHash"].(string))
	}

	if trade["pairName"] != nil {
		t.PairName = trade["pairName"].(string)
	}

	if trade["status"] != nil {
		t.Status = trade["status"].(string)
	}

	if trade["pricepoint"] != nil {
		t.PricePoint = math.ToBigInt(fmt.Sprintf("%v", trade["pricepoint"]))
	}

	if trade["amount"] != nil {
		t.Amount = new(big.Int)
		t.Amount.UnmarshalJSON([]byte(fmt.Sprintf("%v", trade["amount"])))
	}

	if trade["makeFee"] != nil {
		t.MakeFee = new(big.Int)
		t.MakeFee.UnmarshalJSON([]byte(fmt.Sprintf("%v", trade["makeFee"])))
	}
	if trade["takeFee"] != nil {
		t.TakeFee = new(big.Int)
		t.TakeFee.UnmarshalJSON([]byte(fmt.Sprintf("%v", trade["takeFee"])))
	}
	if trade["createdAt"] != nil {
		tm, _ := time.Parse(time.RFC3339Nano, trade["createdAt"].(string))
		t.CreatedAt = tm
	}
	if trade["takerOrderSide"] != nil {
		t.TakerOrderSide = trade["takerOrderSide"].(string)
	}
	if trade["takerOrderType"] != nil {
		t.TakerOrderType = trade["takerOrderType"].(string)
	}
	if trade["makerOrderType"] != nil {
		t.TakerOrderType = trade["makerOrderType"].(string)
	}
	if trade["makerExchange"] != nil {
		t.MakerExchange = common.HexToAddress(trade["makerExchange"].(string))
	}
	if trade["takerExchange"] != nil {
		t.TakerExchange = common.HexToAddress(trade["takerExchange"].(string))
	}

	return nil
}

func (t *Trade) GetBSON() (interface{}, error) {
	tr := TradeRecord{
		ID:             t.ID,
		PairName:       t.PairName,
		Maker:          t.Maker.Hex(),
		Taker:          t.Taker.Hex(),
		BaseToken:      t.BaseToken.Hex(),
		QuoteToken:     t.QuoteToken.Hex(),
		MakerOrderHash: t.MakerOrderHash.Hex(),
		Hash:           t.Hash.Hex(),
		TxHash:         t.TxHash.Hex(),
		TakerOrderHash: t.TakerOrderHash.Hex(),
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
		PricePoint:     t.PricePoint.String(),
		Status:         t.Status,
		Amount:         t.Amount.String(),
		MakeFee:        t.MakeFee.String(),
		TakeFee:        t.TakeFee.String(),
		TakerOrderSide: t.TakerOrderSide,
		TakerOrderType: t.TakerOrderType,
		MakerOrderType: t.MakerOrderType,
		MakerExchange:  t.MakerExchange.Hex(),
		TakerExchange:  t.TakerExchange.Hex(),
	}

	return tr, nil
}

func (t *Trade) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID             bson.ObjectId `json:"id,omitempty" bson:"_id"`
		PairName       string        `json:"pairName" bson:"pairName"`
		Taker          string        `json:"taker" bson:"taker"`
		Maker          string        `json:"maker" bson:"maker"`
		BaseToken      string        `json:"baseToken" bson:"baseToken"`
		QuoteToken     string        `json:"quoteToken" bson:"quoteToken"`
		MakerOrderHash string        `json:"makerOrderHash" bson:"makerOrderHash"`
		TakerOrderHash string        `json:"takerOrderHash" bson:"takerOrderHash"`
		Hash           string        `json:"hash" bson:"hash"`
		TxHash         string        `json:"txHash" bson:"txHash"`
		CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
		UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
		PricePoint     string        `json:"pricepoint" bson:"pricepoint"`
		Status         string        `json:"status" bson:"status"`
		Amount         string        `json:"amount" bson:"amount"`
		MakeFee        string        `json:"makeFee" bson:"makeFee"`
		TakeFee        string        `json:"takeFee" bson:"takeFee"`
		TakerOrderSide string        `json:"takerOrderSide" bson:"takerOrderSide"`
		TakerOrderType string        `json:"takerOrderType" bson:"takerOrderType"`
		MakerOrderType string        `json:"makerOrderType" bson:"makerOrderType"`
		MakerExchange  string        `json:"makerExchange" bson:"makerExchange"`
		TakerExchange  string        `json:"takerExchange" bson:"takerExchange"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	t.ID = decoded.ID
	t.PairName = decoded.PairName
	t.Taker = common.HexToAddress(decoded.Taker)
	t.Maker = common.HexToAddress(decoded.Maker)
	t.BaseToken = common.HexToAddress(decoded.BaseToken)
	t.QuoteToken = common.HexToAddress(decoded.QuoteToken)
	t.MakerOrderHash = common.HexToHash(decoded.MakerOrderHash)
	t.TakerOrderHash = common.HexToHash(decoded.TakerOrderHash)
	t.Hash = common.HexToHash(decoded.Hash)
	t.TxHash = common.HexToHash(decoded.TxHash)
	t.Status = decoded.Status
	t.Amount = math.ToBigInt(decoded.Amount)
	t.PricePoint = math.ToBigInt(decoded.PricePoint)

	t.MakeFee = math.ToBigInt(decoded.MakeFee)
	t.TakeFee = math.ToBigInt(decoded.TakeFee)

	t.CreatedAt = decoded.CreatedAt
	t.UpdatedAt = decoded.UpdatedAt
	t.TakerOrderSide = decoded.TakerOrderSide
	t.TakerOrderType = decoded.TakerOrderType
	t.MakerOrderType = decoded.TakerOrderType
	t.MakerExchange = common.HexToAddress(decoded.MakerExchange)
	t.TakerExchange = common.HexToAddress(decoded.TakerExchange)
	return nil
}

// ComputeHash returns hashes the trade
// The OrderHash, Amount, Taker and TradeNonce attributes must be
// set before attempting to compute the trade hash
func (t *Trade) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(t.MakerOrderHash.Bytes())
	sha.Write(t.TakerOrderHash.Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

func (t *Trade) Pair() (*Pair, error) {
	if (t.BaseToken == common.Address{}) {
		return nil, errors.New("Base token is not set")
	}

	if (t.QuoteToken == common.Address{}) {
		return nil, errors.New("Quote token is set")
	}

	return &Pair{
		BaseTokenAddress:  t.BaseToken,
		QuoteTokenAddress: t.QuoteToken,
	}, nil
}

type TradeBSONUpdate struct {
	*Trade
}

func (t TradeBSONUpdate) GetBSON() (interface{}, error) {
	now := time.Now()

	set := bson.M{
		"taker":          t.Taker.Hex(),
		"maker":          t.Maker.Hex(),
		"baseToken":      t.BaseToken.Hex(),
		"quoteToken":     t.QuoteToken.Hex(),
		"makerOrderHash": t.MakerOrderHash.Hex(),
		"takerOrderHash": t.TakerOrderHash.Hex(),
		"txHash":         t.TxHash.Hex(),
		"pairName":       t.PairName,
		"status":         t.Status,
	}

	if t.PricePoint != nil {
		set["pricepoint"] = t.PricePoint.String()
	}

	if t.Amount != nil {
		set["amount"] = t.Amount.String()
	}

	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"hash":      t.Hash.Hex(),
		"createdAt": now,
	}

	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}

type updateDesc struct {
	UpdatedFields map[string]interface{} `bson:"updatedFields"`
	RemovedFields []string               `bson:"removedFields"`
}
type M bson.M
type evNamespace struct {
	DB   string `bson:"db"`
	Coll string `bson:"coll"`
}
type TradeChangeEvent struct {
	ID                interface{} `bson:"_id"`
	OperationType     string      `bson:"operationType"`
	FullDocument      *Trade      `bson:"fullDocument,omitempty"`
	Ns                evNamespace `bson:"ns"`
	DocumentKey       M           `bson:"documentKey"`
	UpdateDescription *updateDesc `bson:"updateDescription,omitempty"`
}

// UserTrade trade info of user
type UserTrade struct {
	UserAddress      common.Address `json:"userAddress"`
	Count            *big.Int       `json:"count"`
	Volume           *big.Int       `json:"volume"`
	VolumeByQuote    *big.Int       `json:"volumeByQuote"`
	VolumeAskByQuote *big.Int       `json:"volumeAskByQuote"`
	VolumeBidByQuote *big.Int       `json:"volumeBidByQuote"`

	VolumeAsk      *big.Int       `json:"volumeAsk"`
	VolumeBid      *big.Int       `json:"volumeBid"`
	TimeStamp      int64          `json:"timestamp"`
	RelayerAddress common.Address `json:"relayerAddress"`
	BaseToken      common.Address `json:"baseToken"`
	QuoteToken     common.Address `json:"quoteToken"`
}

// RelayerTrade relayer trade
type RelayerTrade struct {
	RelayerAddress common.Hash `json:"relayerAddress"`
	Count          *big.Int    `json:"count"`
	Volume         *big.Int    `json:"volume"`
}

// UserTradeSpec user trade filter
type UserTradeSpec struct {
}

// UserVolume user volume trade
type UserVolume struct {
	UserAddress common.Address `json:"userAddress"`
	Volume      *big.Int       `json:"volume"`
	Rank        int            `json:"rank"`
}

// TradeVolume trade volume info
type TradeVolume struct {
	Trader      *big.Int `json:"trader"`
	TotalVolume *big.Int `json:"totalVolume"`
}

// UserPnL user volume trade
type UserPnL struct {
	UserAddress      common.Address `json:"userAddress"`
	VolumeAskByQuote *big.Int       `json:"volumeAskByQuote"`
	VolumeBidByQuote *big.Int       `json:"volumeBidByQuote"`
	VolumeAsk        *big.Int       `json:"volumeAsk"`
	VolumeBid        *big.Int       `json:"volumeBid"`
	CurrentPrice     *big.Int       `json:"currentPrice"`
	PnL              *big.Int       `json:"currentPnL"`
}
