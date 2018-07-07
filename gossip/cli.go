package main

import (
	"flag"
	"fmt"
	"log"

	"os"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createdata - Create initial data set, require -data")
  fmt.Println("  adddata - Add data to the data set, require -data")
  fmt.Println("  printdata - Print data set.")
	fmt.Println("  startnode - Start a node with ID specified in NODE_ID env. var.")
  fmt.Println("Example:")
  fmt.Println("  ./gossip createdata -data 1")
  fmt.Println("  ./gossip adddata -data asd")
  fmt.Println("  ./gossip printdata")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	createDataCmd := flag.NewFlagSet("createdata", flag.ExitOnError)
	addDataCmd := flag.NewFlagSet("adddata", flag.ExitOnError)
  printDataCmd := flag.NewFlagSet("printdata", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	getDataToCreate := createDataCmd.String("data", "", "The data to create the initial data set")
	getDataToAdd := addDataCmd.String("data", "", "The data to add to the data set")

	switch os.Args[1] {
	case "createdata":
		err := createDataCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "adddata":
		err := addDataCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
  case "printdata":
		err := printDataCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createDataCmd.Parsed() {
		if *getDataToCreate == "" {
			createDataCmd.Usage()
			os.Exit(1)
		}
		cli.createData(*getDataToCreate, nodeID)
	}

	if addDataCmd.Parsed() {
		if *getDataToAdd == "" {
			addDataCmd.Usage()
			os.Exit(1)
		}
		cli.addData(*getDataToAdd, nodeID)
	}

  if printDataCmd.Parsed() {
		cli.printData(nodeID)
	}

	if startNodeCmd.Parsed() {
		cli.startNode(nodeID)
	}
}
