// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"net/http"
	"net/url"
)

type CLIHandler struct{}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=oci authType=ApiKey 
       vault login -method=oci authType=InstancePrincipal 

  The OCI auth method allows users to authenticate with OCI
  credentials. The OCI credentials may be specified in a number of ways,
  listed below:

    1. API Key

    2. Instance Principal

  Authenticate using API key:

		First create a configuration file as explained in https://docs.us-phoenix-1.oraclecloud.com/Content/API/Concepts/sdkconfig.htm
		Then login using the following command:

		$ vault login -method=oci authType=ApiKey role=<RoleName>

  Authenticate using instance principal:
		https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/callingservicesfrominstances.htm
		
		$ vault login -method=oci authType=InstancePrincipal role=<RoleName>

Configuration:
  authType=<string>
      Enter one of following: 
		ApiKey (or) ak		
		InstancePrincipal (or) ip
`
	return strings.TrimSpace(help)
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	role, ok := m["role"]
	if !ok {
		return nil, fmt.Errorf("Enter the role")
	}
	role = strings.ToLower(role)

	path := fmt.Sprintf(PATH_BASE_FORMAT, role)
	signingPath := PATH_VERSION_BASE + path

	loginData, err := CreateLoginData(c.Address(), m, signingPath)
	if err != nil {
		return nil, err
	}

	//Now try to login
	secret, err := c.Logical().Write(path, loginData)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func CreateLoginData(clientAddress string, m map[string]string, path string) (map[string]interface{}, error) {

	authtype, ok := m["authType"]
	if !ok {
		return nil, fmt.Errorf("Enter the authType")
	}

	switch strings.ToLower(authtype) {
	case "ip", "instanceprincipal":
		return createLoginDataForInstancePrincipal(clientAddress, path)
	case "ak", "apikey":
		return createLoginDataForApiKeys(clientAddress, path)
	}

	return nil, fmt.Errorf("Unknown authType")
}

func createLoginDataForApiKeys(clientAddress string, path string) (map[string]interface{}, error) {

	provider := common.DefaultConfigProvider()

	ociClient, err := NewOciClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	return createFinalLoginData(clientAddress, &ociClient, path)
}

func createLoginDataForInstancePrincipal(clientAddress string, path string) (map[string]interface{}, error) {

	ip, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		return nil, err
	}
	ociClient, err := NewOciClientWithConfigurationProvider(ip)
	if err != nil {
		return nil, err
	}
	return createFinalLoginData(clientAddress, &ociClient, path)
}

func createFinalLoginData(clientAddress string, ociClient *OciClient, path string) (map[string]interface{}, error) {

	ociClient.SetHost(clientAddress)
	request, err := ociClient.ConstructLoginRequest(path)
	if err != nil {
		return nil, err
	}

	clientURL, err := url.Parse(clientAddress)
	if err != nil {
		return nil, err
	}
	request.Host = clientURL.Host

	//serialize the request
	serializedRequest := serializeRequest(request)

	//pack it into loginData
	loginData := make(map[string]interface{})
	loginData["requestHeaders"] = serializedRequest

	return loginData, nil
}

func serializeRequest(request http.Request) (map[string][]string) {
	requestHeaders := request.Header
	requestHeaders["host"] = []string{request.Host}
	requestHeaders["(request-target)"] = []string{getRequestTarget(&request)}
	return requestHeaders
}
