package dmserver

import "time"

func serverOutDateClearer() {
	for {
		for k, v := range serverSaved {
			if v.Expire <= time.Now().Unix() {
				v.Delete()
				delete(serverSaved, k)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
