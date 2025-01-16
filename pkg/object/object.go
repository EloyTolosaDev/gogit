package object

type Object interface {
	Hash() string
	Info() string // used to write tree objects to disk
}
