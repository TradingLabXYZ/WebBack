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
	Bin: "0x608060405260286000557371398cad63b47db2e2b00b68a709b64df98e5a29600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550348015620000ce57600080fd5b5060405162001177380380620011778339818101604052810190620000f49190620001e8565b81600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050506200022f565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620001b08262000183565b9050919050565b620001c281620001a3565b8114620001ce57600080fd5b50565b600081519050620001e281620001b7565b92915050565b600080604083850312156200020257620002016200017e565b5b60006200021285828601620001d1565b92505060206200022585828601620001d1565b9150509250929050565b610f38806200023f6000396000f3fe6080604052600436106100295760003560e01c806305d567ca1461002e5780638de6928414610057575b600080fd5b34801561003a57600080fd5b50610055600480360381019061005091906106d0565b610073565b005b610071600480360381019061006c919061075b565b61025f565b005b600081116100b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100ad9061081e565b60405180910390fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dc42c074336040518263ffffffff1660e01b8152600401610111919061085f565b602060405180830381865afa15801561012e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610152919061088f565b811415610194576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161018b9061092e565b60405180910390fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d822745c33836040518363ffffffff1660e01b81526004016101f192919061095d565b600060405180830381600087803b15801561020b57600080fd5b505af115801561021f573d6000803e3d6000fd5b505050507f7f9c6471d241fdd651ec01d49cde6243143db4bbf23aac6d0573bceb342e674e338260405161025492919061095d565b60405180910390a150565b60004290506000600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dc42c074856040518263ffffffff1660e01b81526004016102c191906109e5565b602060405180830381865afa1580156102de573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610302919061088f565b90508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610373576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161036a90610a72565b60405180910390fd5b600081116103b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103ad90610b04565b60405180910390fd5b600083116103f9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f090610b70565b60405180910390fd5b81600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635b75dd8d86336040518363ffffffff1660e01b8152600401610457929190610b90565b608060405180830381865afa158015610474573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104989190610cc2565b60400151106104dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d390610d3b565b60405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d5764858853384876040518563ffffffff1660e01b815260040161053d9493929190610d5b565b600060405180830381600087803b15801561055757600080fd5b505af115801561056b573d6000803e3d6000fd5b505050507f7e70484266444fe9926cee86f5ca8c91acc579733e6c93081c0cbe7ab3487759338585846040516105a49493929190610da0565b60405180910390a1600060646028346105bd9190610e14565b6105c79190610e9d565b9050600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610631573d6000803e3d6000fd5b508473ffffffffffffffffffffffffffffffffffffffff166108fc82346106589190610ece565b9081150290604051600060405180830381858888f19350505050158015610683573d6000803e3d6000fd5b505050505050565b6000604051905090565b600080fd5b6000819050919050565b6106ad8161069a565b81146106b857600080fd5b50565b6000813590506106ca816106a4565b92915050565b6000602082840312156106e6576106e5610695565b5b60006106f4848285016106bb565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610728826106fd565b9050919050565b6107388161071d565b811461074357600080fd5b50565b6000813590506107558161072f565b92915050565b6000806040838503121561077257610771610695565b5b600061078085828601610746565b9250506020610791858286016106bb565b9150509250929050565b600082825260208201905092915050565b7f55534443207072696365206d75737420626520626967676572207468616e203060008201527f2e00000000000000000000000000000000000000000000000000000000000000602082015250565b600061080860218361079b565b9150610813826107ac565b604082019050919050565b60006020820190508181036000830152610837816107fb565b9050919050565b6000610849826106fd565b9050919050565b6108598161083e565b82525050565b60006020820190506108746000830184610850565b92915050565b600081519050610889816106a4565b92915050565b6000602082840312156108a5576108a4610695565b5b60006108b38482850161087a565b91505092915050565b7f55534443207072696365206d75737420626520646966666572656e742074686160008201527f6e2063757272656e74206f6e652e000000000000000000000000000000000000602082015250565b6000610918602e8361079b565b9150610923826108bc565b604082019050919050565b600060208201905081810360008301526109478161090b565b9050919050565b6109578161069a565b82525050565b60006040820190506109726000830185610850565b61097f602083018461094e565b9392505050565b6000819050919050565b60006109ab6109a66109a1846106fd565b610986565b6106fd565b9050919050565b60006109bd82610990565b9050919050565b60006109cf826109b2565b9050919050565b6109df816109c4565b82525050565b60006020820190506109fa60008301846109d6565b92915050565b7f5375627363726962657220616e64207375627363726970746f72206d7573742060008201527f626520646966666572656e742e00000000000000000000000000000000000000602082015250565b6000610a5c602d8361079b565b9150610a6782610a00565b604082019050919050565b60006020820190508181036000830152610a8b81610a4f565b9050919050565b7f41646472657373206973206e6f7420616363657074696e67207375627363726960008201527f7074696f6e732e00000000000000000000000000000000000000000000000000602082015250565b6000610aee60278361079b565b9150610af982610a92565b604082019050919050565b60006020820190508181036000830152610b1d81610ae1565b9050919050565b7f5765656b73206d7573742062652067726561746572207468616e20302e000000600082015250565b6000610b5a601d8361079b565b9150610b6582610b24565b602082019050919050565b60006020820190508181036000830152610b8981610b4d565b9050919050565b6000604082019050610ba560008301856109d6565b610bb26020830184610850565b9392505050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610c0782610bbe565b810181811067ffffffffffffffff82111715610c2657610c25610bcf565b5b80604052505050565b6000610c3961068b565b9050610c458282610bfe565b919050565b600060808284031215610c6057610c5f610bb9565b5b610c6a6080610c2f565b90506000610c7a8482850161087a565b6000830152506020610c8e8482850161087a565b6020830152506040610ca28482850161087a565b6040830152506060610cb68482850161087a565b60608301525092915050565b600060808284031215610cd857610cd7610695565b5b6000610ce684828501610c4a565b91505092915050565b7f546865726520697320612072756e6e696e6720737562736372697074696f6e00600082015250565b6000610d25601f8361079b565b9150610d3082610cef565b602082019050919050565b60006020820190508181036000830152610d5481610d18565b9050919050565b6000608082019050610d7060008301876109d6565b610d7d6020830186610850565b610d8a604083018561094e565b610d97606083018461094e565b95945050505050565b6000608082019050610db56000830187610850565b610dc260208301866109d6565b610dcf604083018561094e565b610ddc606083018461094e565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610e1f8261069a565b9150610e2a8361069a565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615610e6357610e62610de5565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000610ea88261069a565b9150610eb38361069a565b925082610ec357610ec2610e6e565b5b828204905092915050565b6000610ed98261069a565b9150610ee48361069a565b925082821015610ef757610ef6610de5565b5b82820390509291505056fea264697066735822122028ecc07764f926fb82789e5991415b4e30925152ad9799de782d021ee7e016ad64736f6c634300080b0033",
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
