package chainalysis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"go.uber.org/zap"

	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	timeout = time.Minute * 5

	ethSymbol  = "ETH"
	ethAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

// Client is implementation of chainalysis client
type Client struct {
	host   string
	apiKey string
	sugar  *zap.SugaredLogger
	client *http.Client
}

// NewChainAlysisClient creates a new chainalysis client instance.
func NewChainAlysisClient(sugar *zap.SugaredLogger, host, apiKey string) *Client {
	c := &Client{
		host:   host,
		apiKey: apiKey,
		sugar:  sugar,
		client: &http.Client{Timeout: timeout},
	}
	return c
}

type registerData struct {
	RwData  []registerWithdrawal
	RstData []registerSentTransfer
}

type registerWithdrawal struct {
	Asset   string
	Address ethereum.Address
}

type registerSentTransfer struct {
	Asset             string
	TransferReference ethereum.Hash
}

func updateRegisterData(rd registerData, asset string, txHash ethereum.Hash, receiveAdderss ethereum.Address) registerData {
	rd.RwData = append(rd.RwData, registerWithdrawal{
		Asset:   asset,
		Address: receiveAdderss,
	})
	rd.RstData = append(rd.RstData, registerSentTransfer{
		Asset:             asset,
		TransferReference: txHash,
	})
	return rd
}

// PushETHSentTransferEvent push eth sent transfer to chainalysis api
func (c *Client) PushETHSentTransferEvent(tradeLogs []common.TradeLog) error {
	mapRegisterData := make(map[ethereum.Address]registerData)
	for _, log := range tradeLogs {
		var (
			txHash         = log.TransactionHash
			userAddress    = log.UserAddress
			receiveAdderss = log.ReceiveAddress
		)
		if strings.ToLower(log.DestAddress.Hex()) != ethAddress {
			continue
		}

		c.sugar.Debugw("sent transfer data",
			"user addr", userAddress,
			"receive addr", receiveAdderss,
			"tx hash", txHash)
		if rd, ok := mapRegisterData[userAddress]; ok {
			mapRegisterData[userAddress] = updateRegisterData(rd, ethSymbol, txHash, receiveAdderss)
		} else {
			mapRegisterData[userAddress] = registerData{
				RwData: []registerWithdrawal{
					{
						ethSymbol,
						receiveAdderss,
					},
				},
				RstData: []registerSentTransfer{
					{
						ethSymbol,
						txHash,
					},
				},
			}
		}
	}
	for userAddress, registerData := range mapRegisterData {
		if err := c.registerWithdrawalAddress(userAddress, registerData.RwData); err != nil {
			c.sugar.Errorw("got error when register withdrawal address",
				"error", err.Error(),
				"user address", userAddress,
				"register withdrawal data", registerData.RwData)
			return err
		}
		if err := c.registerSentTransfer(userAddress, registerData.RstData); err != nil {
			c.sugar.Errorw("got error when register sent transfer",
				"error", err.Error(),
				"user address", userAddress,
				"register sent transfer data", registerData.RstData)
			return err
		}
	}
	return nil
}

// registerWithdrawalAddress register withdrawal address
func (c *Client) registerWithdrawalAddress(userAddr ethereum.Address, rw []registerWithdrawal) error {
	url := fmt.Sprintf("%s/users/%s/withdrawaladdresses", c.host, userAddr.Hex())
	body, err := json.Marshal(rw)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Token", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			c.sugar.Errorw("failed to close body", "err", cErr.Error())
		}
	}()
	return nil
}

// registerSentTransfer register sent transfer
func (c *Client) registerSentTransfer(userAddr ethereum.Address, rst []registerSentTransfer) error {
	url := fmt.Sprintf("%s/users/%s/transfers/sent", c.host, userAddr.Hex())
	body, err := json.Marshal(rst)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Token", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			c.sugar.Errorw("failed to close body", "err", cErr.Error())
		}
	}()
	return nil
}