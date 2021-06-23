package util

type ComputedValue struct {
	Value      float64
	Percentage float64
}

type KeyValueSelector struct {
	Key string
	Operator string // can be "=" or "!="
	Value string
}
