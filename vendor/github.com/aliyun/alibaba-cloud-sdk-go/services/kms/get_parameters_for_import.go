package kms

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// GetParametersForImport invokes the kms.GetParametersForImport API synchronously
// api document: https://help.aliyun.com/api/kms/getparametersforimport.html
func (client *Client) GetParametersForImport(request *GetParametersForImportRequest) (response *GetParametersForImportResponse, err error) {
	response = CreateGetParametersForImportResponse()
	err = client.DoAction(request, response)
	return
}

// GetParametersForImportWithChan invokes the kms.GetParametersForImport API asynchronously
// api document: https://help.aliyun.com/api/kms/getparametersforimport.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetParametersForImportWithChan(request *GetParametersForImportRequest) (<-chan *GetParametersForImportResponse, <-chan error) {
	responseChan := make(chan *GetParametersForImportResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetParametersForImport(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// GetParametersForImportWithCallback invokes the kms.GetParametersForImport API asynchronously
// api document: https://help.aliyun.com/api/kms/getparametersforimport.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetParametersForImportWithCallback(request *GetParametersForImportRequest, callback func(response *GetParametersForImportResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetParametersForImportResponse
		var err error
		defer close(result)
		response, err = client.GetParametersForImport(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// GetParametersForImportRequest is the request struct for api GetParametersForImport
type GetParametersForImportRequest struct {
	*requests.RpcRequest
	KeyId             string `position:"Query" name:"KeyId"`
	STSToken          string `position:"Query" name:"STSToken"`
	WrappingAlgorithm string `position:"Query" name:"WrappingAlgorithm"`
	WrappingKeySpec   string `position:"Query" name:"WrappingKeySpec"`
}

// GetParametersForImportResponse is the response struct for api GetParametersForImport
type GetParametersForImportResponse struct {
	*responses.BaseResponse
	KeyId           string `json:"KeyId" xml:"KeyId"`
	RequestId       string `json:"RequestId" xml:"RequestId"`
	ImportToken     string `json:"ImportToken" xml:"ImportToken"`
	PublicKey       string `json:"PublicKey" xml:"PublicKey"`
	TokenExpireTime string `json:"TokenExpireTime" xml:"TokenExpireTime"`
}

// CreateGetParametersForImportRequest creates a request to invoke GetParametersForImport API
func CreateGetParametersForImportRequest() (request *GetParametersForImportRequest) {
	request = &GetParametersForImportRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Kms", "2016-01-20", "GetParametersForImport", "kms", "openAPI")
	return
}

// CreateGetParametersForImportResponse creates a response to parse from GetParametersForImport response
func CreateGetParametersForImportResponse() (response *GetParametersForImportResponse) {
	response = &GetParametersForImportResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
