package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/thesky9531/lareina/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	client *http.Client
}

func New() (client *Client) {
	client = &Client{
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig: &tls.Config{
					// TODO:测试时忽略证书校验
					InsecureSkipVerify: true,
				},
			},
		},
	}
	return
}

func NewHttps(caCertPath, certFile, keyFile string) (client *Client) {
	// 创建证书池及各类对象
	var pool *x509.CertPool // 我们要把一部分证书存到这个池中
	var caCrt []byte        // 根证书
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.ErrLog("", err)
		panic(err)
	}
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)

	//var cliCrt tls.Certificate  // 具体的证书加载对象
	//cliCrt, err = tls.LoadX509KeyPair(certFile, keyFile)
	//if err != nil {
	//	return nil, err
	//}

	client = &Client{
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					//RootCAs:      pool,
					//Certificates: []tls.Certificate{cliCrt},
					InsecureSkipVerify: true, //忽略验证证书
				},
			},
		},
	}
	return
}

func (h *Client) Close() {
	h.client.CloseIdleConnections()
	return
}

func (h *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	req = req.WithContext(ctx)
	resp, err := h.client.Do(req)
	if err != nil {
		log.ErrLog("", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	return UnMarshalResponse(resp.Body)
}

func (h *Client) Download(ctx context.Context, req *http.Request) (*File, error) {
	req = req.WithContext(ctx)
	resp, err := h.client.Do(req)
	if err != nil {
		log.ErrLog("", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.New(resp.Status)
	}
	name := resp.Header.Get("Name")
	if name == "" {
		resp.Body.Close()
		return nil, errors.New("name is empty")
	}
	name, err = url.QueryUnescape(name)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	sizeStr := resp.Header.Get("Size")
	if sizeStr == "" {
		resp.Body.Close()
		return nil, errors.New("size is empty")
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	typeStr := resp.Header.Get("Type")
	if typeStr == "" {
		resp.Body.Close()
		return nil, errors.New("type is empty")
	}
	typeStr, err = url.QueryUnescape(typeStr)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	file := &File{
		body: resp.Body,
		name: name,
		size: size,
		tp:   typeStr,
	}
	return file, nil
}
