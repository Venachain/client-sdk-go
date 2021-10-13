package types

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions) together.
type Body struct {
	Transactions []*Transaction
}
