package generic

type OrderingArg struct {
	Field      string
	Descending bool
}

type OrderingArgs []OrderingArg
