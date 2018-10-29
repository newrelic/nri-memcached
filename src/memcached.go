package main

import (
	//sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	//"github.com/newrelic/infra-integrations-sdk/data/event"
	//"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/mitchellh/mapstructure"
	"github.com/newrelic/infra-integrations-sdk/integration"
	//"github.com/newrelic/infra-integrations-sdk/log"
	//"encoding/binary"
	//"bufio"
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"io"
	"net"
	"regexp"
	//"strconv"
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
	println("==== STATS ====")
	stats := getStats("")
	processGeneralStats(stats, nil)

	println("==== ITEMS ====")
	items := getStats("items")
	processItemStats(items, nil)

	println("==== SLABS ====")
	slabs := getStats("slabs")
	processSlabStats(slabs, nil)
	// TODO close the server connection with quit
}

// TODO get inventory with "stats settings"

func processGeneralStats(stats map[string]string, i *integration.Integration) {
	var s GeneralStats
	config := mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		panic("unreachable")
	}
	decoder.Decode(stats)

	fmt.Printf("%#v\n", s)
}

func processItemStats(stats map[string]string, i *integration.Integration) {
	slabs := partitionItemsBySlabID(stats)
	for slabID, slabMetrics := range slabs {
		var s ItemStats
		config := mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &s,
		}

		decoder, err := mapstructure.NewDecoder(&config)
		if err != nil {
			panic("unreachable")
		}
		decoder.Decode(slabMetrics)

		fmt.Printf("%s, %#v\n", slabID, s)
	}
}

func partitionItemsBySlabID(items map[string]string) map[string]map[string]string {
	pattern := regexp.MustCompile(`items:(\d+):([a-z_]+)`)

	returnMap := make(map[string]map[string]string)
	for key, val := range items {
		matches := pattern.FindStringSubmatch(key)
		slabID := matches[1]
		metricName := matches[2]

		// Retrieve the slab metrics. Create it if it doesn't exist
		slab, ok := returnMap[slabID]
		if !ok {
			slab = make(map[string]string)
			returnMap[slabID] = slab
		}

		slab[metricName] = val
	}

	return returnMap
}

func processSlabStats(stats map[string]string, i *integration.Integration) {
	// TODO do something with general statistics
	slabs, _ := partitionSlabsBySlabID(stats)

	for slabID, slabStats := range slabs {
		var s SlabStats
		config := mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &s,
		}

		decoder, err := mapstructure.NewDecoder(&config)
		if err != nil {
			panic("unreachable")
		}
		decoder.Decode(slabStats)

		fmt.Printf("%s: %#v\n", slabID, s)
	}
}

func partitionSlabsBySlabID(slabs map[string]string) (map[string]map[string]string, map[string]string) {
	slabPattern := regexp.MustCompile(`(\d+):([a-z_]+)`)
	generalPattern := regexp.MustCompile(`^[a-z_]+$`)

	statsMap := make(map[string]map[string]string)
	generalMap := make(map[string]string)

	for key, val := range slabs {
		if generalPattern.MatchString(key) {
			generalMap[key] = val
			continue
		}

		matches := slabPattern.FindStringSubmatch(key)
		slabID := matches[1]
		metricName := matches[2]

		// Retrieve the slab metrics. Create it if it doesn't exist
		slab, ok := statsMap[slabID]
		if !ok {
			slab = make(map[string]string)
			statsMap[slabID] = slab
		}

		slab[metricName] = val
	}

	return statsMap, generalMap

}

// TODO authentication
func getStats(key string) map[string]string {

	keyBytes := []byte(key)
	keyLen := len(keyBytes)

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

	metrics := make(map[string]string)

	for {

		headerBytes := make([]byte, headerSize)
		n, err := io.ReadFull(conn, headerBytes)
		if err != nil {
			println(err.Error())
			break
		}
		if n != 24 {
			println("didn't read 24 bytes")
			break
		}

		header := &ResponseHeader{}
		err = struc.Unpack(bytes.NewReader(headerBytes), header)
		if err != nil {
			println(err.Error())
			break
		}

		if int(header.Status) != 0 || int(header.KeyLength) == 0 {
			// TODO handle error cases more gracefully
			break
		}

		body := make([]byte, header.TotalBodyLength)
		n, err = io.ReadFull(conn, body)
		if err != nil {
			println(err.Error())
			break
		}
		if n != header.TotalBodyLength {
			println("didn't read full body")
		}

		key := string(body[header.ExtrasLength : header.KeyLength+header.ExtrasLength])
		value := string(body[header.KeyLength+header.ExtrasLength:])

		metrics[key] = value
	}

	return metrics
}
