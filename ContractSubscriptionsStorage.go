// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SubscriptionsStorageSubscription is an auto generated low-level Go binding around an user-defined struct.
type SubscriptionsStorageSubscription struct {
	Index     *big.Int
	Createdat *big.Int
	Endedat   *big.Int
	Amount    *big.Int
}

// SubscriptionsStorageMetaData contains all meta data concerning the SubscriptionsStorage contract.
var SubscriptionsStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_weeks\",\"type\":\"uint256\"}],\"name\":\"addSubscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"}],\"name\":\"getCountSubscriptionsBySubscriptor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"getSubscription\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"Index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Createdat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Endedat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Amount\",\"type\":\"uint256\"}],\"internalType\":\"structSubscriptionsStorage.Subscription\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubscriptors\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sizeSubscriptions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"subscriptions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"Index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Createdat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Endedat\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"updateAllowedContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555061119c806100616000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c80639974db861161005b5780639974db861461012c578063c72adad81461015c578063c958a0a614610178578063d57648581461019657610088565b80630119cbb41461008d57806309225872146100ab5780635b75dd8d146100c95780636207d31b146100f9575b600080fd5b6100956101b2565b6040516100a29190610c8a565b60405180910390f35b6100b36102a7565b6040516100c09190610ce6565b60405180910390f35b6100e360048036038101906100de9190610d32565b6103b9565b6040516100f09190610dd6565b60405180910390f35b610113600480360381019061010e9190610d32565b61055f565b6040516101239493929190610df1565b60405180910390f35b61014660048036038101906101419190610e36565b61059c565b6040516101539190610c8a565b60405180910390f35b61017660048036038101906101719190610e36565b6106cd565b005b6101806107f9565b60405161018d9190610f21565b60405180910390f35b6101b060048036038101906101ab9190610f6f565b61096f565b005b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061025d5750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b61029c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029390611033565b60405180910390fd5b600480549050905090565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806103525750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610391576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161038890611033565b60405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6103c1610c49565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061046a5750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b6104a9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104a090611033565b60405180910390fd5b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020604051806080016040529081600082015481526020016001820154815260200160028201548152602001600382015481525050905092915050565b6000602052816000526040600020602052806000526040600020600091509150508060000154908060010154908060020154908060030154905084565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806106475750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610686576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067d90611033565b60405180910390fd5b600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806107765750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b6107b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107ac90611033565b60405180910390fd5b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6060600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806108a45750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b6108e3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108da90611033565b60405180910390fd5b600480548060200260200160405190810160405280929190818152602001828054801561096557602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001906001019080831161091b575b5050505050905090565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610a185750600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610a57576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a4e90611033565b60405180910390fd5b60008060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209050600042905080826001018190555062093a8083610af39190611082565b81610afe91906110dc565b8260020181905550838260030181905550600082600001541115610b23575050610c43565b6004869080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060006001600480549050610b9a9190611132565b9050600181610ba991906110dc565b836000018190555060018060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610bfc91906110dc565b600160008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505050505b50505050565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000819050919050565b610c8481610c71565b82525050565b6000602082019050610c9f6000830184610c7b565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610cd082610ca5565b9050919050565b610ce081610cc5565b82525050565b6000602082019050610cfb6000830184610cd7565b92915050565b600080fd5b610d0f81610cc5565b8114610d1a57600080fd5b50565b600081359050610d2c81610d06565b92915050565b60008060408385031215610d4957610d48610d01565b5b6000610d5785828601610d1d565b9250506020610d6885828601610d1d565b9150509250929050565b610d7b81610c71565b82525050565b608082016000820151610d976000850182610d72565b506020820151610daa6020850182610d72565b506040820151610dbd6040850182610d72565b506060820151610dd06060850182610d72565b50505050565b6000608082019050610deb6000830184610d81565b92915050565b6000608082019050610e066000830187610c7b565b610e136020830186610c7b565b610e206040830185610c7b565b610e2d6060830184610c7b565b95945050505050565b600060208284031215610e4c57610e4b610d01565b5b6000610e5a84828501610d1d565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b610e9881610cc5565b82525050565b6000610eaa8383610e8f565b60208301905092915050565b6000602082019050919050565b6000610ece82610e63565b610ed88185610e6e565b9350610ee383610e7f565b8060005b83811015610f14578151610efb8882610e9e565b9750610f0683610eb6565b925050600181019050610ee7565b5085935050505092915050565b60006020820190508181036000830152610f3b8184610ec3565b905092915050565b610f4c81610c71565b8114610f5757600080fd5b50565b600081359050610f6981610f43565b92915050565b60008060008060808587031215610f8957610f88610d01565b5b6000610f9787828801610d1d565b9450506020610fa887828801610d1d565b9350506040610fb987828801610f5a565b9250506060610fca87828801610f5a565b91505092959194509250565b600082825260208201905092915050565b7f4e6f7420616c6c6f7765642e0000000000000000000000000000000000000000600082015250565b600061101d600c83610fd6565b915061102882610fe7565b602082019050919050565b6000602082019050818103600083015261104c81611010565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061108d82610c71565b915061109883610c71565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156110d1576110d0611053565b5b828202905092915050565b60006110e782610c71565b91506110f283610c71565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561112757611126611053565b5b828201905092915050565b600061113d82610c71565b915061114883610c71565b92508282101561115b5761115a611053565b5b82820390509291505056fea26469706673582212207fb8372857e9b2fb41995f3f930c1cc4aea2f381717a5e65d17be2e37c960f4864736f6c634300080b0033",
}

// SubscriptionsStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use SubscriptionsStorageMetaData.ABI instead.
var SubscriptionsStorageABI = SubscriptionsStorageMetaData.ABI

// SubscriptionsStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SubscriptionsStorageMetaData.Bin instead.
var SubscriptionsStorageBin = SubscriptionsStorageMetaData.Bin

// DeploySubscriptionsStorage deploys a new Ethereum contract, binding an instance of SubscriptionsStorage to it.
func DeploySubscriptionsStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SubscriptionsStorage, error) {
	parsed, err := SubscriptionsStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SubscriptionsStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SubscriptionsStorage{SubscriptionsStorageCaller: SubscriptionsStorageCaller{contract: contract}, SubscriptionsStorageTransactor: SubscriptionsStorageTransactor{contract: contract}, SubscriptionsStorageFilterer: SubscriptionsStorageFilterer{contract: contract}}, nil
}

// SubscriptionsStorage is an auto generated Go binding around an Ethereum contract.
type SubscriptionsStorage struct {
	SubscriptionsStorageCaller     // Read-only binding to the contract
	SubscriptionsStorageTransactor // Write-only binding to the contract
	SubscriptionsStorageFilterer   // Log filterer for contract events
}

// SubscriptionsStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubscriptionsStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionsStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubscriptionsStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionsStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubscriptionsStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionsStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubscriptionsStorageSession struct {
	Contract     *SubscriptionsStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SubscriptionsStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubscriptionsStorageCallerSession struct {
	Contract *SubscriptionsStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// SubscriptionsStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubscriptionsStorageTransactorSession struct {
	Contract     *SubscriptionsStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// SubscriptionsStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubscriptionsStorageRaw struct {
	Contract *SubscriptionsStorage // Generic contract binding to access the raw methods on
}

// SubscriptionsStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubscriptionsStorageCallerRaw struct {
	Contract *SubscriptionsStorageCaller // Generic read-only contract binding to access the raw methods on
}

// SubscriptionsStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubscriptionsStorageTransactorRaw struct {
	Contract *SubscriptionsStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubscriptionsStorage creates a new instance of SubscriptionsStorage, bound to a specific deployed contract.
func NewSubscriptionsStorage(address common.Address, backend bind.ContractBackend) (*SubscriptionsStorage, error) {
	contract, err := bindSubscriptionsStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubscriptionsStorage{SubscriptionsStorageCaller: SubscriptionsStorageCaller{contract: contract}, SubscriptionsStorageTransactor: SubscriptionsStorageTransactor{contract: contract}, SubscriptionsStorageFilterer: SubscriptionsStorageFilterer{contract: contract}}, nil
}

// NewSubscriptionsStorageCaller creates a new read-only instance of SubscriptionsStorage, bound to a specific deployed contract.
func NewSubscriptionsStorageCaller(address common.Address, caller bind.ContractCaller) (*SubscriptionsStorageCaller, error) {
	contract, err := bindSubscriptionsStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionsStorageCaller{contract: contract}, nil
}

// NewSubscriptionsStorageTransactor creates a new write-only instance of SubscriptionsStorage, bound to a specific deployed contract.
func NewSubscriptionsStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*SubscriptionsStorageTransactor, error) {
	contract, err := bindSubscriptionsStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionsStorageTransactor{contract: contract}, nil
}

// NewSubscriptionsStorageFilterer creates a new log filterer instance of SubscriptionsStorage, bound to a specific deployed contract.
func NewSubscriptionsStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*SubscriptionsStorageFilterer, error) {
	contract, err := bindSubscriptionsStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubscriptionsStorageFilterer{contract: contract}, nil
}

// bindSubscriptionsStorage binds a generic wrapper to an already deployed contract.
func bindSubscriptionsStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SubscriptionsStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubscriptionsStorage *SubscriptionsStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionsStorage.Contract.SubscriptionsStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubscriptionsStorage *SubscriptionsStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.SubscriptionsStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubscriptionsStorage *SubscriptionsStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.SubscriptionsStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubscriptionsStorage *SubscriptionsStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionsStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubscriptionsStorage *SubscriptionsStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubscriptionsStorage *SubscriptionsStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.contract.Transact(opts, method, params...)
}

// GetAllowedContract is a free data retrieval call binding the contract method 0x09225872.
//
// Solidity: function getAllowedContract() view returns(address)
func (_SubscriptionsStorage *SubscriptionsStorageCaller) GetAllowedContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SubscriptionsStorage.contract.Call(opts, &out, "getAllowedContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAllowedContract is a free data retrieval call binding the contract method 0x09225872.
//
// Solidity: function getAllowedContract() view returns(address)
func (_SubscriptionsStorage *SubscriptionsStorageSession) GetAllowedContract() (common.Address, error) {
	return _SubscriptionsStorage.Contract.GetAllowedContract(&_SubscriptionsStorage.CallOpts)
}

// GetAllowedContract is a free data retrieval call binding the contract method 0x09225872.
//
// Solidity: function getAllowedContract() view returns(address)
func (_SubscriptionsStorage *SubscriptionsStorageCallerSession) GetAllowedContract() (common.Address, error) {
	return _SubscriptionsStorage.Contract.GetAllowedContract(&_SubscriptionsStorage.CallOpts)
}

// GetCountSubscriptionsBySubscriptor is a free data retrieval call binding the contract method 0x9974db86.
//
// Solidity: function getCountSubscriptionsBySubscriptor(address _from) view returns(uint256)
func (_SubscriptionsStorage *SubscriptionsStorageCaller) GetCountSubscriptionsBySubscriptor(opts *bind.CallOpts, _from common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SubscriptionsStorage.contract.Call(opts, &out, "getCountSubscriptionsBySubscriptor", _from)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCountSubscriptionsBySubscriptor is a free data retrieval call binding the contract method 0x9974db86.
//
// Solidity: function getCountSubscriptionsBySubscriptor(address _from) view returns(uint256)
func (_SubscriptionsStorage *SubscriptionsStorageSession) GetCountSubscriptionsBySubscriptor(_from common.Address) (*big.Int, error) {
	return _SubscriptionsStorage.Contract.GetCountSubscriptionsBySubscriptor(&_SubscriptionsStorage.CallOpts, _from)
}

// GetCountSubscriptionsBySubscriptor is a free data retrieval call binding the contract method 0x9974db86.
//
// Solidity: function getCountSubscriptionsBySubscriptor(address _from) view returns(uint256)
func (_SubscriptionsStorage *SubscriptionsStorageCallerSession) GetCountSubscriptionsBySubscriptor(_from common.Address) (*big.Int, error) {
	return _SubscriptionsStorage.Contract.GetCountSubscriptionsBySubscriptor(&_SubscriptionsStorage.CallOpts, _from)
}

// GetSubscription is a free data retrieval call binding the contract method 0x5b75dd8d.
//
// Solidity: function getSubscription(address _from, address _to) view returns((uint256,uint256,uint256,uint256))
func (_SubscriptionsStorage *SubscriptionsStorageCaller) GetSubscription(opts *bind.CallOpts, _from common.Address, _to common.Address) (SubscriptionsStorageSubscription, error) {
	var out []interface{}
	err := _SubscriptionsStorage.contract.Call(opts, &out, "getSubscription", _from, _to)

	if err != nil {
		return *new(SubscriptionsStorageSubscription), err
	}

	out0 := *abi.ConvertType(out[0], new(SubscriptionsStorageSubscription)).(*SubscriptionsStorageSubscription)

	return out0, err

}

// GetSubscription is a free data retrieval call binding the contract method 0x5b75dd8d.
//
// Solidity: function getSubscription(address _from, address _to) view returns((uint256,uint256,uint256,uint256))
func (_SubscriptionsStorage *SubscriptionsStorageSession) GetSubscription(_from common.Address, _to common.Address) (SubscriptionsStorageSubscription, error) {
	return _SubscriptionsStorage.Contract.GetSubscription(&_SubscriptionsStorage.CallOpts, _from, _to)
}

// GetSubscription is a free data retrieval call binding the contract method 0x5b75dd8d.
//
// Solidity: function getSubscription(address _from, address _to) view returns((uint256,uint256,uint256,uint256))
func (_SubscriptionsStorage *SubscriptionsStorageCallerSession) GetSubscription(_from common.Address, _to common.Address) (SubscriptionsStorageSubscription, error) {
	return _SubscriptionsStorage.Contract.GetSubscription(&_SubscriptionsStorage.CallOpts, _from, _to)
}

// GetSubscriptors is a free data retrieval call binding the contract method 0xc958a0a6.
//
// Solidity: function getSubscriptors() view returns(address[])
func (_SubscriptionsStorage *SubscriptionsStorageCaller) GetSubscriptors(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _SubscriptionsStorage.contract.Call(opts, &out, "getSubscriptors")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetSubscriptors is a free data retrieval call binding the contract method 0xc958a0a6.
//
// Solidity: function getSubscriptors() view returns(address[])
func (_SubscriptionsStorage *SubscriptionsStorageSession) GetSubscriptors() ([]common.Address, error) {
	return _SubscriptionsStorage.Contract.GetSubscriptors(&_SubscriptionsStorage.CallOpts)
}

// GetSubscriptors is a free data retrieval call binding the contract method 0xc958a0a6.
//
// Solidity: function getSubscriptors() view returns(address[])
func (_SubscriptionsStorage *SubscriptionsStorageCallerSession) GetSubscriptors() ([]common.Address, error) {
	return _SubscriptionsStorage.Contract.GetSubscriptors(&_SubscriptionsStorage.CallOpts)
}

// SizeSubscriptions is a free data retrieval call binding the contract method 0x0119cbb4.
//
// Solidity: function sizeSubscriptions() view returns(uint256)
func (_SubscriptionsStorage *SubscriptionsStorageCaller) SizeSubscriptions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SubscriptionsStorage.contract.Call(opts, &out, "sizeSubscriptions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SizeSubscriptions is a free data retrieval call binding the contract method 0x0119cbb4.
//
// Solidity: function sizeSubscriptions() view returns(uint256)
func (_SubscriptionsStorage *SubscriptionsStorageSession) SizeSubscriptions() (*big.Int, error) {
	return _SubscriptionsStorage.Contract.SizeSubscriptions(&_SubscriptionsStorage.CallOpts)
}

// SizeSubscriptions is a free data retrieval call binding the contract method 0x0119cbb4.
//
// Solidity: function sizeSubscriptions() view returns(uint256)
func (_SubscriptionsStorage *SubscriptionsStorageCallerSession) SizeSubscriptions() (*big.Int, error) {
	return _SubscriptionsStorage.Contract.SizeSubscriptions(&_SubscriptionsStorage.CallOpts)
}

// Subscriptions is a free data retrieval call binding the contract method 0x6207d31b.
//
// Solidity: function subscriptions(address , address ) view returns(uint256 Index, uint256 Createdat, uint256 Endedat, uint256 Amount)
func (_SubscriptionsStorage *SubscriptionsStorageCaller) Subscriptions(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (struct {
	Index     *big.Int
	Createdat *big.Int
	Endedat   *big.Int
	Amount    *big.Int
}, error) {
	var out []interface{}
	err := _SubscriptionsStorage.contract.Call(opts, &out, "subscriptions", arg0, arg1)

	outstruct := new(struct {
		Index     *big.Int
		Createdat *big.Int
		Endedat   *big.Int
		Amount    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Createdat = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Endedat = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Subscriptions is a free data retrieval call binding the contract method 0x6207d31b.
//
// Solidity: function subscriptions(address , address ) view returns(uint256 Index, uint256 Createdat, uint256 Endedat, uint256 Amount)
func (_SubscriptionsStorage *SubscriptionsStorageSession) Subscriptions(arg0 common.Address, arg1 common.Address) (struct {
	Index     *big.Int
	Createdat *big.Int
	Endedat   *big.Int
	Amount    *big.Int
}, error) {
	return _SubscriptionsStorage.Contract.Subscriptions(&_SubscriptionsStorage.CallOpts, arg0, arg1)
}

// Subscriptions is a free data retrieval call binding the contract method 0x6207d31b.
//
// Solidity: function subscriptions(address , address ) view returns(uint256 Index, uint256 Createdat, uint256 Endedat, uint256 Amount)
func (_SubscriptionsStorage *SubscriptionsStorageCallerSession) Subscriptions(arg0 common.Address, arg1 common.Address) (struct {
	Index     *big.Int
	Createdat *big.Int
	Endedat   *big.Int
	Amount    *big.Int
}, error) {
	return _SubscriptionsStorage.Contract.Subscriptions(&_SubscriptionsStorage.CallOpts, arg0, arg1)
}

// AddSubscription is a paid mutator transaction binding the contract method 0xd5764858.
//
// Solidity: function addSubscription(address _from, address _to, uint256 _amount, uint256 _weeks) returns()
func (_SubscriptionsStorage *SubscriptionsStorageTransactor) AddSubscription(opts *bind.TransactOpts, _from common.Address, _to common.Address, _amount *big.Int, _weeks *big.Int) (*types.Transaction, error) {
	return _SubscriptionsStorage.contract.Transact(opts, "addSubscription", _from, _to, _amount, _weeks)
}

// AddSubscription is a paid mutator transaction binding the contract method 0xd5764858.
//
// Solidity: function addSubscription(address _from, address _to, uint256 _amount, uint256 _weeks) returns()
func (_SubscriptionsStorage *SubscriptionsStorageSession) AddSubscription(_from common.Address, _to common.Address, _amount *big.Int, _weeks *big.Int) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.AddSubscription(&_SubscriptionsStorage.TransactOpts, _from, _to, _amount, _weeks)
}

// AddSubscription is a paid mutator transaction binding the contract method 0xd5764858.
//
// Solidity: function addSubscription(address _from, address _to, uint256 _amount, uint256 _weeks) returns()
func (_SubscriptionsStorage *SubscriptionsStorageTransactorSession) AddSubscription(_from common.Address, _to common.Address, _amount *big.Int, _weeks *big.Int) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.AddSubscription(&_SubscriptionsStorage.TransactOpts, _from, _to, _amount, _weeks)
}

// UpdateAllowedContract is a paid mutator transaction binding the contract method 0xc72adad8.
//
// Solidity: function updateAllowedContract(address _key) returns()
func (_SubscriptionsStorage *SubscriptionsStorageTransactor) UpdateAllowedContract(opts *bind.TransactOpts, _key common.Address) (*types.Transaction, error) {
	return _SubscriptionsStorage.contract.Transact(opts, "updateAllowedContract", _key)
}

// UpdateAllowedContract is a paid mutator transaction binding the contract method 0xc72adad8.
//
// Solidity: function updateAllowedContract(address _key) returns()
func (_SubscriptionsStorage *SubscriptionsStorageSession) UpdateAllowedContract(_key common.Address) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.UpdateAllowedContract(&_SubscriptionsStorage.TransactOpts, _key)
}

// UpdateAllowedContract is a paid mutator transaction binding the contract method 0xc72adad8.
//
// Solidity: function updateAllowedContract(address _key) returns()
func (_SubscriptionsStorage *SubscriptionsStorageTransactorSession) UpdateAllowedContract(_key common.Address) (*types.Transaction, error) {
	return _SubscriptionsStorage.Contract.UpdateAllowedContract(&_SubscriptionsStorage.TransactOpts, _key)
}
