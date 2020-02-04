package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// MessageFeeBinder is an interface implemented by keepers which contains the method "Messages", which returns a string
// slice containing all the messages supported by said keeper.
type MessageFeeBinder interface {
	Messages() []MessageFeeBinding
}

// MessageFeeBinding represent a constant fee to be required by x/ante when handling messages named after Name.
type MessageFeeBinding struct {
	Name string
	Fee  sdk.Dec
}

// StandardFIATFee is the fixed 0.01€ fee Commercio.network requires for each transaction.
// Each Keeper can define a different fee for each message, although most of Commercio.network modules
// will use StandardFIATFee.
var StandardFIATFee = sdk.NewDecWithPrec(1, 2)

// NewStandardBinding returns a new MessageFeeBinding initialized with StandardFIATFee and messageName.
func NewStandardBinding(messageName string) MessageFeeBinding {
	return MessageFeeBinding{
		Name: messageName,
		Fee:  StandardFIATFee,
	}
}
