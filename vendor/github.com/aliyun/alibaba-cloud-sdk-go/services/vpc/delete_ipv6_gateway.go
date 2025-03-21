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

// DeleteIpv6Gateway invokes the vpc.DeleteIpv6Gateway API synchronously
func (client *Client) DeleteIpv6Gateway(request *DeleteIpv6GatewayRequest) (response *DeleteIpv6GatewayResponse, err error) {
	response = CreateDeleteIpv6GatewayResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteIpv6GatewayWithChan invokes the vpc.DeleteIpv6Gateway API asynchronously
func (client *Client) DeleteIpv6GatewayWithChan(request *DeleteIpv6GatewayRequest) (<-chan *DeleteIpv6GatewayResponse, <-chan error) {
	responseChan := make(chan *DeleteIpv6GatewayResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteIpv6Gateway(request)
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

// DeleteIpv6GatewayWithCallback invokes the vpc.DeleteIpv6Gateway API asynchronously
func (client *Client) DeleteIpv6GatewayWithCallback(request *DeleteIpv6GatewayRequest, callback func(response *DeleteIpv6GatewayResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteIpv6GatewayResponse
		var err error
		defer close(result)
		response, err = client.DeleteIpv6Gateway(request)
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

// DeleteIpv6GatewayRequest is the request struct for api DeleteIpv6Gateway
type DeleteIpv6GatewayRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	DryRun               requests.Boolean `position:"Query" name:"DryRun"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	Ipv6GatewayId        string           `position:"Query" name:"Ipv6GatewayId"`
}

// DeleteIpv6GatewayResponse is the response struct for api DeleteIpv6Gateway
type DeleteIpv6GatewayResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteIpv6GatewayRequest creates a request to invoke DeleteIpv6Gateway API
func CreateDeleteIpv6GatewayRequest() (request *DeleteIpv6GatewayRequest) {
	request = &DeleteIpv6GatewayRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "DeleteIpv6Gateway", "vpc", "openAPI")
	request.Method = requests.POST
	return
}

// CreateDeleteIpv6GatewayResponse creates a response to parse from DeleteIpv6Gateway response
func CreateDeleteIpv6GatewayResponse() (response *DeleteIpv6GatewayResponse) {
	response = &DeleteIpv6GatewayResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
