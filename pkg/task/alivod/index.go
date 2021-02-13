package alivod

import (
	"argus/video/pkg/task"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vod"
)

type AliVodTask struct {
	VideoId  string
	cli      *vod.Client
	VodJobId string
}
type AliVodTaskCfg struct {
	VideoId string
}

func (c *AliVodTask) GetVodCli() *vod.Client {
	return c.cli
}

func (c *AliVodTask) Init(cfg interface{}) (err error) {
	switch cfg.(type) {
	case AliVodTaskCfg:
		t := cfg.(AliVodTaskCfg)
		client, err := vod.NewClientWithAccessKey("cn-qingdao", "<accessKeyId>", "<accessSecret>")
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
