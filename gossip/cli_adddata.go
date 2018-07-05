package main

import (
  "fmt"
)

func (cli *CLI) addData(data, nodeID string) {
  dataSet := newDataSet(nodeID)
  dataSet.addData(data, nodeID)
  fmt.Println("Done")
}
