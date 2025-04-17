package euchre

import "fmt"

type card struct {
	denomination denomination
	suit         suit
}

func (c card) String() string {
	return fmt.Sprintf("%s%s", c.denomination, c.suit) //  + " of "
}
