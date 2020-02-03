package types

// MessageLister is an interface implemented by keepers which contains the method "Messages", which returns a string
// slice containing all the messages supported by said keeper.
type MessageLister interface {
	Messages() []string
}
