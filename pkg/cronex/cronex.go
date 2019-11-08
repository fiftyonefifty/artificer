package cronex

import (
	"github.com/robfig/cron/v3"
)

func AddFunc(c *cron.Cron, immediate bool, immediateDoneChannel chan bool, spec string, cmd func()) (cron.EntryID, error) {
	if immediate {
		var entryId *cron.EntryID
		eid, _ := c.AddFunc("@every 1s", func() {
			c.Remove(*entryId)
			if immediateDoneChannel != nil {
				defer func() {
					immediateDoneChannel <- true
				}()
			}
			cmd()
		})
		entryId = &eid
	}
	return c.AddFunc(spec, cmd)
}
