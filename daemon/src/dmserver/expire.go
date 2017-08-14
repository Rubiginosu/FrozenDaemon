package dmserver

import (
	"colorlog"
	"strconv"
	"time"
)

func serverOutDateClearer() {
	for {
		for k, v := range serverSaved {
			if v.Expire <= time.Now().Unix() {
				colorlog.LogPrint("Server" + strconv.Itoa(k) + " was out of date.Delete it now.")
				v.Delete()
				delete(serverSaved, k)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
