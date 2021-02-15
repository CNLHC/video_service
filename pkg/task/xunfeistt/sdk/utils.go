package sdk

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	_ "fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const SliceSize = 10485760

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
	c.cur_sliceid = "aaaaaaaaa"
}

func (*XunfeiSDK) checkBaseResp(r BaseResp) error {
	if r.Ok == 0 {
		return nil
	} else {
		return errors.New("error")
	}
}
func (c *XunfeiSDK) getSliceNum(size int64) int64 {
	return int64(math.Ceil(float64(size) / float64(size)))
}

func (c *XunfeiSDK) GetPrepareReq() (resp PrepareFullReq, err error) {
	fp, err := os.Open(c.file_path)
	if err != nil {
		return
	}
	defer fp.Close()
	info, err := fp.Stat()
	if info != nil {
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
	resp.BaseReq = c.base
	resp.TaskId = c.taskid
	return
}

func GetXunfeiSDK(APPId string, sk string, file_path string) (resp *XunfeiSDK, err error) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	_, err = os.Lstat(file_path)
	if err != nil {
		return
	}
	resp = &XunfeiSDK{
		BaseUrl: "https://raasr.xfyun.cn/api",
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
