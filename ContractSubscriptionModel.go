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

// SubscriptionModelMetaData contains all meta data concerning the SubscriptionModel contract.
var SubscriptionModelMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_subscriptionStorageAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_planStorageAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"ChangePlan\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_weeks\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Subscribe\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"changePlan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_weeks\",\"type\":\"uint256\"}],\"name\":\"subscribe\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405260286000557371398cad63b47db2e2b00b68a709b64df98e5a29600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550348015620000ce57600080fd5b506040516200108d3803806200108d8339818101604052810190620000f49190620001e8565b81600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050506200022f565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620001b08262000183565b9050919050565b620001c281620001a3565b8114620001ce57600080fd5b50565b600081519050620001e281620001b7565b92915050565b600080604083850312156200020257620002016200017e565b5b60006200021285828601620001d1565b92505060206200022585828601620001d1565b9150509250929050565b610e4e806200023f6000396000f3fe6080604052600436106100295760003560e01c806305d567ca1461002e5780638de6928414610057575b600080fd5b34801561003a57600080fd5b5061005560048036038101906100509190610678565b610073565b005b610071600480360381019061006c9190610703565b6102b9565b005b60003373ffffffffffffffffffffffffffffffffffffffff1631116100cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100c4906107a0565b60405180910390fd5b60008111610110576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161010790610832565b60405180910390fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dc42c074336040518263ffffffff1660e01b815260040161016b9190610873565b602060405180830381865afa158015610188573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101ac91906108a3565b8114156101ee576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101e590610942565b60405180910390fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d822745c33836040518363ffffffff1660e01b815260040161024b929190610971565b600060405180830381600087803b15801561026557600080fd5b505af1158015610279573d6000803e3d6000fd5b505050507f7f9c6471d241fdd651ec01d49cde6243143db4bbf23aac6d0573bceb342e674e33826040516102ae929190610971565b60405180910390a150565b60004290506000600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dc42c074856040518263ffffffff1660e01b815260040161031b91906109f9565b602060405180830381865afa158015610338573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061035c91906108a3565b9050600081116103a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161039890610a86565b60405180910390fd5b81600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635b75dd8d86336040518363ffffffff1660e01b81526004016103ff929190610aa6565b608060405180830381865afa15801561041c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104409190610bd8565b6040015110610484576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161047b90610c51565b60405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d5764858853384876040518563ffffffff1660e01b81526004016104e59493929190610c71565b600060405180830381600087803b1580156104ff57600080fd5b505af1158015610513573d6000803e3d6000fd5b505050507f7e70484266444fe9926cee86f5ca8c91acc579733e6c93081c0cbe7ab34877593385858460405161054c9493929190610cb6565b60405180910390a1600060646028346105659190610d2a565b61056f9190610db3565b9050600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f193505050501580156105d9573d6000803e3d6000fd5b508473ffffffffffffffffffffffffffffffffffffffff166108fc82346106009190610de4565b9081150290604051600060405180830381858888f1935050505015801561062b573d6000803e3d6000fd5b505050505050565b6000604051905090565b600080fd5b6000819050919050565b61065581610642565b811461066057600080fd5b50565b6000813590506106728161064c565b92915050565b60006020828403121561068e5761068d61063d565b5b600061069c84828501610663565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106d0826106a5565b9050919050565b6106e0816106c5565b81146106eb57600080fd5b50565b6000813590506106fd816106d7565b92915050565b6000806040838503121561071a5761071961063d565b5b6000610728858286016106ee565b925050602061073985828601610663565b9150509250929050565b600082825260208201905092915050565b7f42616c616e6365206973206e6f7420656e6f7567682e00000000000000000000600082015250565b600061078a601683610743565b915061079582610754565b602082019050919050565b600060208201905081810360008301526107b98161077d565b9050919050565b7f55534443207072696365206d75737420626520626967676572207468616e203060008201527f2e00000000000000000000000000000000000000000000000000000000000000602082015250565b600061081c602183610743565b9150610827826107c0565b604082019050919050565b6000602082019050818103600083015261084b8161080f565b9050919050565b600061085d826106a5565b9050919050565b61086d81610852565b82525050565b60006020820190506108886000830184610864565b92915050565b60008151905061089d8161064c565b92915050565b6000602082840312156108b9576108b861063d565b5b60006108c78482850161088e565b91505092915050565b7f55534443207072696365206d75737420626520646966666572656e742074686160008201527f6e2063757272656e74206f6e652e000000000000000000000000000000000000602082015250565b600061092c602e83610743565b9150610937826108d0565b604082019050919050565b6000602082019050818103600083015261095b8161091f565b9050919050565b61096b81610642565b82525050565b60006040820190506109866000830185610864565b6109936020830184610962565b9392505050565b6000819050919050565b60006109bf6109ba6109b5846106a5565b61099a565b6106a5565b9050919050565b60006109d1826109a4565b9050919050565b60006109e3826109c6565b9050919050565b6109f3816109d8565b82525050565b6000602082019050610a0e60008301846109ea565b92915050565b7f41646472657373206973206e6f7420616363657074696e67207375627363726960008201527f7074696f6e732e00000000000000000000000000000000000000000000000000602082015250565b6000610a70602783610743565b9150610a7b82610a14565b604082019050919050565b60006020820190508181036000830152610a9f81610a63565b9050919050565b6000604082019050610abb60008301856109ea565b610ac86020830184610864565b9392505050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610b1d82610ad4565b810181811067ffffffffffffffff82111715610b3c57610b3b610ae5565b5b80604052505050565b6000610b4f610633565b9050610b5b8282610b14565b919050565b600060808284031215610b7657610b75610acf565b5b610b806080610b45565b90506000610b908482850161088e565b6000830152506020610ba48482850161088e565b6020830152506040610bb88482850161088e565b6040830152506060610bcc8482850161088e565b60608301525092915050565b600060808284031215610bee57610bed61063d565b5b6000610bfc84828501610b60565b91505092915050565b7f546865726520697320612072756e6e696e6720737562736372697074696f6e00600082015250565b6000610c3b601f83610743565b9150610c4682610c05565b602082019050919050565b60006020820190508181036000830152610c6a81610c2e565b9050919050565b6000608082019050610c8660008301876109ea565b610c936020830186610864565b610ca06040830185610962565b610cad6060830184610962565b95945050505050565b6000608082019050610ccb6000830187610864565b610cd860208301866109ea565b610ce56040830185610962565b610cf26060830184610962565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610d3582610642565b9150610d4083610642565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610d7957610d78610cfb565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000610dbe82610642565b9150610dc983610642565b925082610dd957610dd8610d84565b5b828204905092915050565b6000610def82610642565b9150610dfa83610642565b925082821015610e0d57610e0c610cfb565b5b82820390509291505056fea264697066735822122050d61df5e3a79cfd6c3826762054715a629aee8bb5d29a3a7cc621bed9d219d864736f6c634300080b0033",
}

// SubscriptionModelABI is the input ABI used to generate the binding from.
// Deprecated: Use SubscriptionModelMetaData.ABI instead.
var SubscriptionModelABI = SubscriptionModelMetaData.ABI

// SubscriptionModelBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SubscriptionModelMetaData.Bin instead.
var SubscriptionModelBin = SubscriptionModelMetaData.Bin

// DeploySubscriptionModel deploys a new Ethereum contract, binding an instance of SubscriptionModel to it.
func DeploySubscriptionModel(auth *bind.TransactOpts, backend bind.ContractBackend, _subscriptionStorageAddress common.Address, _planStorageAddress common.Address) (common.Address, *types.Transaction, *SubscriptionModel, error) {
	parsed, err := SubscriptionModelMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SubscriptionModelBin), backend, _subscriptionStorageAddress, _planStorageAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SubscriptionModel{SubscriptionModelCaller: SubscriptionModelCaller{contract: contract}, SubscriptionModelTransactor: SubscriptionModelTransactor{contract: contract}, SubscriptionModelFilterer: SubscriptionModelFilterer{contract: contract}}, nil
}

// SubscriptionModel is an auto generated Go binding around an Ethereum contract.
type SubscriptionModel struct {
	SubscriptionModelCaller     // Read-only binding to the contract
	SubscriptionModelTransactor // Write-only binding to the contract
	SubscriptionModelFilterer   // Log filterer for contract events
}

// SubscriptionModelCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubscriptionModelCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionModelTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubscriptionModelTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionModelFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubscriptionModelFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionModelSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubscriptionModelSession struct {
	Contract     *SubscriptionModel // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SubscriptionModelCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubscriptionModelCallerSession struct {
	Contract *SubscriptionModelCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// SubscriptionModelTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubscriptionModelTransactorSession struct {
	Contract     *SubscriptionModelTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// SubscriptionModelRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubscriptionModelRaw struct {
	Contract *SubscriptionModel // Generic contract binding to access the raw methods on
}

// SubscriptionModelCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubscriptionModelCallerRaw struct {
	Contract *SubscriptionModelCaller // Generic read-only contract binding to access the raw methods on
}

// SubscriptionModelTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubscriptionModelTransactorRaw struct {
	Contract *SubscriptionModelTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubscriptionModel creates a new instance of SubscriptionModel, bound to a specific deployed contract.
func NewSubscriptionModel(address common.Address, backend bind.ContractBackend) (*SubscriptionModel, error) {
	contract, err := bindSubscriptionModel(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubscriptionModel{SubscriptionModelCaller: SubscriptionModelCaller{contract: contract}, SubscriptionModelTransactor: SubscriptionModelTransactor{contract: contract}, SubscriptionModelFilterer: SubscriptionModelFilterer{contract: contract}}, nil
}

// NewSubscriptionModelCaller creates a new read-only instance of SubscriptionModel, bound to a specific deployed contract.
func NewSubscriptionModelCaller(address common.Address, caller bind.ContractCaller) (*SubscriptionModelCaller, error) {
	contract, err := bindSubscriptionModel(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionModelCaller{contract: contract}, nil
}

// NewSubscriptionModelTransactor creates a new write-only instance of SubscriptionModel, bound to a specific deployed contract.
func NewSubscriptionModelTransactor(address common.Address, transactor bind.ContractTransactor) (*SubscriptionModelTransactor, error) {
	contract, err := bindSubscriptionModel(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionModelTransactor{contract: contract}, nil
}

// NewSubscriptionModelFilterer creates a new log filterer instance of SubscriptionModel, bound to a specific deployed contract.
func NewSubscriptionModelFilterer(address common.Address, filterer bind.ContractFilterer) (*SubscriptionModelFilterer, error) {
	contract, err := bindSubscriptionModel(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubscriptionModelFilterer{contract: contract}, nil
}

// bindSubscriptionModel binds a generic wrapper to an already deployed contract.
func bindSubscriptionModel(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SubscriptionModelABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubscriptionModel *SubscriptionModelRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionModel.Contract.SubscriptionModelCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubscriptionModel *SubscriptionModelRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.SubscriptionModelTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubscriptionModel *SubscriptionModelRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.SubscriptionModelTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubscriptionModel *SubscriptionModelCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionModel.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubscriptionModel *SubscriptionModelTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubscriptionModel *SubscriptionModelTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.contract.Transact(opts, method, params...)
}

// ChangePlan is a paid mutator transaction binding the contract method 0x05d567ca.
//
// Solidity: function changePlan(uint256 _price) returns()
func (_SubscriptionModel *SubscriptionModelTransactor) ChangePlan(opts *bind.TransactOpts, _price *big.Int) (*types.Transaction, error) {
	return _SubscriptionModel.contract.Transact(opts, "changePlan", _price)
}

// ChangePlan is a paid mutator transaction binding the contract method 0x05d567ca.
//
// Solidity: function changePlan(uint256 _price) returns()
func (_SubscriptionModel *SubscriptionModelSession) ChangePlan(_price *big.Int) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.ChangePlan(&_SubscriptionModel.TransactOpts, _price)
}

// ChangePlan is a paid mutator transaction binding the contract method 0x05d567ca.
//
// Solidity: function changePlan(uint256 _price) returns()
func (_SubscriptionModel *SubscriptionModelTransactorSession) ChangePlan(_price *big.Int) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.ChangePlan(&_SubscriptionModel.TransactOpts, _price)
}

// Subscribe is a paid mutator transaction binding the contract method 0x8de69284.
//
// Solidity: function subscribe(address _to, uint256 _weeks) payable returns()
func (_SubscriptionModel *SubscriptionModelTransactor) Subscribe(opts *bind.TransactOpts, _to common.Address, _weeks *big.Int) (*types.Transaction, error) {
	return _SubscriptionModel.contract.Transact(opts, "subscribe", _to, _weeks)
}

// Subscribe is a paid mutator transaction binding the contract method 0x8de69284.
//
// Solidity: function subscribe(address _to, uint256 _weeks) payable returns()
func (_SubscriptionModel *SubscriptionModelSession) Subscribe(_to common.Address, _weeks *big.Int) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.Subscribe(&_SubscriptionModel.TransactOpts, _to, _weeks)
}

// Subscribe is a paid mutator transaction binding the contract method 0x8de69284.
//
// Solidity: function subscribe(address _to, uint256 _weeks) payable returns()
func (_SubscriptionModel *SubscriptionModelTransactorSession) Subscribe(_to common.Address, _weeks *big.Int) (*types.Transaction, error) {
	return _SubscriptionModel.Contract.Subscribe(&_SubscriptionModel.TransactOpts, _to, _weeks)
}

// SubscriptionModelChangePlanIterator is returned from FilterChangePlan and is used to iterate over the raw logs and unpacked data for ChangePlan events raised by the SubscriptionModel contract.
type SubscriptionModelChangePlanIterator struct {
	Event *SubscriptionModelChangePlan // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SubscriptionModelChangePlanIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionModelChangePlan)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SubscriptionModelChangePlan)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SubscriptionModelChangePlanIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubscriptionModelChangePlanIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubscriptionModelChangePlan represents a ChangePlan event raised by the SubscriptionModel contract.
type SubscriptionModelChangePlan struct {
	Sender common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterChangePlan is a free log retrieval operation binding the contract event 0x7f9c6471d241fdd651ec01d49cde6243143db4bbf23aac6d0573bceb342e674e.
//
// Solidity: event ChangePlan(address sender, uint256 value)
func (_SubscriptionModel *SubscriptionModelFilterer) FilterChangePlan(opts *bind.FilterOpts) (*SubscriptionModelChangePlanIterator, error) {

	logs, sub, err := _SubscriptionModel.contract.FilterLogs(opts, "ChangePlan")
	if err != nil {
		return nil, err
	}
	return &SubscriptionModelChangePlanIterator{contract: _SubscriptionModel.contract, event: "ChangePlan", logs: logs, sub: sub}, nil
}

// WatchChangePlan is a free log subscription operation binding the contract event 0x7f9c6471d241fdd651ec01d49cde6243143db4bbf23aac6d0573bceb342e674e.
//
// Solidity: event ChangePlan(address sender, uint256 value)
func (_SubscriptionModel *SubscriptionModelFilterer) WatchChangePlan(opts *bind.WatchOpts, sink chan<- *SubscriptionModelChangePlan) (event.Subscription, error) {

	logs, sub, err := _SubscriptionModel.contract.WatchLogs(opts, "ChangePlan")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubscriptionModelChangePlan)
				if err := _SubscriptionModel.contract.UnpackLog(event, "ChangePlan", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChangePlan is a log parse operation binding the contract event 0x7f9c6471d241fdd651ec01d49cde6243143db4bbf23aac6d0573bceb342e674e.
//
// Solidity: event ChangePlan(address sender, uint256 value)
func (_SubscriptionModel *SubscriptionModelFilterer) ParseChangePlan(log types.Log) (*SubscriptionModelChangePlan, error) {
	event := new(SubscriptionModelChangePlan)
	if err := _SubscriptionModel.contract.UnpackLog(event, "ChangePlan", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubscriptionModelSubscribeIterator is returned from FilterSubscribe and is used to iterate over the raw logs and unpacked data for Subscribe events raised by the SubscriptionModel contract.
type SubscriptionModelSubscribeIterator struct {
	Event *SubscriptionModelSubscribe // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SubscriptionModelSubscribeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionModelSubscribe)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SubscriptionModelSubscribe)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SubscriptionModelSubscribeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubscriptionModelSubscribeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubscriptionModelSubscribe represents a Subscribe event raised by the SubscriptionModel contract.
type SubscriptionModelSubscribe struct {
	Sender common.Address
	To     common.Address
	Weeks  *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSubscribe is a free log retrieval operation binding the contract event 0x7e70484266444fe9926cee86f5ca8c91acc579733e6c93081c0cbe7ab3487759.
//
// Solidity: event Subscribe(address sender, address to, uint256 _weeks, uint256 amount)
func (_SubscriptionModel *SubscriptionModelFilterer) FilterSubscribe(opts *bind.FilterOpts) (*SubscriptionModelSubscribeIterator, error) {

	logs, sub, err := _SubscriptionModel.contract.FilterLogs(opts, "Subscribe")
	if err != nil {
		return nil, err
	}
	return &SubscriptionModelSubscribeIterator{contract: _SubscriptionModel.contract, event: "Subscribe", logs: logs, sub: sub}, nil
}

// WatchSubscribe is a free log subscription operation binding the contract event 0x7e70484266444fe9926cee86f5ca8c91acc579733e6c93081c0cbe7ab3487759.
//
// Solidity: event Subscribe(address sender, address to, uint256 _weeks, uint256 amount)
func (_SubscriptionModel *SubscriptionModelFilterer) WatchSubscribe(opts *bind.WatchOpts, sink chan<- *SubscriptionModelSubscribe) (event.Subscription, error) {

	logs, sub, err := _SubscriptionModel.contract.WatchLogs(opts, "Subscribe")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubscriptionModelSubscribe)
				if err := _SubscriptionModel.contract.UnpackLog(event, "Subscribe", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSubscribe is a log parse operation binding the contract event 0x7e70484266444fe9926cee86f5ca8c91acc579733e6c93081c0cbe7ab3487759.
//
// Solidity: event Subscribe(address sender, address to, uint256 _weeks, uint256 amount)
func (_SubscriptionModel *SubscriptionModelFilterer) ParseSubscribe(log types.Log) (*SubscriptionModelSubscribe, error) {
	event := new(SubscriptionModelSubscribe)
	if err := _SubscriptionModel.contract.UnpackLog(event, "Subscribe", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
