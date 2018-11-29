// Copyright (c) 2018 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	aliyunECSRamURL      = "http://100.100.100.200/latest/meta-data/ram/security-credentials/"
	expirationTimeFormat = "2006-01-02T15:04:05Z"
)

type SecurityTokenResult struct {
	AccessKeyId     string
	AccessKeySecret string
	Expiration      string
	SecurityToken   string
	Code            string
	LastUpdated     string
}

func getToken() (result []byte, err error) {
	client := http.Client{
		Timeout: time.Second * 3,
	}
	var respList *http.Response
	respList, err = client.Get(aliyunECSRamURL)
	if err != nil {
		fmt.Println("UPDATE_STS_ALARM", " get role list error ", err)
		return nil, err
	}
	defer respList.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(respList.Body)
	if err != nil {
		fmt.Println("UPDATE_STS_ALARM", " parse role list error ", err)
		return nil, err
	}

	bodyStr := string(body)
	bodyStr = strings.TrimSpace(bodyStr)
	roles := strings.Split(bodyStr, "\n")
	role := roles[0]

	var respGet *http.Response
	respGet, err = client.Get(aliyunECSRamURL + role)
	if err != nil {
		fmt.Println("UPDATE_STS_ALARM", " get token error ", err, " role ", role)
		return nil, err
	}
	defer respGet.Body.Close()
	body, err = ioutil.ReadAll(respGet.Body)
	if err != nil {
		fmt.Println("UPDATE_STS_ALARM", " parse token error ", err, " role ", role)
		return nil, err
	}
	return body, nil
}

func UpdateTokenFunction() (accessKeyID, accessKeySecret, securityToken string, expireTime time.Time, err error) {
	var tokenResultBuffer []byte
	for tryTime := 0; tryTime < 3; tryTime++ {
		tokenResultBuffer, err = getToken()
		if err != nil {
			continue
		}
		var tokenResult SecurityTokenResult
		err = json.Unmarshal(tokenResultBuffer, &tokenResult)
		if err != nil {
			fmt.Println("UPDATE_STS_ALARM", " unmarshal token error ", err, " token ", string(tokenResultBuffer))
			continue
		}
		if strings.ToLower(tokenResult.Code) != "success" {
			tokenResult.AccessKeySecret = "xxxxx"
			tokenResult.SecurityToken = "xxxxx"
			fmt.Println("UPDATE_STS_ALARM", " token code not success ", err, " result ", tokenResult)
			continue
		}
		expireTime, err := time.Parse(expirationTimeFormat, tokenResult.Expiration)
		if err != nil {
			tokenResult.AccessKeySecret = "xxxxx"
			tokenResult.SecurityToken = "xxxxx"
			fmt.Println("UPDATE_STS_ALARM ", " parse time error ", err, " result ", tokenResult)
			continue
		}
		fmt.Println("get security token success, id ", tokenResult.AccessKeyId, " expire ", tokenResult.Expiration, " last update", tokenResult.LastUpdated)
		return tokenResult.AccessKeyId, tokenResult.AccessKeySecret, tokenResult.SecurityToken, expireTime, nil
	}
	return accessKeyID, accessKeySecret, securityToken, expireTime, err
}
