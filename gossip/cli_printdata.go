package main

func (cli *CLI) printData(nodeID string) {
  dataSet := newDataSet(nodeID)
  dataSet.printDataSet()
}
