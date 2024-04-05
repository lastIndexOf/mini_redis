package reply

import (
	"bytes"
	"strconv"

	"github.com/lastIndexOf/mini_redis/interface/resp"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

type ErrorReply interface {
	Error() string
	resp.Reply
}

func IsErrReply(r resp.Reply) bool {
	return r.Bytes()[0] == '-'
}

type StandardErrReply struct {
	Status string
}

func (r *StandardErrReply) Bytes() []byte {
	return []byte("-" + r.Status + CRLF)
}

func MakeStandardErrReply(status string) *StandardErrReply {
	return &StandardErrReply{status}
}

type BulkReply struct {
	Arg []byte
}

func (r *BulkReply) Bytes() []byte {
	if len(r.Arg) == 0 {
		return nullBulkReplyBytes
	}

	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

type MultiBulkReply struct {
	Args [][]byte
}

func (r *MultiBulkReply) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(len(r.Args)) + CRLF)

	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}

	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

type StatusReply struct {
	Status string
}

func (r *StatusReply) Bytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{Status: status}
}

type IntReplay struct {
	Code int64
}

func (r *IntReplay) Bytes() []byte {
	return []byte(":" + strconv.Itoa(int(r.Code)) + CRLF)
}

func MakeIntReplay(code int64) *IntReplay {
	return &IntReplay{code}
}
