package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/configor"
	"github.com/toqueteos/ts3"
)

const (
	KeepAliveMins = 5
)

var Config = struct {
	Host              string `default:"127.0.0.1:10011"`
	Whitelisted       bool   `default:"true"`
	QueryUser         string `required:"true"`
	QueryPass         string `required:"true"`
	VirtualServer     uint   `default:"1"`
	BotNickName       string `default:"LoLBot"`
	BotDefaultChannel uint   `default:"0"`

	Region string `required:"true"`
}{}

func main() {
	configor.Load(&Config, "config.json")

	conn, err := ts3.Dial(Config.Host, Config.Whitelisted)
	if err != nil {
		fmt.Printf("Dial error: %s", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	bot(conn)
}

func bot(conn *ts3.Conn) {
	defer conn.Cmd("quit")

	command(conn, fmt.Sprintf("login %s %s", Config.QueryUser, Config.QueryPass), true, false)
	command(conn, fmt.Sprintf("use %d", Config.VirtualServer), true, true)
	clientID := aliveTick(conn)
	command(conn, fmt.Sprintf("clientupdate client_nickname=%s", ts3.Quote(Config.BotNickName)), false, true)
	if Config.BotDefaultChannel > 0 {
		command(conn, fmt.Sprintf("clientmove clid=%d cid=%d", clientID, Config.BotDefaultChannel), false, true)
	}

	conn.NotifyFunc(func(eventType string, data string) {
		if eventType == "notifytextmessage" {
			payload := parseTsDataString(data)
			message := ts3.Unquote(payload["msg"])
			// targetmode invokerid
			if strings.HasPrefix(message, "!") {
				handleCommand(message[1:], func(text string) {
					channelMsg(conn, text)
				})
			}
		}
	})
	command(conn, "servernotifyregister event=textchannel id=0", true, true)

	// Keep
	t := time.NewTicker(time.Minute * KeepAliveMins)
	for {
		<-t.C
		clientID = aliveTick(conn)
	}
}

func aliveTick(conn *ts3.Conn) uint {
	cmdResp, err := conn.Cmd("whoami")
	if err.Id != 0 {
		fmt.Printf("! keepalive err %d: %s\n", err.Id, err.Msg)
		os.Exit(3)
	}
	payload := parseTsDataString(cmdResp)
	clientID, parseErr := strconv.ParseUint(payload["client_id"], 10, 0)
	if parseErr != nil {
		fmt.Printf("! keepalive client_id NaN: %s", parseErr.Error())
		os.Exit(3)
	}
	return uint(clientID)
}

func command(conn *ts3.Conn, cmdReq string, critical bool, verbose bool) {
	cmdResp, err := conn.Cmd(cmdReq)
	hasErr := err.Id != 0
	print := verbose || hasErr
	if print {
		fmt.Printf("> %s\n", cmdReq)
		if len(cmdResp) > 0 {
			fmt.Printf("< %s\n", cmdResp)
		}
	}
	if hasErr {
		fmt.Printf("! err %d: %s\n", err.Id, err.Msg)
		if critical {
			os.Exit(2)
		}
	}
}

func channelMsg(conn *ts3.Conn, msg string) {
	command(conn, fmt.Sprintf("sendtextmessage targetmode=2 msg=%s", ts3.Quote(msg)), false, false)
}
