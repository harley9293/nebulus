package service

import (
	"context"
	"sync"
	"time"
)

func monitorLoop(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	delServices := make([]string, 0)
	for {
		select {
		case <-time.After(16 * time.Millisecond):
			m.rwLock.RLock()
			for _, c := range m.serviceByName {
				if isServiceExist(c) {
					delServices = append(delServices, c.name)
				}
			}
			m.rwLock.RUnlock()
		case <-ctx.Done():
			return
		}

		if len(delServices) > 0 {
			m.rwLock.Lock()
			for _, name := range delServices {
				delete(m.serviceByName, name)
			}
			m.rwLock.Unlock()
		}
		delServices = delServices[:0]
	}
}

func isServiceExist(c *service) bool {
	select {
	case <-c.exit:
		return true
	default:
		return false
	}
}
