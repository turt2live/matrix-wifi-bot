package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"turt2live/matrix-wifi-bot/config"
	"github.com/turt2live/matrix-wifi-bot/logging"
	"matrix-wifi-bot/matrix"
	"strings"
	"os/exec"
	"regexp"
	"time"
)

func main() {
	configPath := flag.String("config", "wifi-bot.yaml", "The path to the configuration")
	flag.Parse()

	config.Path = *configPath

	err := logging.Setup(config.Get().Logging.Directory)
	if err != nil {
		panic(err)
	}

	logrus.Info("Starting monitor bot...")
	client, err := matrix.NewClient(config.Get().Homeserver.Url, config.Get().Homeserver.AccessToken)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Authenticated as ", client.UserID)

	roomId := config.Get().Wifi.AnnounceRoomId
	err = client.JoinRoom(roomId)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Joined announcement room: ", config.Get().Wifi.AnnounceRoomId)

	knownNetworks := make([]string, 0)

	for true {
		nets, err := ScanNetworks()
		if err != nil {
			logrus.Error(err)
			client.SendNotice(roomId, "error", "There was a problem scanning for wifi networks")
		} else {
			for _, ssid := range nets {
				idx := -1
				for i := range knownNetworks {
					if knownNetworks[i] == ssid {
						idx = i
						break
					}
				}

				if idx == -1 {
					logrus.Info("Discovered network: " + ssid)
					client.SendMessage(roomId, "New network discovered: "+ssid)
				}
			}

			for _, ssid := range knownNetworks {
				idx := -1
				for i := range nets {
					if nets[i] == ssid {
						idx = i
						break
					}
				}

				if idx == -1 {
					logrus.Info("Lost network: " + ssid)
					client.SendMessage(roomId, "Network lost: "+ssid)
				}
			}

			knownNetworks = nets
		}
		time.Sleep(60 * time.Second)
	}
}

func ScanNetworks() ([]string, error) {
	logrus.Info("Scanning for wifi networks...")
	cmdFields := strings.Fields(config.Get().Wifi.Command)
	cmd := exec.Command(cmdFields[0], cmdFields[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	result := string(out)

	r := regexp.MustCompile(config.Get().Wifi.SsidSearch)
	matches := r.FindAllStringSubmatch(result, -1)
	ssids := make([]string, 0)
	for i := range matches {
		ssids = append(ssids, strings.TrimSpace(matches[i][1]))
	}
	return ssids, nil
}
