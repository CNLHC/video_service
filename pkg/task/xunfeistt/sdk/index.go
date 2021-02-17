package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

func tomap(obj interface{}) (newMap map[string]string, err error) {
	data, err := json.Marshal(obj) // Convert to a json string
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

func (c *XunfeiSDK) Prepare(req PrepareReq) (resp BaseResp, err error) {
	var fullreq PrepareFullReq
	fullreq, err = c.GetPrepareReq()
	if err != nil {
		return
	}
	fullreq.Language = req.Language
	resp, err = c.request(c.BaseUrl+"/prepare", fullreq)
	err2 := c.checkBaseResp(resp)
	if err == nil && err2 == nil {
		c.taskid = resp.Data
	}
	if err2 != nil {
		err = err2
	}
	return
}

func (c *XunfeiSDK) Upload() (resp BaseResp, err error) {
	file, err := os.Open(c.file_path)
	if err != nil {
		err = errors.Wrap(err, c.file_path)
		return
	}
	var (
		buf     = make([]byte, SliceSize)
		rawreq  *http.Request
		rawresp *http.Response
		body    []byte
	)
	c.ResetSliceId()
	endpoint := c.BaseUrl + "/upload"
	for {
		var (
			nbytes = 0
			b      bytes.Buffer
			fw     io.Writer
			reqmap map[string]string
			w      *multipart.Writer
		)
		if nbytes, err = file.Read(buf); err != nil {
			if errors.Is(err, io.EOF) {
				log.Info().Msgf("task %s upload over")
				err = nil
				return
			}
			return
		}
		w = multipart.NewWriter(&b)
		req := c.GetReq()
		if reqmap, err = tomap(req); err != nil {
			return
		}
		for k, v := range reqmap {
			if fw, err = w.CreateFormField(k); err != nil {
				return
			}
			fw.Write([]byte(v))
		}
		if fw, err = w.CreateFormField("slice_id"); err != nil {
			return
		}
		fw.Write([]byte(c.cur_sliceid))

		header := make(textproto.MIMEHeader)

		log.Info().Msgf("task %s upload : %+v", string(b.Bytes()))
		fn := filepath.Base(c.file_path)
		header.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
				"content", fn))

		header.Set("Content-Type", "application/octet-stream")
		_ = header

		if fw, err = w.CreateFormFile("content", fn); err != nil {
			return
		}

		_ = nbytes
		fw.Write(buf[:nbytes])
		w.Close()

		cli := http.Client{}

		if rawreq, err = http.NewRequest(http.MethodPost, endpoint, &b); err != nil {
			return
		}

		rawreq.Header.Set("Content-Type", w.FormDataContentType())
		rawresp, err = cli.Do(rawreq)
		body, err = ioutil.ReadAll(rawresp.Body)
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return
		}
		if err = c.checkBaseResp(resp); err != nil {
			return
		}
		c.GetNextSliceId()
	}
}

func (c *XunfeiSDK) request(urlstr string, req interface{}) (resp BaseResp, err error) {
	client := &http.Client{}
	data := url.Values{}
	var (
		mapreq  map[string]string
		rawreq  *http.Request
		rawresp *http.Response
	)
	if mapreq, err = tomap(req); err != nil {
		return
	}
	log.Info().Msgf("xunfei raw: %+v \n req: %+v", req, mapreq)
	for k, v := range mapreq {
		data.Set(k, v)
	}

	if rawreq, err = http.NewRequest(
		http.MethodPost,
		urlstr,
		strings.NewReader(data.Encode())); err != nil {
		return
	}
	rawreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rawreq.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	if rawresp, err = client.Do(rawreq); err != nil {
		return
	}

	body, err := ioutil.ReadAll(rawresp.Body)

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}
	err = c.checkBaseResp(resp)
	return
}
func (c *XunfeiSDK) Merge() (resp BaseResp, err error) {
	req := c.GetReq()
	resp, err = c.request(c.BaseUrl+"/merge", req)
	return

}

func (c *XunfeiSDK) GetProgress() (sresp Status, err error) {
	req := c.GetReq()
	resp, err := c.request(c.BaseUrl+"/getProgress", req)
	if err == nil {
		err = json.Unmarshal([]byte(resp.Data), &sresp)
	}
	return
}

func (c *XunfeiSDK) GetResult() (resp BaseResp, err error) {
	req := c.GetReq()
	resp, err = c.request(c.BaseUrl+"/getResult", req)
	return
}
