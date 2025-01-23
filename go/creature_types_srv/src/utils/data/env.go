package data

type Env uint32

//go:generate stringer -type Env
const (
	Local Env = iota
	Dev
	Prod
)
