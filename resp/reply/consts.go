package reply

var (
	pongBytes           = []byte("+PONG\r\n")
	okBytes             = []byte("+OK\r\n")
	nullBulkBytes       = []byte("$-1\r\n")
	emptyMultiBulkBytes = []byte("*0\r\n")
	noReply             = []byte("")
)

type PongReply struct{}

func (r *PongReply) Bytes() (_ []byte) {
	return pongBytes
}

var thePongReply = new(PongReply)

func MakePongReply() *PongReply {
	return thePongReply
}

type OkReply struct{}

func (r *OkReply) Bytes() (_ []byte) {
	return okBytes
}

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

type NullBulkReply struct{}

func (r *NullBulkReply) Bytes() (_ []byte) {
	return nullBulkBytes
}

var theNullBulkReplay = new(NullBulkReply)

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReplay
}

type EmptyMultiBulkReply struct{}

func (r *EmptyMultiBulkReply) Bytes() (_ []byte) {
	return emptyMultiBulkBytes
}

var theEmptyMultiBulkReplay = new(EmptyMultiBulkReply)

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReplay
}

type NoReply struct{}

func (r *NoReply) Bytes() (_ []byte) {
	return noReply
}

var theNoReply = new(NoReply)

func MakeNoReply() *NoReply {
	return theNoReply
}
