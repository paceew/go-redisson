package util

import (
	"bytes"
	"errors"
	"github.com/paceew/go-redisson/pkg/global"
	clog "github.com/paceew/go-redisson/pkg/log"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func PostJson(url string, body map[string]interface{}, timeout int) ([]byte, error) {
	bytesData, err := global.Json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Duration(timeout/2)*time.Second) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)) //设置发送接受数据超时
				return conn, nil
			},
		},
	}
	begin := time.Now()
	resp, err := client.Do(request)
	elapsed := time.Since(begin)
	if err != nil {
		clog.GetDefaultLogSingletons().E("PostJson external costtime:%f,url:%s,body:%+v,%s", elapsed.Seconds(), url, body, err.Error())
		return nil, err
	}
	if resp == nil {
		err := errors.New("resp is nil")
		clog.GetDefaultLogSingletons().E("PostJson external costtime:%f resp is nil, url:%s,body:%+v", elapsed.Seconds(), url, body)
		return nil, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clog.GetDefaultLogSingletons().E("PostJson external costtime:%f url:%s,body:%+v Read err:%s", elapsed.Seconds(), url, body, err.Error())
		return nil, err
	}
	response := respBytes[:]
	if resp.StatusCode != 200 {
		err := errors.New(resp.Status)
		clog.GetDefaultLogSingletons().E("PostJson external costtime:%f StatusCode:%d,url:%s,body:%+v,resp:%s", elapsed.Seconds(), resp.StatusCode, url, body, response)
		return nil, err
	}
	clog.GetDefaultLogSingletons().I("PostJson external costtime:%f %s %s", elapsed.Seconds(), url, response)
	return response, nil
}

func HttpGet(url string, timeout int) (response []byte, err error) {
	begin := time.Now()
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Duration(timeout/2)*time.Second)
				if err != nil {
					return conn, err
				}
				conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
				return conn, nil
			},
		},
	}

	resp, err := client.Get(url)
	elapsed := time.Since(begin)

	if err != nil {
		clog.GetDefaultLogSingletons().E("HttpGet external costtime:%f %s %s %v", elapsed.Seconds(), url, err.Error(), resp)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clog.GetDefaultLogSingletons().E("HttpGet external costtime:%f %s %s", elapsed.Seconds(), url, err.Error())
		return nil, err
	}

	response = body[:]
	if resp.StatusCode != 200 {
		err := errors.New(resp.Status)
		clog.GetDefaultLogSingletons().E("HttpGet external costtime:%f StatusCode:%d,url:%s,resp:%s", elapsed.Seconds(), resp.StatusCode, url, response)
		return nil, err
	}
	clog.GetDefaultLogSingletons().I("HttpGet external costtime:%f,url:%s,resp:%s", elapsed.Seconds(), url, response)
	return response, nil

}
