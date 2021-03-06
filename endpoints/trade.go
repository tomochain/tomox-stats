package endpoints

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-stats/services"
	"github.com/tomochain/tomox-stats/types"
	"github.com/tomochain/tomox-stats/utils/httputils"
)

type tradeEndpoint struct {
	tradeService *services.TradeService
}

// ServeTradeResource sets up the routing of trade endpoints and the corresponding handlers.
// TODO trim down to one single endpoint with the 3 following params: base, quote, address
func ServeTradeResource(
	r *mux.Router,
	tradeService *services.TradeService,
) {
	e := &tradeEndpoint{tradeService}
	r.HandleFunc("/stats/trades/volume", e.handleQueryVolume)
	r.HandleFunc("/stats/trades/total", e.handleQueryVolume)
	r.HandleFunc("/stats/trades/volume24h", e.handleQuery24h)
	r.HandleFunc("/stats/trades/top/pnl", e.handleGetRelayerTopPnLTrades)
	r.HandleFunc("/stats/trades/users/count", e.handleGetNumberUser)
}

// HandleGetTrades is responsible for getting pair's trade history requests
func (e *tradeEndpoint) handleQueryVolume(w http.ResponseWriter, r *http.Request) {

	var baseTokens []common.Address
	var quoteToken common.Address
	var relayerAddress common.Address
	var userAddress common.Address

	var from int64
	var to int64
	var topVolume int
	v := r.URL.Query()
	qt := v.Get("quoteToken")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	rAddress := v.Get("relayerAddress")
	uaddr := v.Get("userAddress")
	top := v.Get("top")
	for _, bt := range v["baseToken"] {
		if bt != "" {
			if common.IsHexAddress(bt) {
				baseTokens = append(baseTokens, common.HexToAddress(bt))
			}
		}
	}
	topVolume = 10

	if qt != "" {
		if !common.IsHexAddress(qt) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid quotetoken address")
			return
		}
		quoteToken = common.HexToAddress(qt)

	}

	if rAddress != "" {
		if !common.IsHexAddress(rAddress) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid relayer address")
			return
		}
		relayerAddress = common.HexToAddress(rAddress)
	}

	if toParam != "" {
		t, _ := strconv.Atoi(toParam)
		to = int64(t)
	}
	if fromParam != "" {
		t, _ := strconv.Atoi(fromParam)
		from = int64(t)
	}
	if top != "" {
		topVolume, _ = strconv.Atoi(top)
	}

	if uaddr != "" {
		if common.IsHexAddress(uaddr) {
			userAddress = common.HexToAddress(uaddr)
		}
	}

	res := e.tradeService.QueryVolume(relayerAddress, userAddress, baseTokens, quoteToken, from, to, topVolume)

	if res == nil {

		httputils.WriteJSON(w, http.StatusOK, []*types.UserVolume{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

// HandleGetTrades is responsible for getting pair's trade history requests
func (e *tradeEndpoint) handleGetRelayerTopPnLTrades(w http.ResponseWriter, r *http.Request) {

	var baseToken common.Address
	var quoteToken common.Address
	var relayerAddress common.Address
	var topVolume int
	v := r.URL.Query()
	bt := v.Get("baseToken")
	qt := v.Get("quoteToken")
	rAddress := v.Get("relayerAddress")
	top := v.Get("top")

	topVolume = 10

	if bt != "" {
		if !common.IsHexAddress(bt) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid basetoken address")
			return
		}
		baseToken = common.HexToAddress(bt)
	}

	if qt != "" {
		if !common.IsHexAddress(qt) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid quotetoken address")
			return
		}
		quoteToken = common.HexToAddress(qt)

	}

	if rAddress != "" {
		if !common.IsHexAddress(rAddress) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid relayer address")
			return
		}
		relayerAddress = common.HexToAddress(rAddress)
	}

	if top != "" {
		topVolume, _ = strconv.Atoi(top)
	}

	res := e.tradeService.GetTopRelayerUserPnL(relayerAddress, baseToken, quoteToken, topVolume)

	if res == nil {

		httputils.WriteJSON(w, http.StatusOK, []*types.UserPnL{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tradeEndpoint) handleQuery24h(w http.ResponseWriter, r *http.Request) {

	var baseTokens []common.Address
	var quoteToken common.Address
	var relayerAddress common.Address
	var userAddress common.Address

	var topVolume int
	v := r.URL.Query()
	qt := v.Get("quoteToken")
	rAddress := v.Get("relayerAddress")
	uaddr := v.Get("userAddress")
	top := v.Get("top")
	for _, bt := range v["baseToken"] {
		if bt != "" {
			if common.IsHexAddress(bt) {
				baseTokens = append(baseTokens, common.HexToAddress(bt))
			}
		}
	}
	topVolume = 10

	if qt != "" {
		if !common.IsHexAddress(qt) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid quotetoken address")
			return
		}
		quoteToken = common.HexToAddress(qt)

	}

	if rAddress != "" {
		if !common.IsHexAddress(rAddress) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid relayer address")
			return
		}
		relayerAddress = common.HexToAddress(rAddress)
	}

	if top != "" {
		topVolume, _ = strconv.Atoi(top)
	}

	if uaddr != "" {
		if common.IsHexAddress(uaddr) {
			userAddress = common.HexToAddress(uaddr)
		}
	}

	res := e.tradeService.Query24hVolume(relayerAddress, userAddress, baseTokens, quoteToken, topVolume)

	if res == nil {

		httputils.WriteJSON(w, http.StatusOK, []*types.UserVolume{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tradeEndpoint) handleQueryTotal(w http.ResponseWriter, r *http.Request) {

	var baseTokens []common.Address
	var quoteToken common.Address
	var relayerAddress common.Address

	var from int64
	var to int64
	v := r.URL.Query()
	qt := v.Get("quoteToken")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	rAddress := v.Get("relayerAddress")
	for _, bt := range v["baseToken"] {
		if bt != "" {
			if common.IsHexAddress(bt) {
				baseTokens = append(baseTokens, common.HexToAddress(bt))
			}
		}
	}

	if qt != "" {
		if !common.IsHexAddress(qt) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid quotetoken address")
			return
		}
		quoteToken = common.HexToAddress(qt)

	}

	if rAddress != "" {
		if !common.IsHexAddress(rAddress) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid relayer address")
			return
		}
		relayerAddress = common.HexToAddress(rAddress)
	}

	if toParam != "" {
		t, _ := strconv.Atoi(toParam)
		to = int64(t)
	}
	if fromParam != "" {
		t, _ := strconv.Atoi(fromParam)
		from = int64(t)
	}

	res := e.tradeService.QueryTotal(relayerAddress, baseTokens, quoteToken, from, to)

	if res == nil {

		httputils.WriteJSON(w, http.StatusOK, []*types.UserVolume{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *tradeEndpoint) handleGetNumberUser(w http.ResponseWriter, r *http.Request) {

	var relayerAddress common.Address
	var baseToken common.Address
	var quoteToken common.Address
	v := r.URL.Query()
	rAddress := v.Get("relayerAddress")
	duration := v.Get("duration")
	bot := v.Get("excludeBot")
	bToken := v.Get("baseToken")
	qToken := v.Get("quoteToken")
	type NumberTrader struct {
		ActiveUser int    `json:"activeUser"`
		Duration   string `json:"duration"`
	}
	var res NumberTrader
	excludeBot := false

	if bot == "true" {
		excludeBot = true
	}
	if rAddress != "" {
		if !common.IsHexAddress(rAddress) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid relayer address")
			return
		}
		relayerAddress = common.HexToAddress(rAddress)
	}
	if bToken != "" {
		if !common.IsHexAddress(bToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid baseToken address")
			return
		}
		baseToken = common.HexToAddress(bToken)
	}
	if qToken != "" {
		if !common.IsHexAddress(qToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid quoteToken address")
			return
		}
		quoteToken = common.HexToAddress(qToken)
	}
	if duration == "" {
		res.Duration = "all"
		res.ActiveUser = e.tradeService.GetNumberTraderByTime(relayerAddress, baseToken, quoteToken, 0, 0, excludeBot)
		httputils.WriteJSON(w, http.StatusOK, res)
		return
	}
	if duration == "1d" {
		res.Duration = duration
		res.ActiveUser = e.tradeService.GetNumberTraderByTime(relayerAddress, baseToken, quoteToken, time.Now().AddDate(0, 0, -1).Unix(), 0, excludeBot)
		httputils.WriteJSON(w, http.StatusOK, res)
		return
	}
	if duration == "7d" {
		res.Duration = duration
		res.ActiveUser = e.tradeService.GetNumberTraderByTime(relayerAddress, baseToken, quoteToken, time.Now().AddDate(0, 0, -7).Unix(), 0, excludeBot)
		httputils.WriteJSON(w, http.StatusOK, res)
		return
	}
	if duration == "30d" {
		res.Duration = duration
		res.ActiveUser = e.tradeService.GetNumberTraderByTime(relayerAddress, baseToken, quoteToken, time.Now().AddDate(0, 0, -30).Unix(), 0, excludeBot)
		httputils.WriteJSON(w, http.StatusOK, res)
		return
	}
	httputils.WriteJSON(w, http.StatusBadRequest, "duration must be empty/1d/7d/30d")
}
