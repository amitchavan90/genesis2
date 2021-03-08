// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package solidity

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SmartContractABI is the input ABI used to generate the binding from.
const SmartContractABI = "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_name\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"batchCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batch\",\"type\":\"uint256\"}],\"name\":\"commitBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// SmartContractFuncSigs maps the 4-byte function signature to its string representation.
var SmartContractFuncSigs = map[string]string{
	"06f13056": "batchCount()",
	"c9bacb66": "commitBatch(uint256)",
	"8da5cb5b": "owner()",
}

// SmartContractBin is the compiled bytecode used for deploying new contracts.
var SmartContractBin = "0x608060405234801561001057600080fd5b5060405161022e38038061022e8339818101604052602081101561003357600080fd5b5051600080546001600160a01b031916331790556001556101d5806100596000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806306f13056146100465780638da5cb5b14610060578063c9bacb6614610084575b600080fd5b61004e6100a3565b60408051918252519081900360200190f35b610068610103565b604080516001600160a01b039092168252519081900360200190f35b6100a16004803603602081101561009a57600080fd5b5035610112565b005b600080546001600160a01b031633146100fc576040805162461bcd60e51b815260206004820152601660248201527529b2b73232b9103737ba1030baba3437b934bd32b21760511b604482015290519081900360640190fd5b5060025490565b6000546001600160a01b031690565b6000546001600160a01b0316331461016a576040805162461bcd60e51b815260206004820152601660248201527529b2b73232b9103737ba1030baba3437b934bd32b21760511b604482015290519081900360640190fd5b600280546001810182556000919091527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace015556fea2646970667358221220a2d459081b245e81f01104b77509c07d5965b11b988c0a31c740ef0d119913f664736f6c634300060c0033"

// DeploySmartContract deploys a new Ethereum contract, binding an instance of SmartContract to it.
func DeploySmartContract(auth *bind.TransactOpts, backend bind.ContractBackend, _name [32]byte) (common.Address, *types.Transaction, *SmartContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SmartContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SmartContractBin), backend, _name)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SmartContract{SmartContractCaller: SmartContractCaller{contract: contract}, SmartContractTransactor: SmartContractTransactor{contract: contract}, SmartContractFilterer: SmartContractFilterer{contract: contract}}, nil
}

// SmartContract is an auto generated Go binding around an Ethereum contract.
type SmartContract struct {
	SmartContractCaller     // Read-only binding to the contract
	SmartContractTransactor // Write-only binding to the contract
	SmartContractFilterer   // Log filterer for contract events
}

// SmartContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type SmartContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmartContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SmartContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmartContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SmartContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmartContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SmartContractSession struct {
	Contract     *SmartContract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SmartContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SmartContractCallerSession struct {
	Contract *SmartContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SmartContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SmartContractTransactorSession struct {
	Contract     *SmartContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SmartContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type SmartContractRaw struct {
	Contract *SmartContract // Generic contract binding to access the raw methods on
}

// SmartContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SmartContractCallerRaw struct {
	Contract *SmartContractCaller // Generic read-only contract binding to access the raw methods on
}

// SmartContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SmartContractTransactorRaw struct {
	Contract *SmartContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSmartContract creates a new instance of SmartContract, bound to a specific deployed contract.
func NewSmartContract(address common.Address, backend bind.ContractBackend) (*SmartContract, error) {
	contract, err := bindSmartContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SmartContract{SmartContractCaller: SmartContractCaller{contract: contract}, SmartContractTransactor: SmartContractTransactor{contract: contract}, SmartContractFilterer: SmartContractFilterer{contract: contract}}, nil
}

// NewSmartContractCaller creates a new read-only instance of SmartContract, bound to a specific deployed contract.
func NewSmartContractCaller(address common.Address, caller bind.ContractCaller) (*SmartContractCaller, error) {
	contract, err := bindSmartContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SmartContractCaller{contract: contract}, nil
}

// NewSmartContractTransactor creates a new write-only instance of SmartContract, bound to a specific deployed contract.
func NewSmartContractTransactor(address common.Address, transactor bind.ContractTransactor) (*SmartContractTransactor, error) {
	contract, err := bindSmartContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SmartContractTransactor{contract: contract}, nil
}

// NewSmartContractFilterer creates a new log filterer instance of SmartContract, bound to a specific deployed contract.
func NewSmartContractFilterer(address common.Address, filterer bind.ContractFilterer) (*SmartContractFilterer, error) {
	contract, err := bindSmartContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SmartContractFilterer{contract: contract}, nil
}

// bindSmartContract binds a generic wrapper to an already deployed contract.
func bindSmartContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SmartContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SmartContract *SmartContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SmartContract.Contract.SmartContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SmartContract *SmartContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmartContract.Contract.SmartContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SmartContract *SmartContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmartContract.Contract.SmartContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SmartContract *SmartContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SmartContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SmartContract *SmartContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SmartContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SmartContract *SmartContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SmartContract.Contract.contract.Transact(opts, method, params...)
}

// BatchCount is a free data retrieval call binding the contract method 0x06f13056.
//
// Solidity: function batchCount() constant returns(uint256)
func (_SmartContract *SmartContractCaller) BatchCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SmartContract.contract.Call(opts, out, "batchCount")
	return *ret0, err
}

// BatchCount is a free data retrieval call binding the contract method 0x06f13056.
//
// Solidity: function batchCount() constant returns(uint256)
func (_SmartContract *SmartContractSession) BatchCount() (*big.Int, error) {
	return _SmartContract.Contract.BatchCount(&_SmartContract.CallOpts)
}

// BatchCount is a free data retrieval call binding the contract method 0x06f13056.
//
// Solidity: function batchCount() constant returns(uint256)
func (_SmartContract *SmartContractCallerSession) BatchCount() (*big.Int, error) {
	return _SmartContract.Contract.BatchCount(&_SmartContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SmartContract *SmartContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SmartContract.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SmartContract *SmartContractSession) Owner() (common.Address, error) {
	return _SmartContract.Contract.Owner(&_SmartContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_SmartContract *SmartContractCallerSession) Owner() (common.Address, error) {
	return _SmartContract.Contract.Owner(&_SmartContract.CallOpts)
}

// CommitBatch is a paid mutator transaction binding the contract method 0xc9bacb66.
//
// Solidity: function commitBatch(uint256 batch) returns()
func (_SmartContract *SmartContractTransactor) CommitBatch(opts *bind.TransactOpts, batch *big.Int) (*types.Transaction, error) {
	return _SmartContract.contract.Transact(opts, "commitBatch", batch)
}

// CommitBatch is a paid mutator transaction binding the contract method 0xc9bacb66.
//
// Solidity: function commitBatch(uint256 batch) returns()
func (_SmartContract *SmartContractSession) CommitBatch(batch *big.Int) (*types.Transaction, error) {
	return _SmartContract.Contract.CommitBatch(&_SmartContract.TransactOpts, batch)
}

// CommitBatch is a paid mutator transaction binding the contract method 0xc9bacb66.
//
// Solidity: function commitBatch(uint256 batch) returns()
func (_SmartContract *SmartContractTransactorSession) CommitBatch(batch *big.Int) (*types.Transaction, error) {
	return _SmartContract.Contract.CommitBatch(&_SmartContract.TransactOpts, batch)
}
