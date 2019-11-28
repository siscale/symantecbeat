// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/marian-craciunescu/symantecbeat/ecs"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

const (
	loginURL       = "https://api.sep.securitycloud.symantec.com/v1/oauth2/tokens"
	eventExportURL = "https://usea1.r3.securitycloud.symantec.com/r3_epmp_i/sccs/v1/events/export"
	eventsearchURL = "https://api.sep.securitycloud.symantec.com/v1/event-search"
)

type SymantecClient struct {
	CustomerID   string
	DomainID     string
	ClientID     string
	ClientSecret string
	oauthToken   string
	mapper       *ecs.Mapper
	logger       *logp.Logger
}

func NewSymantecClient(customerID, domainID, clientID, clientSecret string, mapper *ecs.Mapper) SymantecClient {

	fmt.Printf("Using \ncustomerID=%s\ndomainID=%s\nclientID=%s\nclientSecret=%s\n", customerID, domainID, clientID, clientSecret)

	return SymantecClient{
		CustomerID:   customerID,
		DomainID:     domainID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		mapper:       mapper,
		logger:       logp.NewLogger("symantec_client"),
	}

}

func (s *SymantecClient) GetOauthToken() error {
	client := &http.Client{}

	b64Signature := s.encodeToBase64()

	fmt.Println(b64Signature)

	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("scope", "domain")

	req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		s.logger.Error(err)
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", b64Signature))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-epmp-domain-id", s.DomainID)
	req.Header.Add("x-epmp-customer-id", s.CustomerID)

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error(err)
		return err

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		s.logger.Error(err)
		return err
	}

	var oauthResponse oauthResponse
	err = json.Unmarshal(body, &oauthResponse)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	s.oauthToken = oauthResponse.Token
	s.logger.Infof("Acquired token with %s  valid for %d", s.oauthToken, oauthResponse.Expires)
	return nil
}

func (s *SymantecClient) encodeToBase64() string {
	authorizationRawValue := fmt.Sprintf("%s:%s", s.ClientID, s.ClientSecret)
	authValue := base64.StdEncoding.EncodeToString([]byte(authorizationRawValue))
	return authValue
}

type oauthResponse struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	TokenType string `json:"type"`
	Expires   int    `json:"expires_in"`
}

func (s *SymantecClient) getExportData(jsonValue []byte) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, eventExportURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		s.logger.Errorf("Error doing new request %v", err)
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.oauthToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-epmp-domain-id", s.DomainID)
	req.Header.Add("x-epmp-product", "SAEP")
	req.Header.Add("x-epmp-customer-id", s.CustomerID)

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error(err)
		return nil, err

	}

	s.logger.Infof("Server response=%i body=%s", resp.StatusCode, string(jsonValue))
	s.logger.Debugf("Server response=%i", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Error doing request to server err=%v", err)
		return nil, err
	}
	err = resp.Body.Close()

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Infof("HTTP non2xx.StatusCode=%d,status=%s body=%s",
			resp.StatusCode, resp.Status, string(body))
		return body, nil
	}

	return body, nil
}

func (s *SymantecClient) DoExportRequest(start, end time.Time, t EventType, size int) (mapStrArr []common.MapStr, err error) {

	s.logger.Infof("DoExportRequest for event=%s", t.String())

	requestBody, err := NewEventExportEncoded(start, end, size, t)
	if err != nil {
		s.logger.Errorf("error encoding export event as  json err=%s", err.Error())
		return nil, err
	}

	batches := 0
	noOfEvents := 0
	for {
		response, err := s.getExportData(requestBody)
		if err != nil {
			s.logger.Errorf("error doing  request %s", err.Error())
			return nil, err
		}

		if len(response) == 0 || string(response) == "[]" {
			s.logger.Infof("Finished request no_of_batches=%d ", batches)
			break
		} else {
			reader := bytes.NewReader(response)
			dec := json.NewDecoder(reader)

			var m []map[string]interface{}
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				s.logger.Errorf("error decoding json response err=%s", err.Error())
				return nil, err
			}

			for i := range m {
				newMap := m[i]
				mapStr, err := s.transformToMapStr(newMap)
				if err != nil {
					return nil, err
				} else {
					mapStrArr = append(mapStrArr, mapStr)
					noOfEvents++
				}
			}
		}

		batches++
	}
	s.logger.Infof("For type=%s  got %d in %d  batches", t.String(), noOfEvents, batches)
	return mapStrArr, nil
}

func (s *SymantecClient) getSearchData(jsonValue []byte) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, eventsearchURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		s.logger.Errorf("Error doing new request %v", err)
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.oauthToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-epmp-domain-id", s.DomainID)
	req.Header.Add("x-epmp-customer-id", s.CustomerID)

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Errorf("Error doing POST err=%v", err)
		return nil, err

	}

	s.logger.Infof("Server response=%i body=%s", resp.StatusCode, string(jsonValue))
	s.logger.Debugf("Server response=%i", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Error doing request to server err=%v", err)
		return nil, err
	}
	err = resp.Body.Close()

	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Infof("HTTP non2xx.StatusCode=%d,status=%s body=%s",
			resp.StatusCode, resp.Status, string(body))
		return body, nil
	}

	return body, nil
}

func (s *SymantecClient) DoRetrieveSearchEvents(start time.Time, end time.Time, size int) (mapStrArr []common.MapStr, err error) {
	s.logger.Infof("DoRetrieveSearchEvents for ALL events type")
	next := 0
	noOfEvents := 0
	batches := 0
	for {
		requestBody, err := NewEventSearchEncoded(start, end, size, next, ALL)
		if err != nil {
			s.logger.Errorf("error encoding search event as  json err=%s", err.Error())
			return nil, err
		}

		body, err := s.getSearchData(requestBody)
		if err != nil {
			logp.Err("error doing search request %s", err.Error())
			return nil, err
		}

		var event eventResponse
		err = json.Unmarshal(body, &event)
		if err != nil {
			s.logger.Errorf("error decoding search response json  err=%s", err.Error())
			return nil, err

		}
		noOfBatchEvent := 0
		for i := range event.Events {
			newMap := event.Events[i]
			mapStr, err := s.transformToMapStr(newMap)
			if err != nil {
				return nil, err
			} else {
				mapStrArr = append(mapStrArr, mapStr)
				noOfBatchEvent++
			}
		}
		noOfEvents += noOfBatchEvent
		s.logger.Infof("Total no_of_event=%d next=%d", noOfEvents, next)
		next += noOfBatchEvent
		batches++
		if event.Total <= next {
			break
		}

	}
	s.logger.Infof("Got no_of_event=%d in batches=%d", noOfEvents, batches)
	return mapStrArr, nil
}

func (s *SymantecClient) transformToMapStr(initialMap map[string]interface{}) (common.MapStr, error) {
	mapStr := common.MapStr{}
	s.recurseAndNormalizeMap("", mapStr, initialMap)
	return mapStr, nil
}

func (s *SymantecClient) recurseAndNormalizeMap(parentKey string, result common.MapStr, initialMap map[string]interface{}) {

	for k := range initialMap {
		switch innerType := initialMap[k].(type) {
		case map[string]interface{}:
			{
				parentKey := fmt.Sprintf("%s.", k)
				s.recurseAndNormalizeMap(parentKey, result, innerType)
			}
		case float32, float64, int, int8, int16, int32, int64, string, bool:
			actualKey := parentKey + k
			ecsField := s.mapper.EcsField(actualKey)

			_, err := result.Put(ecsField, initialMap[k])
			if err != nil {
				logp.Err("error puting field in map err=%s", err.Error())
			}
		default:
			fmt.Println(innerType)
		}
	}

}

type eventResponse struct {
	Total  int                      `json:"total"`
	Next   int                      `json:"next"`
	Events []map[string]interface{} `json:"events"`
}
