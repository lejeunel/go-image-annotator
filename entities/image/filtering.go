package image

type Order int

type FilterQueryStr = string
type OrderingStr = string

const (
	AscOrder Order = iota
	DescOrder
)

type OrderingArg struct {
	Field string
	Order
}

type OrderingArgs []OrderingArg

type Filtering struct {
	Collection *string
}
