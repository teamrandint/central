package main

import (
	"fmt"
	"seng468/triggerserver/quote"
)

func main() {
	dec, err := quoteclient.Query("username", "AAA", 55)
	fmt.Println(dec, err)
}
