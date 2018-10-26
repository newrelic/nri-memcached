package main

import (
	//sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	//"github.com/newrelic/infra-integrations-sdk/data/event"
	//"github.com/newrelic/infra-integrations-sdk/data/metric"
	//"github.com/newrelic/infra-integrations-sdk/integration"
	//"encoding/binary"
	"bufio"
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"net"
)

// More info about the binary protocol here: https://github.com/memcached/memcached/wiki/BinaryProtocolRevamped

// RequestHeader is a marshalling struct for the memcached binary header format
type RequestHeader struct {
	Magic           byte
	Opcode          byte
	KeyLength       int `struc:"int16"`
	ExtrasLength    int `struc:"int8"`
	DataType        byte
	VBucketID       int `struc:"int16"`
	TotalBodyLength int `struc:"int32"`
	Opaque          [4]byte
	CAS             [8]byte
}

// ResponseHeader is a marshalling struct for the memcached binary header format
type ResponseHeader struct {
	Magic           byte
	Opcode          byte
	KeyLength       int `struc:"int16"`
	ExtrasLength    int `struc:"int8"`
	DataType        byte
	Status          int `struc:"int16"`
	TotalBodyLength int `struc:"int32"`
	Opaque          [4]byte
	CAS             [8]byte
}

const (
	headerSize = 24
)

func main() {
	items := "items"
	slabs := "slabs"
	getStats(nil)
	getStats(&items)
	getStats(&slabs)
}

// TODO authentication
func getStats(key *string) {

	var keyBytes []byte
	var keyLen int
	if key != nil {
		keyBytes = []byte(*key)
		keyLen = len(keyBytes)
	} else {
		keyBytes = []byte{}
		keyLen = 0
	}

	sendHeader := &RequestHeader{
		Magic:           0x80,
		Opcode:          0x10,
		KeyLength:       keyLen,
		ExtrasLength:    0x00,
		DataType:        0x00,
		VBucketID:       0x0000,
		TotalBodyLength: keyLen,
		Opaque:          [4]byte{0x01, 0x01, 0x01, 0x01},
		CAS:             [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	var sendBuf bytes.Buffer

	err := struc.Pack(&sendBuf, sendHeader)
	if err != nil {
		println(err.Error())
	}
	_, err = sendBuf.Write(keyBytes)
	if err != nil {
		println(err.Error())
	}

	conn, err := net.Dial("tcp", "mc-14-1:11211")
	if err != nil {
		println(err.Error())
	}

	_, err = conn.Write(sendBuf.Bytes())
	if err != nil {
		println(err.Error())
	}

	// TODO investigate reading into the buffer dynamically
	responseBuf := bufio.NewReaderSize(conn, 8000)

	for {

		headerBytes := make([]byte, headerSize)
		_, err := responseBuf.Read(headerBytes)
		if err != nil {
			println(err.Error())
		}

		header := &ResponseHeader{}
		err = struc.Unpack(bytes.NewReader(headerBytes), header)
		if err != nil {
			println(err.Error())
		}

		if int(header.Status) != 0 || int(header.KeyLength) == 0 {
			// TODO handle error cases more gracefully
			break
		}

		body := make([]byte, header.TotalBodyLength)
		_, err = responseBuf.Read(body)
		if err != nil {
			println(err.Error())
		}

		key := body[header.ExtrasLength:header.KeyLength]
		value := body[header.KeyLength+header.ExtrasLength:]

		fmt.Printf("%s: %s\n", string(key), string(value))
	}
}
