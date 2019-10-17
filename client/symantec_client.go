package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	loginURL = "/oauth2/tokens"
	eventURL = "/sccs/v1/events/export"
)

type SymantecClient struct {
	ApiURL       string
	CustomerID   string
	DomainID     string
	ClientID     string
	ClientSecret string
	oauthToken   string
	logger       *logp.Logger
}

func NewSymantecClient(apiURL, customerID, domainID, clientID, clientSecret string) SymantecClient {

	fmt.Printf("Using \ncustomerID=%s\ndomainID=%s\nclientID=%s\nclientSecret=%s\n", customerID, domainID, clientID, clientSecret)

	return SymantecClient{
		ApiURL:       apiURL,
		CustomerID:   customerID,
		DomainID:     domainID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		logger:       logp.NewLogger("symantec_client"),
	}

}

func (s *SymantecClient) GetOauthToken() error {
	client := &http.Client{}

	b64Signature := s.encodeToBase64()

	fmt.Println(b64Signature)

	uri := s.ApiURL + loginURL

	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("scope", "domain")

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBufferString(data.Encode()))
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

func (s *SymantecClient) getData(jsonValue []byte) ([]byte, error) {

	client := &http.Client{}

	uri := s.ApiURL + eventURL

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonValue))
	if err != nil {
		s.logger.Error(err)
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

	s.logger.Debugf("Server response=%i", resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		s.logger.Error(err)
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return body, nil
}

func (s *SymantecClient) DoRequest(start, end time.Time, t EventType, size int) (mapStrArr []common.MapStr, err error) {

	logp.Info("DoRequest for event=%s", t.String())

	requestBody, err := NewEventEncoded(start, end, size, t)
	if err != nil {
		fmt.Println(err)
	}

	batches := 0
	noOfEvents := 0
	for {
		response, err := s.getData(requestBody)
		if err != nil {
			logp.Err("error doing  request %s", err.Error())
			return nil, err
		}

		if len(response) == 0 || string(response) == "[]" {
			fmt.Println("Finished request no_of_batches", batches)
			break
		} else {
			reader := bytes.NewReader(response)
			dec := json.NewDecoder(reader)

			var m []map[string]interface{}
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				logp.Err("error decoding json response err=%s", err.Error())
				return nil, err
			}

			for i := range m {
				newMap := m[i]
				mapStr, err := transformToMapStr(newMap)
				if err != nil {
					return nil, err
				} else {
					mapStr.Put("event_type", t.String())
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

func transformToMapStr(intialMap map[string]interface{}) (common.MapStr, error) {
	mapStr := common.MapStr{}
	for k := range intialMap {
		_, err := mapStr.Put(k, intialMap[k])
		if err != nil {
			logp.Err("error puting field in map err=%s", err.Error())
			return nil, err
		}
	}
	return mapStr, nil
}
