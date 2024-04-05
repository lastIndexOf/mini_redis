package parser

import (
	"bufio"
	"errors"
	"io"
	"strconv"

	"github.com/lastIndexOf/mini_redis/interface/resp"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	multiLine         bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           uint64 // for bulk string
}

func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && s.expectedArgsCount == len(s.args)
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)

	go parse0(reader, ch)

	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {}

func readLine(reader *bufio.Reader, state *readState) ([]byte, bool, error) {

	if state.bulkLen == 0 {
		// no ($num) prefix
		line, err := reader.ReadBytes('\n')

		if err != nil {
			return nil, true, err
		}

		if len(line) < 2 || line[len(line)-2] != '\r' {
			return nil, false, errors.New("protocol error: invalid line ending (" + string(line) + ")")
		}

		return line, false, err
	}

	line := make([]byte, state.bulkLen+2)
	_, err := io.ReadFull(reader, line)

	if err != nil {
		return nil, true, err
	}

	if len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
		return nil, false, errors.New("protocol error: invalid line ending (" + string(line) + ")")
	}

	state.bulkLen = 0

	return line, false, err
}

func parseMultiBulkHeader(line []byte, state *readState) error {
	expectedLen, err := strconv.ParseInt(string(line[1:len(line)-2]), 10, 64)

	if err != nil {
		return errors.New("protocol error: invalid multibulk length (" + err.Error() + ")")
	}

	if expectedLen < 0 {
		return errors.New("protocol error: invalid multibulk length (" + string(line) + ")")
	}

	if expectedLen == 0 {
		state.expectedArgsCount = 0
		return nil
	}

	state.expectedArgsCount = int(expectedLen)
	state.args = make([][]byte, expectedLen)
	state.msgType = line[0]
	state.multiLine = true

	return nil
}
