package protocol

import (
	"fmt"
	"github.com/Schumix/semver"
	"net"
	"strconv"
)

type Settings struct {
	Conn net.Conn
}

var settings Settings

func Setup(user_sett *Settings) {
	settings = *user_sett
}

func SendPing() {
	msg := strconv.Itoa(CMSG_PING) + PACKET_SEPARATOR
	fmt.Fprint(settings.Conn, msg)
}

func SendPong() {
	msg := strconv.Itoa(CMSG_PONG) + PACKET_SEPARATOR
	fmt.Fprint(settings.Conn, msg)
}

func SendCloseSignal() {
	msg := strconv.Itoa(CMSG_CLOSE_CONNECTION) + PACKET_SEPARATOR +
		"uh. stomachache. shutting down for now." + PACKET_SEPARATOR
	fmt.Fprint(settings.Conn, msg)
}

func RegConnection() {
	msg := strconv.Itoa(CMSG_REQUEST_AUTH) + PACKET_SEPARATOR +
		"schumix webadmin (reg GUID)" + PACKET_SEPARATOR +
		md5_gen("schumix") + PACKET_SEPARATOR
	fmt.Fprint(settings.Conn, msg)
}

func RequestVersion() {
	msg := strconv.Itoa(CMSG_SCHUMIX_VERSION) + PACKET_SEPARATOR
	fmt.Fprint(settings.Conn, msg)
}

func CheckVersion(ver string) bool {
	v1, _ := semver.New(MIN_SCHUMIX_VERSION)
	v2, _ := semver.New(ver)

	if v2.Compare(v1) == 0 || v2.Compare(v1) == 1 {
		fmt.Println("Version check OK.")
		fmt.Println("[VERSION] Webadmin:", VERSION, "Min Schumix:",
			MIN_SCHUMIX_VERSION, "Schumix connected:", v2)
		return true
	} else {
		fmt.Println("Schumix version is too low...")
	}
	return false
}
