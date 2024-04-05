package reply

var (
	unknownErrBytes   = []byte("-Err unknown\r\n")
	syntaxErrBytes    = []byte("-Err syntax error\r\n")
	wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
)

type UnknownErrReply struct{}

func (r *UnknownErrReply) Error() (_ string) {
	return "Err unknown"
}

func (r *UnknownErrReply) Bytes() (_ []byte) {
	return unknownErrBytes
}

var theUnknownErrReply = &UnknownErrReply{}

func MakeUnknownErrReply() *UnknownErrReply {
	return theUnknownErrReply
}

type ArgNumErrReply struct {
	Cmd string
}

func (r *ArgNumErrReply) Error() (_ string) {
	return "-Err wrong number of arguments for '" + r.Cmd + "' command\r\n"
}

func (r *ArgNumErrReply) Bytes() (_ []byte) {
	return []byte("-Err wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{Cmd: cmd}
}

type SyntaxErrReply struct{}

func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (r *SyntaxErrReply) Bytes() []byte {
	return syntaxErrBytes
}

var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

type WrongTypeErrReply struct{}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (r *WrongTypeErrReply) Bytes() []byte {
	return wrongTypeErrBytes
}

type ProtocolErrReply struct {
	Msg string
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error '" + r.Msg + "' command"
}

func (r *ProtocolErrReply) Bytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}
