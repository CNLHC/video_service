package video

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"os/exec"
)

type ProberResp struct {
	Format ProbeFormat `json:"format"`
}

type ProbeFormat struct {
	Filename       string    `json:"filename"`
	NbStreams      int64     `json:"nb_streams"`
	NbPrograms     int64     `json:"nb_programs"`
	FormatName     string    `json:"format_name"`
	FormatLongName string    `json:"format_long_name"`
	StartTime      string    `json:"start_time"`
	Duration       string    `json:"duration"`
	Size           string    `json:"size"`
	BitRate        string    `json:"bit_rate"`
	ProbeScore     int64     `json:"probe_score"`
	Tags           ProbeTags `json:"tags"`
}

type ProbeTags struct {
	MajorBrand       string `json:"major_brand"`
	MinorVersion     string `json:"minor_version"`
	CompatibleBrands string `json:"compatible_brands"`
	Encoder          string `json:"encoder"`
}

//ffprobe https://publicstatic.cnworkshop.xyz/index.mp4 -of json -v quiet  -show_format
type Prober struct {
	cmd exec.Cmd
}

func (Prober) Probe(fp string) (resp ProberResp, err error) {
	var (
		stderr io.ReadCloser
		stdout io.ReadCloser
		buf    []byte
		errbuf []byte
	)
	cmd := exec.Command("ffprobe",
		"-of", "json",
		"-v", "quiet",
		"-show_format",
		fp,
	)

	stdout, err = cmd.StdoutPipe()
	stderr, err = cmd.StderrPipe()
	log.Info().Msgf("probe of %s. cmd is: %s", fp, cmd.String())
	err = cmd.Start()
	buf, err = ioutil.ReadAll(stdout)
	errbuf, err = ioutil.ReadAll(stderr)
	err = cmd.Wait()
	if err != nil {
		return resp, err
	}
	if err != nil {
		return resp, err
	}
	log.Info().Msgf("probe of %s is %s", fp, string(buf))
	log.Info().Msgf("stderr of %s is  %s", fp, string(errbuf))
	err = json.Unmarshal(buf, &resp)
	return resp, err
}
