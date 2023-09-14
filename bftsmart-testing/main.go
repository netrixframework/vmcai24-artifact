package main

import "github.com/netrixframework/bftsmart-testing/cmd"

func main() {
	c := cmd.RootCmd()
	c.Execute()
}
