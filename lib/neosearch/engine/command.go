package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

// Command defines a NeoSearch internal command.
// This command describes a single operation in the index storage and is
// decomposed in the following parts:
//   - Index
//   - Database
//   - Key
//   - KeyType
//   - Value
//   - ValueType
//   - Batch
type Command struct {
	Index     string
	Database  string
	Command   string
	Key       []byte
	KeyType   uint8
	Value     []byte
	ValueType uint8

	Batch bool
}

func (c Command) Println() {
	line := c.Reverse()
	fmt.Println(line)
}

func (c Command) Reverse() string {
	var (
		keyStr string
		valStr string
		line   string
	)

	if c.Key != nil {
		if c.KeyType == TypeString {
			keyStr = `'` + string(c.Key) + `'`
		} else if c.KeyType == TypeUint {
			keyStr = `uint(` + strconv.Itoa(int(utils.BytesToUint64(c.Key))) + `)`
		} else if c.KeyType == TypeInt {
			keyStr = `int(` + strconv.Itoa(int(utils.BytesToInt64(c.Key))) + `)`
		} else if c.KeyType == TypeFloat {
			keyStr = `float(` + strconv.FormatFloat(utils.BytesToFloat64(c.Key), 'f', -1, 64) + `)`
		} else if c.KeyType == TypeBool {
			keyStr = `bool(` + string(c.Key) + `)`
		} else {
			fmt.Printf("Command error: %+v", c)
			panic(fmt.Errorf("Invalid command key type: %d - %+v", c.KeyType, string(c.Key)))
		}
	}

	if c.Value != nil {
		if c.ValueType == TypeString {
			valStr = `'` + string(c.Value) + `'`
		} else if c.ValueType == TypeUint {
			valStr = `uint(` + strconv.Itoa(int(utils.BytesToUint64(c.Value))) + `)`
		} else if c.ValueType == TypeInt {
			valStr = `int(` + strconv.Itoa(int(utils.BytesToInt64(c.Value))) + `)`
		} else if c.ValueType == TypeFloat {
			valStr = `float(` + strconv.FormatFloat(utils.BytesToFloat64(c.Value), 'f', -1, 64) + `)`
		} else if c.ValueType == TypeBool {
			valStr = `bool(` + string(c.Value) + `)`
		} else {
			panic(fmt.Errorf("Invalid command key type: %d", c.ValueType))
		}
	}

	switch strings.ToUpper(c.Command) {
	case "SET", "MERGESET":
		line = fmt.Sprintf("USING %s.%s %s %s %s;", c.Index, c.Database, strings.ToUpper(c.Command), keyStr, valStr)
	case "BATCH", "flushbatch":
		line = fmt.Sprintf("USING %s.%s %s;", c.Index, c.Database, strings.ToUpper(c.Command))
	case "GET", "DELETE":
		line = fmt.Sprintf("USING %s.%s %s %s;", c.Index, c.Database, strings.ToUpper(c.Command), keyStr)
	default:
		panic(fmt.Errorf("Invalid command: %s: %v", strings.ToUpper(c.Command), c))
	}

	return line
}
