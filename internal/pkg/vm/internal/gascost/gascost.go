package gascost

import (
	"fmt"
	"github.com/filecoin-project/go-filecoin/internal/pkg/vm/gas"
	"github.com/filecoin-project/go-filecoin/internal/pkg/vm/internal/message"
	"github.com/filecoin-project/specs-actors/actors/abi"
)

// Pricelist provides prices for operations in the VM.
//
// Note: this interface should be APPEND ONLY since last chain checkpoint
type Pricelist interface {
	// OnChainMessage returns the gas used for storing a message of a given size in the chain.
	OnChainMessage(msgSize int) gas.Unit
	// OnChainReturnValue returns the gas used for storing the response of a message in the chain.
	OnChainReturnValue(receipt *message.Receipt) gas.Unit

	// OnMethodInvocation returns the gas used when invoking a method.
	OnMethodInvocation(value abi.TokenAmount, methodNum abi.MethodNum) gas.Unit

	// OnIpldGet returns the gas used for storing an object
	OnIpldGet(dataSize int) gas.Unit
	// OnIpldPut returns the gas used for storing an object
	OnIpldPut(dataSize int) gas.Unit

	// OnCreateActor returns the gas used for creating an actor
	OnCreateActor() gas.Unit
	// OnDeleteActor returns the gas used for deleting an actor
	OnDeleteActor() gas.Unit
}

var prices = map[abi.ChainEpoch]Pricelist{
	abi.ChainEpoch(0): &pricelistV0{
		onChainMessageBase:        gas.Zero,
		onChainMessagePerByte:     gas.NewGas(2),
		onChainReturnValuePerByte: gas.NewGas(8),
		sendBase:                  gas.NewGas(5),
		sendTransferFunds:         gas.NewGas(5),
		sendInvokeMethod:          gas.NewGas(10),
		ipldGetBase:               gas.NewGas(10),
		ipldGetPerByte:            gas.NewGas(1),
		ipldPutBase:               gas.NewGas(20),
		ipldPutPerByte:            gas.NewGas(2),
		createActorBase:           gas.NewGas(40), // IPLD put + 20
		createActorExtra:          gas.NewGas(500),
		deleteActor:               gas.NewGas(-500), // -createActorExtra
	},
}

// PricelistByEpoch finds the latest prices for the given epoch
func PricelistByEpoch(epoch abi.ChainEpoch) Pricelist {
	// since we are storing the prices as map or epoch to price
	// we need to get the price with the highest epoch that is lower or equal to the `epoch` arg
	bestEpoch := abi.ChainEpoch(0)
	bestPrice := prices[bestEpoch]
	for e, pl := range prices {
		// if `e` happened after `bestEpoch` and `e` is earlier or equal to the target `epoch`
		if e > bestEpoch && e <= epoch {
			bestEpoch = e
			bestPrice = pl
		}
	}
	if bestPrice == nil {
		panic(fmt.Sprintf("bad setup: no gas prices available for epoch %d", epoch))
	}
	return bestPrice
}
