package sdk

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	_ "fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const SliceSize = 10485760

func (c *XunfeiSDK) SetTaskID(id string) {
	c.taskid = id
}
func (c *XunfeiSDK) GetNextSliceId() string {
	t := []byte(c.cur_sliceid)
	j := len(t) - 1
	for j >= 0 {
		if t[j] != 'z' {
			t[j] = t[j] + 1
			break
		} else {
			t[j] = 'a'
			j = j - 1
		}
	}
	c.cur_sliceid = string(t)
	return c.cur_sliceid
}

func (c *XunfeiSDK) ResetSliceId() {
	c.cur_sliceid = "aaaaaaaaaa"
}

func (*XunfeiSDK) checkBaseResp(r BaseResp) error {
	if r.Ok == 0 {
		return nil
	} else {
		msg := fmt.Sprintf("error %+v", r)
		return errors.New(msg)
	}
}
func (c *XunfeiSDK) getSliceNum(size int64) int64 {
	return int64(math.Ceil(float64(size) / float64(SliceSize)))
}

func (c *XunfeiSDK) GetPrepareReq() (resp PrepareFullReq, err error) {
	fp, err := os.Open(c.file_path)
	if err != nil {
		err = errors.Wrap(err, c.file_path)
		return
	}
	defer fp.Close()
	info, err := fp.Stat()
	if err != nil {
		return
	}
	size := info.Size()

	resp.FileLen = strconv.FormatInt(size, 10)
	resp.FileName = info.Name()
	resp.SliceNum = strconv.FormatInt(c.getSliceNum(size), 10)
	resp.BaseReq = c.base
	return
}

func (c *XunfeiSDK) GetReq() (resp TaskIdReq) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	signa := signature(c.base.AppID, c.sk, ts)
	c.base.Ts = ts
	c.base.Signa = signa
	resp.BaseReq = c.base
	resp.TaskId = c.taskid
	return
}

func GetXunfeiSDK(APPId string, sk string, file_path string) (resp *XunfeiSDK, err error) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	_, err = os.Lstat(file_path)
	if err != nil {
		err = errors.Wrap(err, file_path)
		return
	}
	resp = &XunfeiSDK{
		BaseUrl:   "http://raasr.xfyun.cn/api",
		file_path: file_path,
		aid:       APPId,
		sk:        sk,
		base: BaseReq{
			AppID: APPId,
			Ts:    ts,
			Signa: signature(APPId, sk, ts),
		},
	}
	return
}

func signature(id string, sk string, ts string) (signa string) {
	base := id + ts
	t := md5.Sum([]byte(base))
	hashbase := hex.EncodeToString(t[:])
	h := hmac.New(sha1.New, []byte(sk))
	h.Write([]byte(hashbase))
	signa = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
