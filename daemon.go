package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/xerrors"
)

type Daemon struct {
	Info            *Info
	Duration        int
	MaxOfflineCount int
	MaxOnlineCount  int
}

type HBInfo struct {
	TvSec  int    `json:"tv_sec"`
	TvUSec int    `json:"tv_usec"`
	Pid    int    `json:"pid"`
	PPid   int    `json:"ppid"`
	Sign   string `json:"sign"`
}

func (i HBInfo) String() string {
	return fmt.Sprintf("{Sec:%d USec:%d Pid:%d PPid:%d}", i.TvSec, i.TvUSec, i.Pid, i.PPid)
}

func NewDaemon(info *Info, duration int, maxOfflineCount int, maxOnlineCount int) *Daemon {
	return &Daemon{
		Info:            info,
		Duration:        duration,
		MaxOfflineCount: maxOfflineCount,
		MaxOnlineCount:  maxOnlineCount,
	}
}

func (d *Daemon) GetInfo() (HBInfo, error) {
	url := "http://localhost:17181/api/node/heartbeat"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HBInfo{}, xerrors.Errorf("http.NewRequest error: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return HBInfo{}, xerrors.Errorf("request %s error: %w", url, err)
	}

	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		return HBInfo{}, xerrors.Errorf("(%s), request failed: StatusCode=%s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return HBInfo{}, xerrors.Errorf("(%s), read all error: %w", err)
	}

	// {"errno":200,"errmsg":"","data":{"tv_sec":1630720889,"tv_usec":674877,"pid":35019,"ppid":35013,"sign":"196e73d13a34d142490a44e5f9ebb88574765718"}}
	m := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &m); err != nil {
		return HBInfo{}, xerrors.Errorf("json unmarshal error: %w", err)
	}

	errCode, has1 := m["errno"]
	dataStr, has2 := m["data"]
	if has1 == false || has2 == false {
		return HBInfo{}, xerrors.Errorf("invalid data package, server error maybe: %+v", m)
	}

	errno := GetInt(errCode)
	if errno != 200 {
		return HBInfo{}, xerrors.Errorf("none-200 http response, errno=%d, errstr=%+v", errno, dataStr)
	}

	dMap := dataStr.(map[string]interface{})
	var info = HBInfo{
		TvSec:  GetInt(dMap["tv_sec"]),
		TvUSec: GetInt(dMap["tv_usec"]),
		Pid:    GetInt(dMap["pid"]),
		PPid:   GetInt(dMap["ppid"]),
		Sign:   GetString(dMap["sign"]),
	}

	// check the sign to make sure the
	// response bytes are from the qark-client api
	seed := fmt.Sprintf("%d", info.TvSec)
	vSign, err := BuildSign(d.Info.Key, string(d.Info.Id), seed, seed, "__heartbeat__")
	if err != nil {
		return HBInfo{}, xerrors.Errorf("failed build sign: %w", err)
	}

	if vSign != info.Sign {
		return HBInfo{}, xerrors.Errorf("sign check error, make sure the qark-client is now running")
	}

	return info, nil
}

// Monitor Daemon Coroutine checking
func (d *Daemon) Monitor(ctx context.Context, offlineCallback func(count int, daemon *Daemon), onlineCallback func(count int, daemon *Daemon)) {
	log.Infof("qark-client monitor started")
	var offlineCounter int
	var onlineCounter int
	var duration = time.Duration(d.Duration) * time.Second
	for {
		if _, err := d.GetInfo(); err != nil {
			offlineCounter++
			onlineCounter = 0
			log.Errorf("offline count %d, offlineCallback if (count >= %d), err: %+v", offlineCounter, d.MaxOfflineCount, err)
		} else {
			onlineCounter++
			offlineCounter = 0
		}

		// check and invoke the callback function
		if offlineCounter >= d.MaxOfflineCount {
			offlineCallback(offlineCounter, d)
		}

		if onlineCounter >= d.MaxOnlineCount && onlineCounter <= d.MaxOnlineCount+2 {
			onlineCallback(onlineCounter, d)
		}

		select {
		case <-ctx.Done():
			return // gracefully shutdown
		case <-time.After(duration):
			continue
		}
	}
}
