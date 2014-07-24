package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Schumix/semver"
	"io"
	"net"
	"time"
)

var conn net.Conn
var connectionState = make(chan int)
var state = STATE_CLOSED
var mHost string

func connectToClient(host string) {
	connectionState <- STATE_OPENING

	mHost = host
	fmt.Print("[SOCKET] Connecting to ", host, "...\n")
	var err error

	conn, err = net.Dial("tcp", host)

	go reConnect()

	if err != nil {
		connectionState <- STATE_CLOSED

		fmt.Println(err)
		fmt.Println("[SOCKET] Fail.")
	} else {
		connectionState <- STATE_OPEN
		fmt.Print("[SOCKET] Done. ")

		go regConnection()

		listenToSocket()

		defer conn.Close()
	}
}

func listenToSocket() {
	fmt.Printf("Listening...\n")
	buffer := make([]byte, MAX_BUFFER_SIZE)
	for {
		if state == STATE_CLOSED {
			break
		}
		n, err := conn.Read(buffer[:])
		if err != nil {
			fmt.Println(err)
		}
		if err == io.EOF {
			fmt.Println("[SOCKET] Remote server closed connection.")
			connectionState <- STATE_CLOSED
			break
		}
		handlePacket(buffer[:n], n)
	}
}

func handlePacket(data []byte, size int) {
	// separate packet to its elements
	packet := bytes.Split(data, PACKET_SEPARATOR)
	if packet[0] == "" {
		fmt.Print("Empty packet.")
		return
	}
	fmt.Print("-- START PACKET -- ", size, " bytes")
	fmt.Print(" -- Opcode: ", packet[0], " -- ")
	buf := bytes.NewReader(packet[0])
	opcode, err := binary.ReadVarint(buf)
	if err != nil {
		fmt.Println(err)
	}
	switch opcode {
	case SMSG_AUTH_APPROVED:
		fmt.Println("Auth request approved.")
		requestVersion()
	case SMSG_AUTH_DENIED:
		fmt.Println("Auth request denied.")
	case SMSG_CLOSE_CONNECTION:
		connectionState <- STATE_CLOSING
		fmt.Println("Server sent closing signal. Connection closed.")
		conn.Close()
	case SMSG_PING:
		fmt.Println("SMSG_PING")
		sendPong()
	case SMSG_PONG:
		fmt.Println("SMSG_PONG")
	case SMSG_SCHUMIX_VERSION:
		checkVersion(packet[1])
	default:
		fmt.Println("Unknown opcode.")
	}
	fmt.Println(packet[1:])
	fmt.Println("-- END PACKET --")
}

func reConnect() {
	for which := range connectionState {
		switch which {
		case STATE_OPEN:
			state = STATE_OPEN
		case STATE_OPENING:
		case STATE_CLOSING:
		case STATE_CLOSED:
			state = STATE_CLOSED
			dur, err := time.ParseDuration(config["Timeout"])
			if err != nil {
				fmt.Println(err)
				break
			}
			time.Sleep(dur)
			fmt.Println("[SOCET] Reconnecting...")
			go connectToClient(mHost)
		}
	}
}

func shutdownSocket() {
	if state == STATE_OPEN || state == STATE_OPENING {
		fmt.Println("Shutting down socket connection...")
		sendCloseSignal()
		conn.Close()
	}
}

func sendPing() {
	msg := strconv.Itoa(CMSG_PING) + PACKET_SEPARATOR
	fmt.Fprint(conn, msg)
}

func sendPong() {
	msg := strconv.Itoa(CMSG_PONG) + PACKET_SEPARATOR
	fmt.Fprint(conn, msg)
}

func sendCloseSignal() {
	msg := strconv.Itoa(CMSG_CLOSE_CONNECTION) + PACKET_SEPARATOR +
		"uh. stomachache. shutting down for now." + PACKET_SEPARATOR
	fmt.Fprint(conn, msg)
}

func regConnection() {
	msg := strconv.Itoa(CMSG_REQUEST_AUTH) + PACKET_SEPARATOR +
		"schumix webadmin (reg GUID)" + PACKET_SEPARATOR + md5_gen("schumix") + PACKET_SEPARATOR
	fmt.Fprint(conn, msg)
}

func requestVersion() {
	msg := strconv.Itoa(CMSG_SCHUMIX_VERSION) + PACKET_SEPARATOR
	fmt.Fprint(conn, msg)
}

func checkVersion(ver string) {
	v1, _ := semver.New(MIN_SCHUMIX_VERSION)
	v2, _ := semver.New(ver)

	if v2.Compare(v1) == 0 || v2.Compare(v1) == 1 {
		fmt.Println("Version check OK.")
		fmt.Println("[VERSION] Webadmin:", VERSION, "Min Schumix:",
			MIN_SCHUMIX_VERSION, "Schumix connected:", v2)
	} else {
		fmt.Println("Schumix version is too low...")
		shutdownSocket()
	}
}
