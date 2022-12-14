package utils

import "github.com/jaevor/go-nanoid"

type IDGenerator interface {
	Generate() string
}

type nanoIDGenerator struct {
	gen func() string
}

func NewNanoIDGenerator(length int) IDGenerator {
	gen, err := nanoid.Standard(length)
	if err != nil {
		panic(err)
	}
	return &nanoIDGenerator{
		gen,
	}
}

func (n *nanoIDGenerator) Generate() string {
	return n.gen()
}
