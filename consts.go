package protocol

const MAX_BUFFER_SIZE = 262144 // 8^6 byte = 256 MB
var PACKET_SEPARATOR = []byte("|;|")

const (
	SCMSG_PACKET_NULL     = 0x0
	CMSG_REQUEST_AUTH     = 0x01
	SMSG_AUTH_APPROVED    = 0x02
	SMSG_AUTH_DENIED      = 0x03
	CMSG_CLOSE_CONNECTION = 0x04
	SMSG_CLOSE_CONNECTION = 0x05
	CMSG_PING             = 0x06
	SMSG_PING             = 0x07
	CMSG_PONG             = 0x08
	SMSG_PONG             = 0x09
	CMSG_SCHUMIX_VERSION  = 0x10
	SMSG_SCHUMIX_VERSION  = 0x11
)

// connection states
const (
	STATE_OPENING = iota
	STATE_OPEN
	STATE_CLOSING
	STATE_CLOSED
)
