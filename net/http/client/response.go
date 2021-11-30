package client

import (
	"encoding/json"
	"errors"
	"github.com/thesky9531/lareina/log"
	"io"
	"io/ioutil"
)

var unMarshalError = errors.New("解析body失败,未识别的返回值格式")

type Response struct {
	Data   json.RawMessage
	Status int
	Error  string
}

func UnMarshalResponse(body io.Reader) (*Response, error) {
	resp := new(Response)
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.ErrLog("", err)
		return resp, err
	}
	set := make(map[string]json.RawMessage)
	err = json.Unmarshal(bodyBytes, &set)
	if err != nil {
		log.ErrLog(string(bodyBytes), err)
		return resp, err
	}
	stu, ok := set["status"]
	if !ok {
		err = unMarshalError
		log.ErrLog(string(bodyBytes), err)
		return resp, err
	}
	err = json.Unmarshal(stu, &resp.Status)
	if err != nil {
		log.ErrLog("", err)
		return resp, unMarshalError
	}
	resp.Error = string(set["error"])
	resp.Data = set["data"]
	return resp, nil
}

type File struct {
	body io.ReadCloser
	name string
	size int64
	tp   string
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.body.Read(p)
}

func (f *File) Close() error {
	return f.body.Close()
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int64 {
	return f.size
}

func (f *File) Type() string {
	return f.tp
}
