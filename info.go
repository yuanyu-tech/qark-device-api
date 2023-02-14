package device

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	logging "github.com/ipfs/go-log/v2"
)

type ID string

const EmptyID = ID("0")

var EmptyInfo = &Info{
	Key: "empty key",
	Uid: "empty uid",
	Id:  ID("0"),
}

type Info struct {
	Key string
	Uid string
	Id  ID
}

func (i Info) String() string {
	return fmt.Sprintf("{Uid:%s Id:%s}", i.Uid, i.Id)
}

var deviceInfo *Info

var log = logging.Logger("device")

func GetGlobalInfo() (*Info, error) {
	// check and initialize the DeviceID
	if deviceInfo == nil {
		var err error
		deviceInfo, err = GetInfo()
		if err != nil {
			return nil, fmt.Errorf("GetGlobalInfo: %w", err)
		}

		log.Infof("Global DeviceInfo were init to %+v", deviceInfo)
	}

	return deviceInfo, nil
}

// GetInfo Get the device info from qark-client node protocol json file
func GetInfo() (*Info, error) {
	var jsonFile = "/var/run/qark-client/node.protocol.json"
	f, err := os.OpenFile(jsonFile, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to get device info, please startup the device service first")
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to load device info with error: %w", err)
	}

	var gInfo Info
	err = json.Unmarshal(data, &gInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal device info info: %w", err)
	}

	return &gInfo, nil
}
