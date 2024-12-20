package vpc

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

// ModifyExpressConnectTrafficQos invokes the vpc.ModifyExpressConnectTrafficQos API synchronously
func (client *Client) ModifyExpressConnectTrafficQos(request *ModifyExpressConnectTrafficQosRequest) (response *ModifyExpressConnectTrafficQosResponse, err error) {
	response = CreateModifyExpressConnectTrafficQosResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyExpressConnectTrafficQosWithChan invokes the vpc.ModifyExpressConnectTrafficQos API asynchronously
func (client *Client) ModifyExpressConnectTrafficQosWithChan(request *ModifyExpressConnectTrafficQosRequest) (<-chan *ModifyExpressConnectTrafficQosResponse, <-chan error) {
	responseChan := make(chan *ModifyExpressConnectTrafficQosResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyExpressConnectTrafficQos(request)
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

// ModifyExpressConnectTrafficQosWithCallback invokes the vpc.ModifyExpressConnectTrafficQos API asynchronously
func (client *Client) ModifyExpressConnectTrafficQosWithCallback(request *ModifyExpressConnectTrafficQosRequest, callback func(response *ModifyExpressConnectTrafficQosResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyExpressConnectTrafficQosResponse
		var err error
		defer close(result)
		response, err = client.ModifyExpressConnectTrafficQos(request)
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

// ModifyExpressConnectTrafficQosRequest is the request struct for api ModifyExpressConnectTrafficQos
type ModifyExpressConnectTrafficQosRequest struct {
	*requests.RpcRequest
	ClientToken          string                                              `position:"Query" name:"ClientToken"`
	RemoveInstanceList   *[]ModifyExpressConnectTrafficQosRemoveInstanceList `position:"Query" name:"RemoveInstanceList"  type:"Repeated"`
	AddInstanceList      *[]ModifyExpressConnectTrafficQosAddInstanceList    `position:"Query" name:"AddInstanceList"  type:"Repeated"`
	QosId                string                                              `position:"Query" name:"QosId"`
	ResourceOwnerAccount string                                              `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string                                              `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer                                    `position:"Query" name:"OwnerId"`
	QosName              string                                              `position:"Query" name:"QosName"`
	QosDescription       string                                              `position:"Query" name:"QosDescription"`
}

// ModifyExpressConnectTrafficQosRemoveInstanceList is a repeated param struct in ModifyExpressConnectTrafficQosRequest
type ModifyExpressConnectTrafficQosRemoveInstanceList struct {
	InstanceId   string `name:"InstanceId"`
	InstanceType string `name:"InstanceType"`
}

// ModifyExpressConnectTrafficQosAddInstanceList is a repeated param struct in ModifyExpressConnectTrafficQosRequest
type ModifyExpressConnectTrafficQosAddInstanceList struct {
	InstanceId   string `name:"InstanceId"`
	InstanceType string `name:"InstanceType"`
}

// ModifyExpressConnectTrafficQosResponse is the response struct for api ModifyExpressConnectTrafficQos
type ModifyExpressConnectTrafficQosResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyExpressConnectTrafficQosRequest creates a request to invoke ModifyExpressConnectTrafficQos API
func CreateModifyExpressConnectTrafficQosRequest() (request *ModifyExpressConnectTrafficQosRequest) {
	request = &ModifyExpressConnectTrafficQosRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ModifyExpressConnectTrafficQos", "vpc", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyExpressConnectTrafficQosResponse creates a response to parse from ModifyExpressConnectTrafficQos response
func CreateModifyExpressConnectTrafficQosResponse() (response *ModifyExpressConnectTrafficQosResponse) {
	response = &ModifyExpressConnectTrafficQosResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
