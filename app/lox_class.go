package main

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

func (c *LoxClass) String() string {
	return c.name
}
