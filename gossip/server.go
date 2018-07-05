package main

import (
  "bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
  "os"
  "bufio"
  "strconv"
  "github.com/boltdb/bolt"
)

var nodeAddress string
var neighborListFile = "neighbors_%s.dat"

const seedNode = "localhost:3000"
const protocol = "tcp"
const commandLength = 12

// StartServer starts a node
func StartServer(nodeID string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
  neighborListFile = fmt.Sprintf(neighborListFile, nodeID)

	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	ds := newDataSet(nodeID)

  if nodeAddress != seedNode {
    if _, err := os.Stat(neighborListFile); err != nil {
      createNeighborList()
    }

    sendGetaddr(seedNode)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, ds)
	}
}

func createNeighborList() {
  file, err := os.Create(neighborListFile)
  if err != nil {
    log.Panic(err)
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  fmt.Fprintln(w, seedNode)
  w.Flush()
}

func addToNeighborList(addr string) {
  file, err := os.OpenFile(neighborListFile, os.O_APPEND, 0666)
  if err != nil {
    log.Panic(err)
  }
  defer file.Close()

  nodeList := getNeighbors()

  if len(nodeList) >= 2 {
    return
  }

  for _, node := range nodeList {
    if node == addr || addr == nodeAddress {
      return
    }
  }

  w := bufio.NewWriter(file)
  fmt.Fprintln(w, addr)
  w.Flush()
}

func getNeighbors() []string {
  file, err := os.Open(neighborListFile)
  if err != nil {
    log.Panic(err)
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines
}

func handleConnection(conn net.Conn, ds *DataSet) {
  request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
  case "getaddr":
    handleGetaddr(request, ds)
  case "addr":
    handleAddr(request, ds)
  case "version":
    handleVersion(request, ds)
  case "getdata":
    handleGetdata(request, ds)
  case "data":
    handleData(request, ds)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

type addr struct {
	AddrList []string
}

type getaddr struct {
	AddrFrom string
}

type version struct {
	Height     int
	AddrFrom   string
}

type getdata struct {
	AddrFrom string
}

type data struct {
  Data []string
}

func sendGetaddr(address string) {
  data := getaddr{nodeAddress}
  payload := gobEncode(data)
  request := append(commandToBytes("getaddr"), payload...)

  sendTo(address, request)
}

func handleGetaddr(request []byte, ds *DataSet) {
  var buff bytes.Buffer
	var payload getaddr

  buff.Write(request[commandLength:])
  dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

  foreignerAddr := payload.AddrFrom
  addToNeighborList(foreignerAddr)

  sendAddr(foreignerAddr)
}

func sendAddr(address string) {
  neighbors := getNeighbors()
  data := addr{neighbors}
  payload := gobEncode(data)
  request := append(commandToBytes("addr"), payload...)

  sendTo(address, request)
}

func handleAddr(request []byte, ds *DataSet) {
  var buff bytes.Buffer
  var payload addr

  buff.Write(request[commandLength:])
  dec := gob.NewDecoder(&buff)
  err := dec.Decode(&payload)
  if err != nil {
    log.Panic(err)
  }

  addresses := payload.AddrList

  for _, address := range addresses {
    addToNeighborList(address)
  }

  neighbors := getNeighbors()

  for _, neighbor := range neighbors {
    sendVersion(neighbor, ds)
  }
}

func sendVersion(addr string, ds *DataSet) {
  height := ds.getHeight()
  data := version{height, nodeAddress}
  payload := gobEncode(data)
  request := append(commandToBytes("version"), payload...)

  sendTo(addr, request)
}

func handleVersion(request []byte, ds *DataSet) {
  var buff bytes.Buffer
	var payload version

  buff.Write(request[commandLength:])
  dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

  foreignerHeight := payload.Height
  foreignerAddr := payload.AddrFrom
  localHeight := ds.getHeight()

  if localHeight > foreignerHeight {
    sendVersion(foreignerAddr, ds)
  } else if localHeight < foreignerHeight {
    sendGetdata(foreignerAddr)
  }
}

func sendGetdata(addr string) {
  data := getdata{nodeAddress}
  payload := gobEncode(data)
  request := append(commandToBytes("getdata"), payload...)

  sendTo(addr, request)
}

func handleGetdata(request []byte, ds *DataSet) {
  var buff bytes.Buffer
  var payload getdata

  buff.Write(request[commandLength:])
  dec := gob.NewDecoder(&buff)
  err := dec.Decode(&payload)
  if err != nil {
    log.Panic(err)
  }

  foreignerAddr := payload.AddrFrom

  sendData(foreignerAddr, ds)
}

func sendData(addr string, ds *DataSet) {
  var d []string

  ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
    h := b.Get([]byte("h"))
    height, _ := strconv.Atoi(string(h))

    for i := 0; i < height; i++ {
      _d := b.Get([]byte(fmt.Sprintf("%v", i)))
      d = append(d, string(_d))
    }

		return nil
	})

  dataToSend := data{d}
  payload := gobEncode(dataToSend)
  request := append(commandToBytes("data"), payload...)

  sendTo(addr, request)
}

func handleData(request []byte, ds *DataSet) {
  var buff bytes.Buffer
  var payload data

  buff.Write(request[commandLength:])
  dec := gob.NewDecoder(&buff)
  err := dec.Decode(&payload)
  if err != nil {
    log.Panic(err)
  }

  newData := payload.Data

  ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

    for index, newDatum := range newData {
      b.Put([]byte(fmt.Sprintf("%v", index)), []byte(newDatum))
    }

    b.Put([]byte("h"), []byte(fmt.Sprintf("%v", len(newData))))

		return nil
	})
}

func sendTo(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		// var updatedNodes []string
    //
		// for _, node := range knownNodes {
		// 	if node != addr {
		// 		updatedNodes = append(updatedNodes, node)
		// 	}
		// }
    //
		// knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}
