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

// BindMFADevice invokes the ram.BindMFADevice API synchronously
// api document: https://help.aliyun.com/api/ram/bindmfadevice.html
func (client *Client) BindMFADevice(request *BindMFADeviceRequest) (response *BindMFADeviceResponse, err error) {
	response = CreateBindMFADeviceResponse()
	err = client.DoAction(request, response)
	return
}

// BindMFADeviceWithChan invokes the ram.BindMFADevice API asynchronously
// api document: https://help.aliyun.com/api/ram/bindmfadevice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) BindMFADeviceWithChan(request *BindMFADeviceRequest) (<-chan *BindMFADeviceResponse, <-chan error) {
	responseChan := make(chan *BindMFADeviceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.BindMFADevice(request)
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

// BindMFADeviceWithCallback invokes the ram.BindMFADevice API asynchronously
// api document: https://help.aliyun.com/api/ram/bindmfadevice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) BindMFADeviceWithCallback(request *BindMFADeviceRequest, callback func(response *BindMFADeviceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *BindMFADeviceResponse
		var err error
		defer close(result)
		response, err = client.BindMFADevice(request)
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

// BindMFADeviceRequest is the request struct for api BindMFADevice
type BindMFADeviceRequest struct {
	*requests.RpcRequest
	SerialNumber        string `position:"Query" name:"SerialNumber"`
	AuthenticationCode2 string `position:"Query" name:"AuthenticationCode2"`
	AuthenticationCode1 string `position:"Query" name:"AuthenticationCode1"`
	UserName            string `position:"Query" name:"UserName"`
}

// BindMFADeviceResponse is the response struct for api BindMFADevice
type BindMFADeviceResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateBindMFADeviceRequest creates a request to invoke BindMFADevice API
func CreateBindMFADeviceRequest() (request *BindMFADeviceRequest) {
	request = &BindMFADeviceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ram", "2015-05-01", "BindMFADevice", "ram", "openAPI")
	return
}

// CreateBindMFADeviceResponse creates a response to parse from BindMFADevice response
func CreateBindMFADeviceResponse() (response *BindMFADeviceResponse) {
	response = &BindMFADeviceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
