package services

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/tomochain/viewdex/daos"
	"github.com/tomochain/viewdex/relayer"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/viewdex/app"
	"github.com/tomochain/viewdex/types"
)

// RelayerService struct
type RelayerService struct {
	relayer    *relayer.Relayer
	tokenDao   *daos.TokenDao
	pairDao    *daos.PairDao
	relayerDao *daos.RelayerDao
}

// NewRelayerService returns a new instance of orderservice
func NewRelayerService(
	relaye *relayer.Relayer,
	tokenDao *daos.TokenDao,
	pairDao *daos.PairDao,
	relayerDao *daos.RelayerDao,
) *RelayerService {
	return &RelayerService{
		relaye,
		tokenDao,
		pairDao,
		relayerDao,
	}
}

func (s *RelayerService) GetByAddress(addr common.Address) (*types.Relayer, error) {
	return s.relayerDao.GetByAddress(addr)
}

func (s *RelayerService) UpdateNameByAddress(addr common.Address, name string, url string) error {
	return s.relayerDao.UpdateNameByAddress(addr, name, url)
}

func (s *RelayerService) GetRelayerAddress(r *http.Request) common.Address {
	v := r.URL.Query()
	relayerAddress := v.Get("relayerAddress")

	if relayerAddress == "" {
		relayer, _ := s.relayerDao.GetByHost(r.Host)
		if relayer != nil {
			relayerAddress = relayer.Address.Hex()
		}
	}

	if relayerAddress == "" {
		relayerAddress = app.Config.Tomochain["exchange_address"]
	}

	return common.HexToAddress(relayerAddress)
}

func (s *RelayerService) updatePairRelayer(relayerInfo *relayer.RInfo) error {
	currentPairs, err := s.pairDao.GetAllByCoinbase(relayerInfo.Address)
	fmt.Println("UpdatePairRelayer starting...", relayerInfo.Address.Hex())
	if err != nil {
		return err
	}

	for _, newpair := range relayerInfo.Pairs {
		found := false
		for _, currentPair := range currentPairs {
			if newpair.BaseToken == currentPair.BaseTokenAddress && newpair.QuoteToken == currentPair.QuoteTokenAddress {
				found = true
			}
		}
		if !found {
			pairBaseData := relayerInfo.Tokens[newpair.BaseToken]
			pairQuoteData := relayerInfo.Tokens[newpair.QuoteToken]
			pair := &types.Pair{
				BaseTokenSymbol:    pairBaseData.Symbol,
				BaseTokenAddress:   newpair.BaseToken,
				BaseTokenDecimals:  int(pairBaseData.Decimals),
				QuoteTokenSymbol:   pairQuoteData.Symbol,
				QuoteTokenAddress:  newpair.QuoteToken,
				QuoteTokenDecimals: int(pairQuoteData.Decimals),
				RelayerAddress:     relayerInfo.Address,
				Active:             true,
				MakeFee:            big.NewInt(int64(relayerInfo.MakeFee)),
				TakeFee:            big.NewInt(int64(relayerInfo.TakeFee)),
			}
			fmt.Println("Create Pair:", pair.BaseTokenAddress.Hex(), pair.QuoteTokenAddress.Hex(), relayerInfo.Address.Hex())
			err := s.pairDao.Create(pair)
			if err != nil {

			}
		}
	}

	for _, currentPair := range currentPairs {
		found := false
		for _, newpair := range relayerInfo.Pairs {
			if currentPair.BaseTokenAddress == newpair.BaseToken && currentPair.QuoteTokenAddress == newpair.QuoteToken {
				found = true
			}
		}
		if !found {
			fmt.Println("Delete Pair:", currentPair.BaseTokenAddress.Hex(), currentPair.QuoteTokenAddress.Hex())
			err := s.pairDao.DeleteByTokenAndCoinbase(currentPair.BaseTokenAddress, currentPair.QuoteTokenAddress, relayerInfo.Address)
			if err != nil {

			}
		}
	}
	return nil
}

func (s *RelayerService) updateRelayers(relayerInfos []*relayer.RInfo, lendingRelayerInfos []*relayer.LendingRInfo) error {
	currentRelayers, err := s.relayerDao.GetAll()
	if err != nil {
		return err
	}

	found := false
	for _, r := range relayerInfos {
		found = false
		for _, v := range currentRelayers {
			if v.Address.Hex() == r.Address.Hex() {
				found = true
				break
			}
		}
		lendingFee := uint16(0)
		for _, l := range lendingRelayerInfos {
			if l.Address.Hex() == r.Address.Hex() {
				lendingFee = l.Fee
				break
			}
		}
		relayer := &types.Relayer{
			RID:        r.RID,
			Owner:      r.Owner,
			Deposit:    r.Deposit,
			Address:    r.Address,
			Resign:     r.Resign,
			LockTime:   r.LockTime,
			MakeFee:    big.NewInt(int64(r.MakeFee)),
			TakeFee:    big.NewInt(int64(r.TakeFee)),
			LendingFee: big.NewInt(int64(lendingFee)),
		}
		if !found {
			fmt.Println("Create relayer:", r.Address.Hex())
			err = s.relayerDao.Create(relayer)
			if err != nil {

			}
		} else {
			fmt.Println("Update relayer:", r.Address.Hex())
			err = s.relayerDao.UpdateByAddress(r.Address, relayer)
			if err != nil {

			}
		}
	}

	for _, r := range currentRelayers {
		found = false
		for _, v := range relayerInfos {
			if v.Address.Hex() == r.Address.Hex() {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Delete relayer:", r.Address.Hex)
			err = s.relayerDao.DeleteByAddress(r.Address)
			if err != nil {

			}
		}
	}
	return nil
}

func (s *RelayerService) updateRelayer(relayerInfo *relayer.RInfo, lendingRelayerInfo *relayer.LendingRInfo) error {
	currentRelayer, err := s.relayerDao.GetByAddress(relayerInfo.Address)
	if err != nil {
		return err
	}

	found := false
	if currentRelayer != nil {
		found = true
	}

	lendingFee := lendingRelayerInfo.Fee
	relayer := &types.Relayer{
		RID:        relayerInfo.RID,
		Owner:      relayerInfo.Owner,
		Deposit:    relayerInfo.Deposit,
		Address:    relayerInfo.Address,
		Resign:     relayerInfo.Resign,
		LockTime:   relayerInfo.LockTime,
		MakeFee:    big.NewInt(int64(relayerInfo.MakeFee)),
		TakeFee:    big.NewInt(int64(relayerInfo.TakeFee)),
		LendingFee: big.NewInt(int64(lendingFee)),
	}

	if !found {
		fmt.Println("Create relayer:", relayerInfo.Address.Hex())
		err = s.relayerDao.Create(relayer)
		if err != nil {

		}
	} else {
		fmt.Println("Update relayer:", relayerInfo.Address.Hex())
		err = s.relayerDao.UpdateByAddress(relayerInfo.Address, relayer)
		if err != nil {

		}
	}

	return nil
}

func (s *RelayerService) updateTokenRelayer(relayerInfo *relayer.RInfo) error {
	currentTokens, err := s.tokenDao.GetAllByCoinbase(relayerInfo.Address)
	if err != nil {
		return err
	}

	for ntoken, v := range relayerInfo.Tokens {
		found := false
		for _, ctoken := range currentTokens {
			if ntoken.Hex() == ctoken.ContractAddress.Hex() {
				found = true
			}
		}
		token := &types.Token{
			Symbol:          v.Symbol,
			ContractAddress: ntoken,
			RelayerAddress:  relayerInfo.Address,
			Decimals:        int(v.Decimals),
			MakeFee:         big.NewInt(int64(relayerInfo.MakeFee)),
			TakeFee:         big.NewInt(int64(relayerInfo.TakeFee)),
		}
		if !found {
			fmt.Println("Create Token:", token.ContractAddress.Hex())
			err = s.tokenDao.Create(token)
			if err != nil {

			}
		} else {
			fmt.Println("Update Token:", token.ContractAddress.Hex())
			err = s.tokenDao.UpdateByTokenAndCoinbase(ntoken, relayerInfo.Address, token)
		}
		for _, ctoken := range currentTokens {
			found = false
			for ntoken, v = range relayerInfo.Tokens {

				if ctoken.ContractAddress.Hex() == ntoken.Hex() {
					found = true
				}
			}
			if !found {
				fmt.Println("Delete Token:", ctoken.ContractAddress.Hex)
				err = s.tokenDao.DeleteByTokenAndCoinbase(ctoken.ContractAddress, relayerInfo.Address)
				if err != nil {

				}
			}
		}
	}
	return nil
}

// UpdateRelayer get the total number of orders amount created by a user
func (s *RelayerService) UpdateRelayer(coinbase common.Address) error {
	relayerInfo, err := s.relayer.GetRelayer(coinbase)
	if err != nil {
		return err
	}
	s.updateTokenRelayer(relayerInfo)
	s.updatePairRelayer(relayerInfo)

	relayerLendingInfo, err := s.relayer.GetLending()
	if err != nil {
		return err
	}
	s.updateRelayer(relayerInfo, relayerLendingInfo)
	return nil
}

func (s *RelayerService) UpdateRelayers() error {
	relayerInfos, err := s.relayer.GetRelayers()
	if err != nil {
		return err
	}
	for _, relayerInfo := range relayerInfos {
		s.updateTokenRelayer(relayerInfo)
		s.updatePairRelayer(relayerInfo)
	}

	relayerLendingInfos, err := s.relayer.GetLendings()
	if err != nil {
		return err
	}

	s.updateRelayers(relayerInfos, relayerLendingInfos)
	return nil
}
