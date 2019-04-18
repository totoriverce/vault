package ram

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

// GetPublicKey invokes the ram.GetPublicKey API synchronously
// api document: https://help.aliyun.com/api/ram/getpublickey.html
func (client *Client) GetPublicKey(request *GetPublicKeyRequest) (response *GetPublicKeyResponse, err error) {
	response = CreateGetPublicKeyResponse()
	err = client.DoAction(request, response)
	return
}

// GetPublicKeyWithChan invokes the ram.GetPublicKey API asynchronously
// api document: https://help.aliyun.com/api/ram/getpublickey.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetPublicKeyWithChan(request *GetPublicKeyRequest) (<-chan *GetPublicKeyResponse, <-chan error) {
	responseChan := make(chan *GetPublicKeyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetPublicKey(request)
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

// GetPublicKeyWithCallback invokes the ram.GetPublicKey API asynchronously
// api document: https://help.aliyun.com/api/ram/getpublickey.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetPublicKeyWithCallback(request *GetPublicKeyRequest, callback func(response *GetPublicKeyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetPublicKeyResponse
		var err error
		defer close(result)
		response, err = client.GetPublicKey(request)
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

// GetPublicKeyRequest is the request struct for api GetPublicKey
type GetPublicKeyRequest struct {
	*requests.RpcRequest
	UserPublicKeyId string `position:"Query" name:"UserPublicKeyId"`
	UserName        string `position:"Query" name:"UserName"`
}

// GetPublicKeyResponse is the response struct for api GetPublicKey
type GetPublicKeyResponse struct {
	*responses.BaseResponse
	RequestId string    `json:"RequestId" xml:"RequestId"`
	PublicKey PublicKey `json:"PublicKey" xml:"PublicKey"`
}

// CreateGetPublicKeyRequest creates a request to invoke GetPublicKey API
func CreateGetPublicKeyRequest() (request *GetPublicKeyRequest) {
	request = &GetPublicKeyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "GetPublicKey", "ram", "openAPI")
	return
}

// CreateGetPublicKeyResponse creates a response to parse from GetPublicKey response
func CreateGetPublicKeyResponse() (response *GetPublicKeyResponse) {
	response = &GetPublicKeyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
