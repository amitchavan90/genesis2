package blockchain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"genesis/config"
	"genesis/solidity"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/ninja-software/terror"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BinsPerBatch for blockchain - the amount of transactions that can fit in a single batch
const BinsPerBatch int = 16

// ErrBlockchainConnectionIssue for blockchain
var ErrBlockchainConnectionIssue = errors.New("unable to connect to smart contact")

// ErrBlockchainOutOfGas for blockchain
var ErrBlockchainOutOfGas = errors.New("smart contract out of gas")

// ErrSmartContractAlreadyDeployed for blockchain
var ErrSmartContractAlreadyDeployed = errors.New("contract already deployed")

// Service struct for Blockchain
type Service struct {
	Client            *ethclient.Client
	TransactionReader ethereum.TransactionReader
	ChainReader       ethereum.ChainReader
	Chain             bind.ContractBackend
	Auth              *bind.TransactOpts
	AuthCall          *bind.CallOpts
	Mux               sync.Mutex // blockchain must be operating single threaded, in correct series
	ctx               context.Context

	Mailer *mailgun.MailgunImpl
	Config *config.PlatformConfig
}

// MakeAddress for Blockchain
func MakeAddress() (common.Address, *ecdsa.PrivateKey, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return [20]byte{}, nil, terror.New(err, "")
	}
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return [20]byte{}, nil, terror.New(fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey"), "")
	}
	addr := crypto.FromECDSAPub(publicKeyECDSA)
	result := [20]byte{}
	copy(result[:], addr)
	return result, key, nil
}

// New Blockchain service
func New(
	ctx context.Context,
	host string,
	key *ecdsa.PrivateKey,
	sim bool,
	mailer *mailgun.MailgunImpl,
	config *config.PlatformConfig,
) (*Service, error) {
	if key == nil {
		return nil, terror.New(fmt.Errorf("blockchain private key is nil"), "")
	}

	var chain bind.ContractBackend
	var err error
	var auth *bind.TransactOpts
	if sim {
		fmt.Println("WARNING: RUNNING SIMULATION")
		ticker := time.NewTicker(1 * time.Second)
		key, _ := crypto.GenerateKey()
		auth = bind.NewKeyedTransactor(key)
		address := auth.From
		authCall := &bind.CallOpts{
			From:    auth.From,
			Pending: true,
			Context: ctx,
		}
		gAlloc := map[common.Address]core.GenesisAccount{
			address: {Balance: big.NewInt(10000000000)},
		}

		simchain := backends.NewSimulatedBackend(gAlloc, 0)
		go func(chain *backends.SimulatedBackend) {
			for {
				select {
				case <-ticker.C:
					chain.Commit()
				case <-ctx.Done():
					fmt.Println("cancelled")
					return
				}
			}
		}(simchain)
		chain = simchain
		return &Service{
			Client:            nil,
			TransactionReader: nil,
			ChainReader:       nil,
			Chain:             chain,
			Auth:              auth,
			AuthCall:          authCall,
			ctx:               ctx,
			Mailer:            mailer,
			Config:            config,
		}, nil
	}
	auth = bind.NewKeyedTransactor(key)
	authCall := &bind.CallOpts{
		From:    auth.From,
		Pending: true,
		Context: ctx,
	}
	client, err := ethclient.DialContext(ctx, host)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return &Service{
		Client:            client,
		TransactionReader: client,
		ChainReader:       client,
		Chain:             client,
		Auth:              auth,
		AuthCall:          authCall,
		ctx:               ctx,
		Mailer:            mailer,
		Config:            config,
	}, nil
}

// GetBalance returns to current account's balance
func (s *Service) GetBalance(ctx context.Context) (*big.Int, error) {
	result, err := s.Client.BalanceAt(ctx, s.Auth.From, nil)
	if err != nil {
		return nil, terror.New(err, "")
	}
	return result, nil
}

// CheckBalance checks the current account balance and send a notification email if below a certain amount
func (s *Service) CheckBalance(ctx context.Context) error {
	if s.Config.Blockchain.EthLowNotifyEmail == "" {
		return nil
	}

	// Check balance
	balance, err := s.GetBalance(ctx)
	if err != nil {
		return terror.New(err, "failed to get eth balance")
	}

	gwei := new(big.Int).Div(balance, big.NewInt(10).Exp(big.NewInt(10), big.NewInt(9), nil))
	eth := float64(gwei.Int64()) / math.Pow10(9)

	if eth > s.Config.Blockchain.EthLowNotifyAmount {
		return nil
	}

	// Create email
	sender := s.Config.Email.Sender
	subject := "Genesis - ETH is Low"

	message := s.Mailer.NewMessage(sender, subject, "", strings.Split(s.Config.Blockchain.EthLowNotifyEmail, ",")...)
	message.SetTemplate("portal_basic")
	message.AddVariable("message", fmt.Sprintf("Your account balance is currently: %f ETH", eth))

	// Send Email
	emailCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err = s.Mailer.Send(emailCtx, message)
	if err != nil {
		return terror.New(err, "failed to send email")
	}

	return nil
}

// DeploySmartContract creates a blockchain smart contract
func (s *Service) DeploySmartContract(ctx context.Context, name string) (common.Address, *types.Transaction, *solidity.SmartContract, error) {
	err := s.UpdateGas(ctx)
	if err != nil {
		return common.Address{}, nil, nil, terror.New(err, "")
	}

	var _name [32]byte
	copy(_name[:], []byte(name))
	address, tx, contract, err := solidity.DeploySmartContract(s.Auth, s.Chain, _name)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "sender doesn't have enough funds") {
			return address, tx, contract, terror.New(ErrBlockchainOutOfGas, "")
		} else if strings.Contains(msg, "No connection could be made") {
			return address, tx, contract, terror.New(ErrBlockchainConnectionIssue, "")
		}
		return address, tx, contract, terror.New(err, "deploy smart contract")
	}

	return address, tx, contract, nil
}

// UpdateGas calculates gas and updates the service's Auth
func (s *Service) UpdateGas(ctx context.Context) error {
	nonce, gasPrice, err := s.CalculateGas(ctx)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "connect: connection refused") {
			return terror.New(ErrBlockchainConnectionIssue, "")
		}
		return terror.New(err, "smart contract: calculate gas")
	}

	s.Auth.GasPrice = gasPrice
	s.Auth.Nonce = nonce

	return nil
}

// CalculateGas for the contract
func (s *Service) CalculateGas(ctx context.Context) (*big.Int, *big.Int, error) {
	nonce, err := s.Chain.PendingNonceAt(ctx, s.Auth.From)
	if err != nil {
		return nil, nil, terror.New(err, "")
	}

	gasPrice, err := s.Chain.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, terror.New(err, "")
	}

	return big.NewInt(int64(nonce)), gasPrice, nil
}

// GetContract for the product
func (s *Service) GetContract(addr common.Address) (*solidity.SmartContract, error) {
	return solidity.NewSmartContract(addr, s.Chain)
}
