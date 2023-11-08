package parser

type Parser interface {
	Parse(buf []byte) error
	AddVar(name string, val string)
}
