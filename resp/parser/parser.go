package parser

import (
	"bufio"
	"errors"
	"github.com/lastIndexOf/mini_redis/lib/logger"
	"github.com/lastIndexOf/mini_redis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"

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
	bulkLen           int64 // for bulk string
}

func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && s.expectedArgsCount == len(s.args)
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)

	go parse0(reader, ch)

	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()

	bufReader := bufio.NewReader(reader)
	var state readState
	for {
		line, isIoErr, err := readLine(bufReader, &state)

		if err != nil {
			if isIoErr {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}

			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}

		if !state.multiLine {
			switch line[0] {
			case '*':
				err := parseMultiBulkHeader(line, &state)

				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(line)),
					}
					state = readState{}
					continue
				}

				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: reply.MakeEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			case '$':
				err := parseBulkHeader(line, &state)

				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(line)),
					}
					state = readState{}
					continue
				}

				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: &reply.NullBulkReply{},
					}
					state = readState{}
					continue
				}
			default:
				res, err := parseSingleLineReply(line)

				ch <- &Payload{
					Data: res,
					Err:  err,
				}

				state = readState{}
				continue
			}
		} else {
			err := readBody(line, &state)

			if err != nil {
				ch <- &Payload{
					Err: err,
				}
				state = readState{}
				continue
			}

			if state.finished() {
				switch state.msgType {
				case '*':
					ch <- &Payload{
						Data: reply.MakeMultiBulkReply(state.args),
						Err:  err,
					}
				case '$':
					ch <- &Payload{
						Data: reply.MakeBulkReply(state.args[0]),
						Err:  err,
					}
				}

				state = readState{bulkLen: 0}
			}
		}
	}
}

// ParseStream parses RESP stream and sends parsed payloads to the channel.
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

// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// *3\r\n
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

// $5\r\nvalue\r\n
func parseBulkHeader(line []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(line[1:len(line)-2]), 10, 64)

	if err != nil {
		return errors.New("protocol error: invalid bulk length (" + err.Error() + ")")
	}

	if state.bulkLen == -1 {
		return nil
	}

	if state.bulkLen <= 0 {
		return errors.New("protocol error: invalid bulk length (" + string(line) + ")")
	}

	state.msgType = line[0]
	state.multiLine = true
	state.expectedArgsCount = 1
	state.args = make([][]byte, 0, 1)

	return nil
}

// +OK\r\n
// -ERR\r\n
// :5\r\n
func parseSingleLineReply(line []byte) (resp.Reply, error) {
	msg := strings.TrimSuffix(string(line), "\r\n")

	var ret resp.Reply
	switch msg[0] {
	case '+':
		ret = reply.MakeStatusReply(msg[1:])
	case '-':
		ret = reply.MakeStandardErrReply(msg[1:])
	case ':':
		val, err := strconv.ParseInt(msg[1:], 10, 64)

		if err != nil {
			return nil, errors.New("protocol error: invalid integer (" + err.Error() + ")")
		}

		ret = reply.MakeIntReply(val)
	}

	return ret, nil
}

// $3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// $3\r\n
// nvalue\r\n
func readBody(line []byte, state *readState) error {
	msg := line[:len(line)-2]

	var err error
	switch msg[0] {
	case '$':
		state.bulkLen, err = strconv.ParseInt(string(msg[1:]), 10, 64)

		if err != nil {
			return errors.New("protocol error: invalid bulk length (" + err.Error() + ")")
		}

		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	default:
		state.args = append(state.args, msg)
	}

	return nil
}
