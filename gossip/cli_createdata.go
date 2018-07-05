package main

import (
	"fmt"
)

func (cli *CLI) createData(data, nodeID string) {
  createDataSet(data, nodeID)
	fmt.Println("Done!")
}
