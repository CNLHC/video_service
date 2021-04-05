package alivod

import (
	"argus/video/pkg/config"
	"argus/video/pkg/task"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/mts"
)

type AliVodTask struct {
	VideoId  string
	cli      *mts.Client
	VodJobId string
}
type AliVodTaskCfg struct {
	VideoId string
}

func (c *AliVodTask) GetVodCli() *mts.Client {
	return c.cli
}

func (c *AliVodTask) Init(cfg interface{}) (err error) {
	switch cfg.(type) {
	case AliVodTaskCfg:
		t := cfg.(AliVodTaskCfg)
		client, err := mts.NewClientWithAccessKey("cn-beijing",
			config.Get("Ali_AID"),
			config.Get(("Ali_AKEY")))
		if err != nil {
			return err
		}
		c.cli = client
		c.VideoId = t.VideoId
	default:
		return task.ErrWrongCfg
	}
	return err
}
