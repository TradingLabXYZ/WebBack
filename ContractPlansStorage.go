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

// PlansStorageMetaData contains all meta data concerning the PlansStorage contract.
var PlansStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"addPlan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"containsPlan\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowedContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"getPlanByAddress\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getPlanByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPlans\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"plans\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"Index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Value\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"removePlan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sizePlans\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_key\",\"type\":\"address\"}],\"name\":\"updateAllowedContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550611489806100616000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c8063c72adad811610066578063c72adad81461016e578063d822745c1461018a578063d94a862b146101a6578063dc42c074146101c4578063de4bb795146101f45761009e565b806309225872146100a35780637ad71e73146100c157806384fb8176146100dd57806386d86d841461010d578063b97e9d8d1461013e575b600080fd5b6100ab610212565b6040516100b8919061103a565b60405180910390f35b6100db60048036038101906100d69190611086565b610324565b005b6100f760048036038101906100f29190611086565b61066a565b60405161010491906110ce565b60405180910390f35b61012760048036038101906101229190611086565b6107a0565b604051610135929190611102565b60405180910390f35b61015860048036038101906101539190611157565b6107c4565b6040516101659190611184565b60405180910390f35b61018860048036038101906101839190611086565b610954565b005b6101a4600480360381019061019f919061119f565b610a80565b005b6101ae610c5b565b6040516101bb919061129d565b60405180910390f35b6101de60048036038101906101d99190611086565b610dd1565b6040516101eb9190611184565b60405180910390f35b6101fc610f04565b6040516102099190611184565b60405180910390f35b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806102bd5750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b6102fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102f39061131c565b60405180910390fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806103cd5750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b61040c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104039061131c565b60405180910390fd5b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020905060008160000154141561046057600080fd5b6003805490508160000154111561047657600080fd5b600060018260000154610489919061136b565b90506000600160038054905061049f919061136b565b90506001826104ae919061139f565b600080600384815481106104c5576104c46113f5565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000018190555060038181548110610545576105446113f5565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1660038381548110610584576105836113f5565b5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060038054806105de576105dd611424565b5b6001900381819060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905590556000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000808201600090556001820160009055505050505050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806107155750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610754576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161074b9061131c565b60405180910390fd5b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000154119050919050565b60006020528060005260406000206000915090508060000154908060010154905082565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061086f5750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b6108ae576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108a59061131c565b60405180910390fd5b60008210156108bc57600080fd5b60038054905082106108cd57600080fd5b600080600384815481106108e4576108e36113f5565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101549050919050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614806109fd5750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610a3c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a339061131c565b60405180910390fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610b295750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610b68576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b5f9061131c565b60405180910390fd5b60008060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209050818160010181905550600081600001541115610bc65750610c57565b6003839080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060006001600380549050610c3d919061136b565b9050600181610c4c919061139f565b826000018190555050505b5050565b6060600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610d065750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610d45576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d3c9061131c565b60405180910390fd5b6003805480602002602001604051908101604052809291908181526020018280548015610dc757602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610d7d575b5050505050905090565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610e7c5750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610ebb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610eb29061131c565b60405180910390fd5b6000808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600101549050919050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610faf5750600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b610fee576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610fe59061131c565b60405180910390fd5b600380549050905090565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061102482610ff9565b9050919050565b61103481611019565b82525050565b600060208201905061104f600083018461102b565b92915050565b600080fd5b61106381611019565b811461106e57600080fd5b50565b6000813590506110808161105a565b92915050565b60006020828403121561109c5761109b611055565b5b60006110aa84828501611071565b91505092915050565b60008115159050919050565b6110c8816110b3565b82525050565b60006020820190506110e360008301846110bf565b92915050565b6000819050919050565b6110fc816110e9565b82525050565b600060408201905061111760008301856110f3565b61112460208301846110f3565b9392505050565b611134816110e9565b811461113f57600080fd5b50565b6000813590506111518161112b565b92915050565b60006020828403121561116d5761116c611055565b5b600061117b84828501611142565b91505092915050565b600060208201905061119960008301846110f3565b92915050565b600080604083850312156111b6576111b5611055565b5b60006111c485828601611071565b92505060206111d585828601611142565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b61121481611019565b82525050565b6000611226838361120b565b60208301905092915050565b6000602082019050919050565b600061124a826111df565b61125481856111ea565b935061125f836111fb565b8060005b83811015611290578151611277888261121a565b975061128283611232565b925050600181019050611263565b5085935050505092915050565b600060208201905081810360008301526112b7818461123f565b905092915050565b600082825260208201905092915050565b7f4e6f7420616c6c6f7765642e0000000000000000000000000000000000000000600082015250565b6000611306600c836112bf565b9150611311826112d0565b602082019050919050565b60006020820190508181036000830152611335816112f9565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611376826110e9565b9150611381836110e9565b9250828210156113945761139361133c565b5b828203905092915050565b60006113aa826110e9565b91506113b5836110e9565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156113ea576113e961133c565b5b828201905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfea2646970667358221220f6755a4581a2ea01a632df481ae8361a22988aa7a7b1b20a4992184325de4abd64736f6c634300080b0033",
}

// PlansStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use PlansStorageMetaData.ABI instead.
var PlansStorageABI = PlansStorageMetaData.ABI

// PlansStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PlansStorageMetaData.Bin instead.
var PlansStorageBin = PlansStorageMetaData.Bin

// DeployPlansStorage deploys a new Ethereum contract, binding an instance of PlansStorage to it.
func DeployPlansStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PlansStorage, error) {
	parsed, err := PlansStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PlansStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PlansStorage{PlansStorageCaller: PlansStorageCaller{contract: contract}, PlansStorageTransactor: PlansStorageTransactor{contract: contract}, PlansStorageFilterer: PlansStorageFilterer{contract: contract}}, nil
}

// PlansStorage is an auto generated Go binding around an Ethereum contract.
type PlansStorage struct {
	PlansStorageCaller     // Read-only binding to the contract
	PlansStorageTransactor // Write-only binding to the contract
	PlansStorageFilterer   // Log filterer for contract events
}

// PlansStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type PlansStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlansStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PlansStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlansStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PlansStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlansStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PlansStorageSession struct {
	Contract     *PlansStorage     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PlansStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PlansStorageCallerSession struct {
	Contract *PlansStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PlansStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PlansStorageTransactorSession struct {
	Contract     *PlansStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PlansStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type PlansStorageRaw struct {
	Contract *PlansStorage // Generic contract binding to access the raw methods on
}

// PlansStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PlansStorageCallerRaw struct {
	Contract *PlansStorageCaller // Generic read-only contract binding to access the raw methods on
}

// PlansStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PlansStorageTransactorRaw struct {
	Contract *PlansStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPlansStorage creates a new instance of PlansStorage, bound to a specific deployed contract.
func NewPlansStorage(address common.Address, backend bind.ContractBackend) (*PlansStorage, error) {
	contract, err := bindPlansStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PlansStorage{PlansStorageCaller: PlansStorageCaller{contract: contract}, PlansStorageTransactor: PlansStorageTransactor{contract: contract}, PlansStorageFilterer: PlansStorageFilterer{contract: contract}}, nil
}

// NewPlansStorageCaller creates a new read-only instance of PlansStorage, bound to a specific deployed contract.
func NewPlansStorageCaller(address common.Address, caller bind.ContractCaller) (*PlansStorageCaller, error) {
	contract, err := bindPlansStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PlansStorageCaller{contract: contract}, nil
}

// NewPlansStorageTransactor creates a new write-only instance of PlansStorage, bound to a specific deployed contract.
func NewPlansStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*PlansStorageTransactor, error) {
	contract, err := bindPlansStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PlansStorageTransactor{contract: contract}, nil
}

// NewPlansStorageFilterer creates a new log filterer instance of PlansStorage, bound to a specific deployed contract.
func NewPlansStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*PlansStorageFilterer, error) {
	contract, err := bindPlansStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PlansStorageFilterer{contract: contract}, nil
}

// bindPlansStorage binds a generic wrapper to an already deployed contract.
func bindPlansStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PlansStorageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PlansStorage *PlansStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PlansStorage.Contract.PlansStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PlansStorage *PlansStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlansStorage.Contract.PlansStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PlansStorage *PlansStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PlansStorage.Contract.PlansStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PlansStorage *PlansStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PlansStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PlansStorage *PlansStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlansStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PlansStorage *PlansStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PlansStorage.Contract.contract.Transact(opts, method, params...)
}

// ContainsPlan is a free data retrieval call binding the contract method 0x84fb8176.
//
// Solidity: function containsPlan(address _key) view returns(bool)
func (_PlansStorage *PlansStorageCaller) ContainsPlan(opts *bind.CallOpts, _key common.Address) (bool, error) {
	var out []interface{}
	err := _PlansStorage.contract.Call(opts, &out, "containsPlan", _key)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ContainsPlan is a free data retrieval call binding the contract method 0x84fb8176.
//
// Solidity: function containsPlan(address _key) view returns(bool)
func (_PlansStorage *PlansStorageSession) ContainsPlan(_key common.Address) (bool, error) {
	return _PlansStorage.Contract.ContainsPlan(&_PlansStorage.CallOpts, _key)
}

// ContainsPlan is a free data retrieval call binding the contract method 0x84fb8176.
//
// Solidity: function containsPlan(address _key) view returns(bool)
func (_PlansStorage *PlansStorageCallerSession) ContainsPlan(_key common.Address) (bool, error) {
	return _PlansStorage.Contract.ContainsPlan(&_PlansStorage.CallOpts, _key)
}

// GetAllowedContract is a free data retrieval call binding the contract method 0x09225872.
//
// Solidity: function getAllowedContract() view returns(address)
func (_PlansStorage *PlansStorageCaller) GetAllowedContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PlansStorage.contract.Call(opts, &out, "getAllowedContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAllowedContract is a free data retrieval call binding the contract method 0x09225872.
//
// Solidity: function getAllowedContract() view returns(address)
func (_PlansStorage *PlansStorageSession) GetAllowedContract() (common.Address, error) {
	return _PlansStorage.Contract.GetAllowedContract(&_PlansStorage.CallOpts)
}

// GetAllowedContract is a free data retrieval call binding the contract method 0x09225872.
//
// Solidity: function getAllowedContract() view returns(address)
func (_PlansStorage *PlansStorageCallerSession) GetAllowedContract() (common.Address, error) {
	return _PlansStorage.Contract.GetAllowedContract(&_PlansStorage.CallOpts)
}

// GetPlanByAddress is a free data retrieval call binding the contract method 0xdc42c074.
//
// Solidity: function getPlanByAddress(address _key) view returns(uint256)
func (_PlansStorage *PlansStorageCaller) GetPlanByAddress(opts *bind.CallOpts, _key common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PlansStorage.contract.Call(opts, &out, "getPlanByAddress", _key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlanByAddress is a free data retrieval call binding the contract method 0xdc42c074.
//
// Solidity: function getPlanByAddress(address _key) view returns(uint256)
func (_PlansStorage *PlansStorageSession) GetPlanByAddress(_key common.Address) (*big.Int, error) {
	return _PlansStorage.Contract.GetPlanByAddress(&_PlansStorage.CallOpts, _key)
}

// GetPlanByAddress is a free data retrieval call binding the contract method 0xdc42c074.
//
// Solidity: function getPlanByAddress(address _key) view returns(uint256)
func (_PlansStorage *PlansStorageCallerSession) GetPlanByAddress(_key common.Address) (*big.Int, error) {
	return _PlansStorage.Contract.GetPlanByAddress(&_PlansStorage.CallOpts, _key)
}

// GetPlanByIndex is a free data retrieval call binding the contract method 0xb97e9d8d.
//
// Solidity: function getPlanByIndex(uint256 _index) view returns(uint256)
func (_PlansStorage *PlansStorageCaller) GetPlanByIndex(opts *bind.CallOpts, _index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PlansStorage.contract.Call(opts, &out, "getPlanByIndex", _index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlanByIndex is a free data retrieval call binding the contract method 0xb97e9d8d.
//
// Solidity: function getPlanByIndex(uint256 _index) view returns(uint256)
func (_PlansStorage *PlansStorageSession) GetPlanByIndex(_index *big.Int) (*big.Int, error) {
	return _PlansStorage.Contract.GetPlanByIndex(&_PlansStorage.CallOpts, _index)
}

// GetPlanByIndex is a free data retrieval call binding the contract method 0xb97e9d8d.
//
// Solidity: function getPlanByIndex(uint256 _index) view returns(uint256)
func (_PlansStorage *PlansStorageCallerSession) GetPlanByIndex(_index *big.Int) (*big.Int, error) {
	return _PlansStorage.Contract.GetPlanByIndex(&_PlansStorage.CallOpts, _index)
}

// GetPlans is a free data retrieval call binding the contract method 0xd94a862b.
//
// Solidity: function getPlans() view returns(address[])
func (_PlansStorage *PlansStorageCaller) GetPlans(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _PlansStorage.contract.Call(opts, &out, "getPlans")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetPlans is a free data retrieval call binding the contract method 0xd94a862b.
//
// Solidity: function getPlans() view returns(address[])
func (_PlansStorage *PlansStorageSession) GetPlans() ([]common.Address, error) {
	return _PlansStorage.Contract.GetPlans(&_PlansStorage.CallOpts)
}

// GetPlans is a free data retrieval call binding the contract method 0xd94a862b.
//
// Solidity: function getPlans() view returns(address[])
func (_PlansStorage *PlansStorageCallerSession) GetPlans() ([]common.Address, error) {
	return _PlansStorage.Contract.GetPlans(&_PlansStorage.CallOpts)
}

// Plans is a free data retrieval call binding the contract method 0x86d86d84.
//
// Solidity: function plans(address ) view returns(uint256 Index, uint256 Value)
func (_PlansStorage *PlansStorageCaller) Plans(opts *bind.CallOpts, arg0 common.Address) (struct {
	Index *big.Int
	Value *big.Int
}, error) {
	var out []interface{}
	err := _PlansStorage.contract.Call(opts, &out, "plans", arg0)

	outstruct := new(struct {
		Index *big.Int
		Value *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Index = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Value = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Plans is a free data retrieval call binding the contract method 0x86d86d84.
//
// Solidity: function plans(address ) view returns(uint256 Index, uint256 Value)
func (_PlansStorage *PlansStorageSession) Plans(arg0 common.Address) (struct {
	Index *big.Int
	Value *big.Int
}, error) {
	return _PlansStorage.Contract.Plans(&_PlansStorage.CallOpts, arg0)
}

// Plans is a free data retrieval call binding the contract method 0x86d86d84.
//
// Solidity: function plans(address ) view returns(uint256 Index, uint256 Value)
func (_PlansStorage *PlansStorageCallerSession) Plans(arg0 common.Address) (struct {
	Index *big.Int
	Value *big.Int
}, error) {
	return _PlansStorage.Contract.Plans(&_PlansStorage.CallOpts, arg0)
}

// SizePlans is a free data retrieval call binding the contract method 0xde4bb795.
//
// Solidity: function sizePlans() view returns(uint256)
func (_PlansStorage *PlansStorageCaller) SizePlans(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PlansStorage.contract.Call(opts, &out, "sizePlans")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SizePlans is a free data retrieval call binding the contract method 0xde4bb795.
//
// Solidity: function sizePlans() view returns(uint256)
func (_PlansStorage *PlansStorageSession) SizePlans() (*big.Int, error) {
	return _PlansStorage.Contract.SizePlans(&_PlansStorage.CallOpts)
}

// SizePlans is a free data retrieval call binding the contract method 0xde4bb795.
//
// Solidity: function sizePlans() view returns(uint256)
func (_PlansStorage *PlansStorageCallerSession) SizePlans() (*big.Int, error) {
	return _PlansStorage.Contract.SizePlans(&_PlansStorage.CallOpts)
}

// AddPlan is a paid mutator transaction binding the contract method 0xd822745c.
//
// Solidity: function addPlan(address _key, uint256 _value) returns()
func (_PlansStorage *PlansStorageTransactor) AddPlan(opts *bind.TransactOpts, _key common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PlansStorage.contract.Transact(opts, "addPlan", _key, _value)
}

// AddPlan is a paid mutator transaction binding the contract method 0xd822745c.
//
// Solidity: function addPlan(address _key, uint256 _value) returns()
func (_PlansStorage *PlansStorageSession) AddPlan(_key common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PlansStorage.Contract.AddPlan(&_PlansStorage.TransactOpts, _key, _value)
}

// AddPlan is a paid mutator transaction binding the contract method 0xd822745c.
//
// Solidity: function addPlan(address _key, uint256 _value) returns()
func (_PlansStorage *PlansStorageTransactorSession) AddPlan(_key common.Address, _value *big.Int) (*types.Transaction, error) {
	return _PlansStorage.Contract.AddPlan(&_PlansStorage.TransactOpts, _key, _value)
}

// RemovePlan is a paid mutator transaction binding the contract method 0x7ad71e73.
//
// Solidity: function removePlan(address _key) returns()
func (_PlansStorage *PlansStorageTransactor) RemovePlan(opts *bind.TransactOpts, _key common.Address) (*types.Transaction, error) {
	return _PlansStorage.contract.Transact(opts, "removePlan", _key)
}

// RemovePlan is a paid mutator transaction binding the contract method 0x7ad71e73.
//
// Solidity: function removePlan(address _key) returns()
func (_PlansStorage *PlansStorageSession) RemovePlan(_key common.Address) (*types.Transaction, error) {
	return _PlansStorage.Contract.RemovePlan(&_PlansStorage.TransactOpts, _key)
}

// RemovePlan is a paid mutator transaction binding the contract method 0x7ad71e73.
//
// Solidity: function removePlan(address _key) returns()
func (_PlansStorage *PlansStorageTransactorSession) RemovePlan(_key common.Address) (*types.Transaction, error) {
	return _PlansStorage.Contract.RemovePlan(&_PlansStorage.TransactOpts, _key)
}

// UpdateAllowedContract is a paid mutator transaction binding the contract method 0xc72adad8.
//
// Solidity: function updateAllowedContract(address _key) returns()
func (_PlansStorage *PlansStorageTransactor) UpdateAllowedContract(opts *bind.TransactOpts, _key common.Address) (*types.Transaction, error) {
	return _PlansStorage.contract.Transact(opts, "updateAllowedContract", _key)
}

// UpdateAllowedContract is a paid mutator transaction binding the contract method 0xc72adad8.
//
// Solidity: function updateAllowedContract(address _key) returns()
func (_PlansStorage *PlansStorageSession) UpdateAllowedContract(_key common.Address) (*types.Transaction, error) {
	return _PlansStorage.Contract.UpdateAllowedContract(&_PlansStorage.TransactOpts, _key)
}

// UpdateAllowedContract is a paid mutator transaction binding the contract method 0xc72adad8.
//
// Solidity: function updateAllowedContract(address _key) returns()
func (_PlansStorage *PlansStorageTransactorSession) UpdateAllowedContract(_key common.Address) (*types.Transaction, error) {
	return _PlansStorage.Contract.UpdateAllowedContract(&_PlansStorage.TransactOpts, _key)
}
