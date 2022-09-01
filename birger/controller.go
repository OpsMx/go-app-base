/*
 * Copyright 2022 OpsMx, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package birger

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/OpsMx/go-app-base/httputil"
	"github.com/OpsMx/go-app-base/util"
)

// ControllerManager checks the services available on the controller,
// and fetches new tokens for newly discovered services.  It will
// update the ArgoManager with new endpoints, and remove old ones.
type ControllerManager struct {
	UpdateChan        chan ServiceUpdate
	conf              Config
	serviceTypes      []string
	shutdownWorker    chan bool
	shutdownCount     sync.WaitGroup
	updateRate        time.Duration
	healthcheckStatus error
	services          map[string]controllerService
}

type controllerService struct {
	URL         string
	Name        string
	Type        string
	Annotations map[string]string
	AgentName   string
	Token       string
}

// MakeControllerManager returns a new ControllerManager which will periodically poll
// the controller for services, and send
func MakeControllerManager(conf Config, serviceTypes []string) *ControllerManager {
	conf.applyDefaults()
	m := ControllerManager{
		conf:              conf,
		serviceTypes:      serviceTypes,
		shutdownWorker:    make(chan bool),
		updateRate:        time.Duration(conf.UpdateFrequencySeconds) * time.Second,
		services:          map[string]controllerService{},
		healthcheckStatus: fmt.Errorf("controller is not yet synced"),
		UpdateChan:        make(chan ServiceUpdate, 10),
	}

	m.shutdownCount.Add(1)
	go m.worker()

	return &m
}

// Shutdown tells the manager to stop doing updates and causes all
// goprocs started to exit as cleanly as possible.
func (m *ControllerManager) Shutdown() {
	m.shutdownWorker <- true
	close(m.UpdateChan)
	m.shutdownCount.Wait()
}

func (m *ControllerManager) worker() {
	// Initialize but stop the timer before it triggers.
	t := time.NewTimer(1 * time.Hour)
	t.Stop()

	m.reloadFromController()
	t.Reset(m.updateRate)

	for {
		select {
		case <-m.shutdownWorker:
			t.Stop()
			m.shutdownCount.Done()
			return
		case <-t.C:
			m.reloadFromController()
			t.Reset(m.updateRate)
		}
	}
}

func (m *ControllerManager) reloadFromController() {
	services, err := m.getArgoServices()
	if err != nil {
		m.healthcheckStatus = err
		return
	}
	m.healthcheckStatus = nil

	// compare existing services to the new list.  We can assume that if we have an entry,
	// we do not need to refresh tokens and the URL cannot change when talking to the
	// controller.  If these change, we will want a restart.
	for key, fetchedService := range services {
		if svc, found := m.services[key]; found {
			if annotationsDifferent(svc, fetchedService) {
				m.services[key] = fetchedService
				m.sendUpdate(fetchedService)
			}
			continue
		}
		url, token, err := m.getTokenAndURL(fetchedService)
		if err != nil {
			m.healthcheckStatus = err
			log.Printf("unable to fetch service credentials from controller: %v", err)
			return
		}
		fetchedService.URL = url
		fetchedService.Token = token
		m.services[key] = fetchedService
		m.sendUpdate(fetchedService)
	}

	// now, remove any we don't currently see.
	for key, service := range m.services {
		if _, found := services[key]; found {
			continue
		}
		m.sendDelete(service)
		delete(m.services, key)
	}
}

func annotationsDifferent(a controllerService, b controllerService) bool {
	if len(a) != len(b) {
		return true
	}
	for k, v := range a.Annotations {
		if v != b.Annotations[k] {
			return true
		}
	}
	return false
}

func (m *ControllerManager) sendUpdate(s controllerService) {
	m.UpdateChan <- ServiceUpdate{
		Operation:   "update",
		Name:        s.Name,
		Type:        s.Type,
		AgentName:   s.AgentName,
		Annotations: s.Annotations,
		URL:         s.URL,
		Token:       s.Token,
	}
}

func (m *ControllerManager) sendDelete(s controllerService) {
	m.UpdateChan <- ServiceUpdate{
		Operation: "delete",
		Name:      s.Name,
		Type:      s.Type,
		AgentName: s.AgentName,
	}
}

// Check returns the last error received during a sync, if any.
// Used for a healthcheck status.
func (m *ControllerManager) Check() error {
	return m.healthcheckStatus
}

type connectedAgentsResponse struct {
	ConnectedAgents []connectedAgent `json:"connectedAgents,omitempty"`
}

type connectedAgent struct {
	Name         string            `json:"name,omitempty"`
	Annnotations map[string]string `json:"annotations,omitempty"`
	Endpoints    []agentEndpoint   `json:"endpoints,omitempty"`
}

type agentEndpoint struct {
	Name         string            `json:"name,omitempty"`
	Type         string            `json:"type,omitempty"`
	Annnotations map[string]string `json:"annotations,omitempty"`
	Configured   bool              `json:"configured,omitempty"`
}

type controllerServiceCredentialsRequest struct {
	AgentName string `json:"agentName,omitempty"`
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
}

type controllerServiceCredentialResponse struct {
	AgentName      string `json:"agentName,omitempty"`
	Name           string `json:"name,omitempty"`
	Type           string `json:"type,omitempty"`
	CredentialType string `json:"credentialType,omitempty"`
	Credential     struct {
		Password string `json:"password,omitempty"`
	} `json:"credential,omitempty"`
	URL string `json:"url,omitempty"`
}

func (m *ControllerManager) getTokenAndURL(s controllerService) (serviceUrl string, serviceToken string, err error) {
	url, err := url.JoinPath(m.conf.URL, "/api/v1/generateServiceCredentials")

	client, err := m.getTLSClient()
	if err != nil {
		return "", "", fmt.Errorf("making TLS client: %v", err)
	}

	credentialsRequest := controllerServiceCredentialsRequest{
		AgentName: s.AgentName,
		Name:      s.Name,
		Type:      s.Type,
	}

	d, err := json.Marshal(credentialsRequest)
	r := bytes.NewReader(d)
	resp, err := client.Post(url, "application/json", r)
	if err != nil {
		return "", "", fmt.Errorf("fetching service credentials: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("fetching service credentials: http status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("reading body: %v", err)
	}

	var creds controllerServiceCredentialResponse
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return "", "", fmt.Errorf("cannot decode service credentials JSON: %v", err)
	}

	return creds.URL, creds.Credential.Password, nil
}

func (m *ControllerManager) getArgoServices() (map[string]controllerService, error) {
	url, err := url.JoinPath(m.conf.URL, "/api/v1/getAgentStatistics")
	if err != nil {
		return map[string]controllerService{}, fmt.Errorf("joining url: %v", err)
	}

	client, err := m.getTLSClient()
	if err != nil {
		return map[string]controllerService{}, fmt.Errorf("making TLS client: %v", err)
	}

	resp, err := client.Get(url)
	if err != nil {
		return map[string]controllerService{}, fmt.Errorf("fetching connected agents: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return map[string]controllerService{}, fmt.Errorf("fetching connnected agents: http status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]controllerService{}, fmt.Errorf("reading body: %v", err)
	}

	return m.parseAgentStatistics(data)
}

func (m *ControllerManager) parseAgentStatistics(data []byte) (map[string]controllerService, error) {
	var ca connectedAgentsResponse
	err := json.Unmarshal(data, &ca)
	if err != nil {
		return map[string]controllerService{}, fmt.Errorf("cannot decode connected agent JSON: %v", err)
	}

	endpoints := map[string]controllerService{}

	for _, a := range ca.ConnectedAgents {
		for _, ep := range a.Endpoints {
			if !ep.Configured || !util.Contains(m.serviceTypes, ep.Type) {
				continue
			}
			key := a.Name + ":" + ep.Name + ":" + ep.Type
			endpoints[key] = controllerService{AgentName: a.Name, Name: ep.Name, Type: ep.Type, Annotations: ep.Annnotations}
		}
	}

	return endpoints, nil
}

func (m *ControllerManager) getTLSClient() (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(m.conf.CertificatePath, m.conf.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("loading certificate and key: %v", err)
	}

	caCert, err := os.ReadFile(m.conf.CAPath)
	if err != nil {
		return nil, fmt.Errorf("loading certificate authority: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return httputil.NewHTTPClient(&tlsConfig), nil
}
