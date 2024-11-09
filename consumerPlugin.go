////////////////////////////////////////////////////////////////////////////////
//	consumerPlugin.go  -  Sep-25-2024  -  aldebap
//
//	Kong consumer Plugins
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// kong consumer ID payload
type KongConsumerID struct {
	Id string `json:"id"`
}

// kong Basic Auth config attributes
type KongBasicAuthConfig struct {
	userName string
	password string
}

// create a new kong Basic Auth config
func NewKongBasicAuthConfig(userName string, password string) *KongBasicAuthConfig {

	return &KongBasicAuthConfig{
		userName: userName,
		password: password,
	}
}

// kong basic auth request payload
type KongBasicAuthRequest struct {
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// kong response payload
type KongBasicAuthResponse struct {
	Id string `json:"id"`
}

// add a new consumer to Kong
func (ks *KongServerDomain) AddConsumerBasicAuth(id string, newKongBasicAuthConfig *KongBasicAuthConfig, options Options) error {

	var consumerPluginURL string = fmt.Sprintf("%s/%s/%s/%s", ks.ServerURL(), consumersResource, id, basicAuthPlugins)

	payload, err := json.Marshal(KongBasicAuthRequest{
		UserName: newKongBasicAuthConfig.userName,
		Password: newKongBasicAuthConfig.password,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", consumerPluginURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("consumer not found")
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("fail sending add consumer basic auth command to Kong: " + resp.Status)
	}

	//	parse response payload
	var respPayload []byte

	respPayload, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var consumerBasicAuthResp KongBasicAuthResponse

	err = json.Unmarshal(respPayload, &consumerBasicAuthResp)
	if err != nil {
		return err
	}

	if options.jsonOutput {
		fmt.Printf("%s\n%s\n", resp.Status, string(respPayload))
	} else {
		if options.verbose {
			fmt.Printf("http response status code: %s\nnew plugin ID: %s\n", resp.Status, consumerBasicAuthResp.Id)
		} else {
			fmt.Printf("%s\n", consumerBasicAuthResp.Id)
		}
	}

	return nil
}

// kong KeyAuth config attributes
type KongKeyAuthConfig struct {
	key string
	ttl int64
}

// create a new kong KeyAuth config
func NewKongKeyAuthConfig(key string, ttl int64) *KongKeyAuthConfig {

	return &KongKeyAuthConfig{
		key: key,
		ttl: ttl,
	}
}

// kong consumer KeyAuth request payload
type KongKeyAuthRequest struct {
	Key string `json:"key,omitempty"`
	Ttl int64  `json:"ttl,omitempty"`
}

// kong consumer KeyAuth response payload
type KongKeyAuthResponse struct {
	Id string `json:"id"`
}

func (ks *KongServerDomain) AddConsumerKeyAuth(id string, newKongKeyAuthConfig *KongKeyAuthConfig, options Options) error {

	var consumerPluginURL string = fmt.Sprintf("%s/%s/%s/%s", ks.ServerURL(), consumersResource, id, keyAuthPlugins)

	payload, err := json.Marshal(KongKeyAuthRequest{
		Key: newKongKeyAuthConfig.key,
		Ttl: newKongKeyAuthConfig.ttl,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", consumerPluginURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("consumer not found")
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("fail sending add consumer keyAuth command to Kong: " + resp.Status)
	}

	//	parse response payload
	var respPayload []byte

	respPayload, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var consumerKeyAuthResp KongKeyAuthResponse

	err = json.Unmarshal(respPayload, &consumerKeyAuthResp)
	if err != nil {
		return err
	}

	if options.jsonOutput {
		fmt.Printf("%s\n%s\n", resp.Status, string(respPayload))
	} else {
		if options.verbose {
			fmt.Printf("http response status code: %s\nnew plugin ID: %s\n", resp.Status, consumerKeyAuthResp.Id)
		} else {
			fmt.Printf("%s\n", consumerKeyAuthResp.Id)
		}
	}

	return nil
}

// kong JWT config attributes
type KongJWTConfig struct {
	algorithm string
	key       string
	secret    string
}

// create a new kong JWT config
func NewKongJWTConfig(algorithm string, key string, secret string) *KongJWTConfig {

	return &KongJWTConfig{
		algorithm: algorithm,
		key:       key,
		secret:    secret,
	}
}

// kong consumer JWT request payload
type KongJWTRequest struct {
	Algorithm string `json:"algorithm,omitempty"`
	Key       string `json:"key,omitempty"`
	Secret    string `json:"secret,omitempty"`
}

// kong consumer JWT response payload
type KongJWTResponse struct {
	Id        string          `json:"id"`
	Consumer  *KongConsumerID `json:"service,omitempty"`
	Algorithm string          `json:"algorithm,omitempty"`
	Key       string          `json:"key,omitempty"`
	Secret    string          `json:"secret,omitempty"`
	Tags      []string        `json:"tags"`
}

func (ks *KongServerDomain) AddConsumerJWT(id string, newKongJWTConfig *KongJWTConfig, options Options) error {

	var consumerPluginURL string = fmt.Sprintf("%s/%s/%s/%s", ks.ServerURL(), consumersResource, id, jwtPlugins)

	payload, err := json.Marshal(KongJWTRequest{
		Algorithm: newKongJWTConfig.algorithm,
		Key:       newKongJWTConfig.key,
		Secret:    newKongJWTConfig.secret,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", consumerPluginURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("consumer not found")
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("fail sending add consumer JWT command to Kong: " + resp.Status)
	}

	//	parse response payload
	var respPayload []byte

	respPayload, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var consumerJWTResp KongJWTResponse

	err = json.Unmarshal(respPayload, &consumerJWTResp)
	if err != nil {
		return err
	}

	if options.jsonOutput {
		fmt.Printf("%s\n%s\n", resp.Status, string(respPayload))
	} else {
		if options.verbose {
			fmt.Printf("http response status code: %s\nnew plugin ID: %s\n", resp.Status, consumerJWTResp.Id)
		} else {
			fmt.Printf("%s\n", consumerJWTResp.Id)
		}
	}

	return nil
}

// kong IP Restriction config attributes
type KongIPRestrictionConfig struct {
	allow []string
	deny  []string
}

// kong IP Restriction plugin attributes
type KongIPRestrictionPlugin struct {
	name   string
	config *KongIPRestrictionConfig
}

// create a new kong IP Restriction plugin
func NewKongIPRestrictionPlugin(name string, allow []string, deny []string) *KongIPRestrictionPlugin {

	return &KongIPRestrictionPlugin{
		name: name,
		config: &KongIPRestrictionConfig{
			allow: allow,
			deny:  deny,
		},
	}
}

// kong plugin request config
type KongIPRestrictionRequestConfig struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}

// kong consumer JWT request payload
type KongIPRestrictionRequest struct {
	Name         string                          `json:"name,omitempty"`
	InstanceName string                          `json:"instance_name,omitempty"`
	Config       *KongIPRestrictionRequestConfig `json:"config,omitempty"`
}

// kong consumer JWT response payload
type KongIPRestrictionResponse struct {
	Id           string                         `json:"id"`
	Name         string                         `json:"name,omitempty"`
	InstanceName string                         `json:"instance_name,omitempty"`
	Config       KongIPRestrictionRequestConfig `json:"config,omitempty"`
}

func (ks *KongServerDomain) AddConsumerIPRestriction(id string, newKongIPRestrictionConfig *KongIPRestrictionPlugin, options Options) error {

	var consumerPluginURL string = fmt.Sprintf("%s/%s/%s/%s", ks.ServerURL(), consumersResource, id, pluginsResource)

	payload, err := json.Marshal(KongIPRestrictionRequest{
		Name:         IPRestrictionPlugins,
		InstanceName: newKongIPRestrictionConfig.name,
		Config: &KongIPRestrictionRequestConfig{
			Allow: newKongIPRestrictionConfig.config.allow,
			Deny:  newKongIPRestrictionConfig.config.deny,
		},
	})
	if err != nil {
		return err
	}

	//log.Printf("[debug] URL: %s", consumerPluginURL)
	//log.Printf("[debug] post payload: %s", payload)

	req, err := http.NewRequest("POST", consumerPluginURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("consumer not found")
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("fail sending add consumer IP Restriction command to Kong: " + resp.Status)
	}

	//	parse response payload
	var respPayload []byte

	respPayload, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var consumerIPRestrictionResp KongIPRestrictionResponse

	err = json.Unmarshal(respPayload, &consumerIPRestrictionResp)
	if err != nil {
		return err
	}

	if options.jsonOutput {
		fmt.Printf("%s\n%s\n", resp.Status, string(respPayload))
	} else {
		if options.verbose {
			fmt.Printf("http response status code: %s\nnew plugin ID: %s\n", resp.Status, consumerIPRestrictionResp.Id)
		} else {
			fmt.Printf("%s\n", consumerIPRestrictionResp.Id)
		}
	}

	return nil
}

// kong Rate Limiting config attributes
type KongRateLimitingConfig struct {
	second       int32
	minute       int32
	hour         int32
	errorCode    int32
	errorMessage string
}

// kong Rate Limiting plugin attributes
type KongRateLimitingPlugin struct {
	name   string
	config *KongRateLimitingConfig
}

// create a new kong Rate Limiting plugin
func NewKongRateLimitingPlugin(name string, second int32, minute int32, hour int32, errorCode int32, errorMessage string) *KongRateLimitingPlugin {

	return &KongRateLimitingPlugin{
		name: name,
		config: &KongRateLimitingConfig{
			second:       second,
			minute:       minute,
			hour:         hour,
			errorCode:    errorCode,
			errorMessage: errorMessage,
		},
	}
}

// kong plugin request config
type KongRateLimitingRequestConfig struct {
	Second       int32  `json:"second,omitempty"`
	Minute       int32  `json:"minute,omitempty"`
	Hour         int32  `json:"hour,omitempty"`
	ErrorCode    int32  `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// kong consumer JWT request payload
type KongRateLimitingRequest struct {
	Name         string                         `json:"name,omitempty"`
	InstanceName string                         `json:"instance_name,omitempty"`
	Config       *KongRateLimitingRequestConfig `json:"config,omitempty"`
}

// kong consumer JWT response payload
type KongRateLimitingResponse struct {
	Id           string                        `json:"id"`
	Name         string                        `json:"name,omitempty"`
	InstanceName string                        `json:"instance_name,omitempty"`
	Config       KongRateLimitingRequestConfig `json:"config,omitempty"`
}

func (ks *KongServerDomain) AddConsumerRateLimiting(id string, newKongRateLimitingPlugin *KongRateLimitingPlugin, options Options) error {

	var consumerPluginURL string = fmt.Sprintf("%s/%s/%s/%s", ks.ServerURL(), consumersResource, id, pluginsResource)

	payload, err := json.Marshal(KongRateLimitingRequest{
		Name:         RateLimitingPlugins,
		InstanceName: newKongRateLimitingPlugin.name,
		Config: &KongRateLimitingRequestConfig{
			Second:       newKongRateLimitingPlugin.config.second,
			Minute:       newKongRateLimitingPlugin.config.minute,
			Hour:         newKongRateLimitingPlugin.config.hour,
			ErrorCode:    newKongRateLimitingPlugin.config.errorCode,
			ErrorMessage: newKongRateLimitingPlugin.config.errorMessage,
		},
	})
	if err != nil {
		return err
	}

	//log.Printf("[debug] URL: %s", consumerPluginURL)
	//log.Printf("[debug] post payload: %s", payload)

	req, err := http.NewRequest("POST", consumerPluginURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("consumer not found")
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("fail sending add consumer IP Restriction command to Kong: " + resp.Status)
	}

	//	parse response payload
	var respPayload []byte

	respPayload, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var consumerRateLimitingResponse KongRateLimitingResponse

	err = json.Unmarshal(respPayload, &consumerRateLimitingResponse)
	if err != nil {
		return err
	}

	if options.jsonOutput {
		fmt.Printf("%s\n%s\n", resp.Status, string(respPayload))
	} else {
		if options.verbose {
			fmt.Printf("http response status code: %s\nnew plugin ID: %s\n", resp.Status, consumerRateLimitingResponse.Id)
		} else {
			fmt.Printf("%s\n", consumerRateLimitingResponse.Id)
		}
	}

	return nil
}
