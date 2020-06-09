package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/projectriff/streaming-http-adapter/pkg/proxy/mocks"
	"github.com/projectriff/streaming-http-adapter/pkg/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_invokeGrpc_input_startFrame(t *testing.T) {
	riffClient, invokeClient := mockRiffClient()
	p := &proxy{riffClient: riffClient}

	request, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	request.Header.Add("accept", "text/plain")
	p.invokeGrpc(httptest.NewRecorder(), request)

	inputSignals := inputSignals(invokeClient.Calls)
	startFrame := inputSignals[0].GetStart()
	assert.Equal(t, []string{"text/plain"}, startFrame.ExpectedContentTypes)
	assert.Equal(t, []string{"in"}, startFrame.InputNames)
	assert.Equal(t, []string{"out"}, startFrame.OutputNames)
}

func Test_invokeGrpc_input_dataFrame(t *testing.T) {
	riffClient, invokeClient := mockRiffClient()
	p := &proxy{riffClient: riffClient}

	request, _ := http.NewRequest("POST", "/", strings.NewReader("some body"))
	request.Header.Add("content-type", "text/plain")
	request.Header.Add("x-custom-header", "header-value")
	p.invokeGrpc(httptest.NewRecorder(), request)

	inputSignals := inputSignals(invokeClient.Calls)
	dataFrame := inputSignals[1].GetData()
	assert.Equal(t, dataFrame.Headers[XHttpMethodHeader], "POST")
	assert.Equal(t, dataFrame.Headers[XHttpPathHeader], "/")
	assert.Equal(t, dataFrame.Headers[XHttpQueryHeader], "")
	assert.Equal(t, dataFrame.Headers[XHttpProtoHeader], "HTTP/1.1")
	assert.Equal(t, "some body", string(dataFrame.Payload))
	assert.Equal(t, "text/plain", dataFrame.ContentType)
	assert.Contains(t, dataFrame.Headers, "X-Custom-Header")
	assert.Equal(t, dataFrame.Headers["X-Custom-Header"], "header-value")
}

func Test_invokeGrpc_output(t *testing.T) {
	riffClient, _ := mockRiffClientWithResponse("<data>some response</data>", "application/xml", map[string]string{
		XHttpStatusHeader: fmt.Sprintf("%d", http.StatusCreated),
	})
	p := &proxy{riffClient: riffClient}

	request, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, "<data>some response</data>", responseRecorder.Body.String())
	assert.Equal(t, "application/xml", responseRecorder.Header().Get("Content-Type"))
}

func Test_invokeGrpc_wiring(t *testing.T) {
	riffClient, invokeClient := mockRiffClient()
	p := &proxy{riffClient: riffClient}

	request, _ := http.NewRequest("POST", "/", strings.NewReader("some body"))
	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	riffClient.AssertExpectations(t)
	invokeClient.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func Test_not_acceptable_media_type(t *testing.T) {
	accept := "text/zglorbf"
	errorMsg := fmt.Sprintf("Invoker: Not Acceptable: unrecognized output #0's content-type %s", accept)
	riffClient, _ := mockRiffClientWithError(codes.InvalidArgument, errorMsg)
	p := &proxy{riffClient: riffClient}
	request, _ := http.NewRequest("POST", "/", strings.NewReader("some body"))
	request.Header.Set("Accept", accept)

	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusNotAcceptable, responseRecorder.Code)
	assert.Equal(t, "text/plain", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, errorMsg+"\n", responseRecorder.Body.String())
}

func Test_unsupported_request_method(t *testing.T) {
	riffClient, _ := mockRiffClient()
	p := &proxy{riffClient: riffClient}

	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
}

func Test_unsupported_request_path(t *testing.T) {
	riffClient, _ := mockRiffClient()
	p := &proxy{riffClient: riffClient}

	request, _ := http.NewRequest("POST", "/nope/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
}

func Test_unsupported_content_type(t *testing.T) {
	contentType := "text/zglorbf"
	errorMsg := fmt.Sprintf("Invoker: Unsupported Media Type: unsupported input #0's content-type %s", contentType)
	riffClient, _ := mockRiffClientWithError(codes.InvalidArgument, errorMsg)
	p := &proxy{riffClient: riffClient}
	request, _ := http.NewRequest("POST", "/", strings.NewReader("some body"))
	request.Header.Set("Content-Type", contentType)

	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusUnsupportedMediaType, responseRecorder.Code)
	assert.Equal(t, "text/plain", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, errorMsg+"\n", responseRecorder.Body.String())
}

func Test_unsupported_content_type_json(t *testing.T) {
	contentType := "text/zglorbf"
	errorMsg := fmt.Sprintf("Invoker: Unsupported Media Type: unsupported input #0's content-type %s", contentType)
	riffClient, _ := mockRiffClientWithError(codes.InvalidArgument, errorMsg)
	p := &proxy{riffClient: riffClient}
	request, _ := http.NewRequest("POST", "/", strings.NewReader("some body"))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "application/json")

	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusUnsupportedMediaType, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"Invoker: Unsupported Media Type: unsupported input #0's content-type text/zglorbf\"}", responseRecorder.Body.String())
}

func Test_error_ordered_accept_text(t *testing.T) {
	contentType := "text/zglorbf"
	errorMsg := fmt.Sprintf("Invoker: Unsupported Media Type: unsupported input #0's content-type %s", contentType)
	riffClient, _ := mockRiffClientWithError(codes.InvalidArgument, errorMsg)
	p := &proxy{riffClient: riffClient}
	request, _ := http.NewRequest("POST", "/", strings.NewReader("some body"))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "text/plain, application/json")

	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusUnsupportedMediaType, responseRecorder.Code)
	assert.Equal(t, "text/plain", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "Invoker: Unsupported Media Type: unsupported input #0's content-type text/zglorbf\n", responseRecorder.Body.String())
}

func Test_error_weighted_accept_json(t *testing.T) {
	contentType := "text/zglorbf"
	errorMsg := fmt.Sprintf("Invoker: Unsupported Media Type: unsupported input #0's content-type %s", contentType)
	riffClient, _ := mockRiffClientWithError(codes.InvalidArgument, errorMsg)
	p := &proxy{riffClient: riffClient}
	request, _ := http.NewRequest("POST", "/", strings.NewReader("some body"))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "something/madeup;q=0.5, application/json; charset=utf-8, text/plain")

	responseRecorder := httptest.NewRecorder()
	p.invokeGrpc(responseRecorder, request)

	assert.Equal(t, http.StatusUnsupportedMediaType, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"Invoker: Unsupported Media Type: unsupported input #0's content-type text/zglorbf\"}", responseRecorder.Body.String())
}

func inputSignals(calls []mock.Call) []*rpc.InputSignal {
	var inputSignals []*rpc.InputSignal
	for _, call := range calls {
		if call.Method == "Send" {
			signal := call.Arguments.Get(0).(*rpc.InputSignal)
			inputSignals = append(inputSignals, signal)
		}
	}
	return inputSignals
}

func mockRiffClient() (*mocks.RiffClient, *mocks.Riff_InvokeClient) {
	return mockRiffClientWithResponse("", "", map[string]string{})
}

func mockRiffClientWithResponse(outputBody string, contentType string, headers map[string]string) (*mocks.RiffClient, *mocks.Riff_InvokeClient) {
	riffClient := &mocks.RiffClient{}
	invokeClient := &mocks.Riff_InvokeClient{}
	riffClient.On("Invoke", context.Background()).Return(invokeClient, nil)
	invokeClient.On("Send", mock.Anything).Return(nil)
	invokeClient.On("CloseSend").Return(nil)
	invokeClient.On("Recv").Return(outputSignal(outputBody, contentType, headers), nil).Once()
	invokeClient.On("Recv").Return(nil, io.EOF)
	return riffClient, invokeClient
}

func mockRiffClientWithError(code codes.Code, msg string) (*mocks.RiffClient, *mocks.Riff_InvokeClient) {
	riffClient := &mocks.RiffClient{}
	invokeClient := &mocks.Riff_InvokeClient{}
	riffClient.On("Invoke", context.Background()).Return(invokeClient, nil)
	invokeClient.On("Send", mock.MatchedBy(isStartSignal)).Return(nil)
	invokeClient.On("Send", mock.MatchedBy(isDataSignal)).Return(status.Error(code, msg))
	return riffClient, invokeClient
}

func isDataSignal(inputSignal *rpc.InputSignal) bool {
	return inputSignal.GetData() != nil
}

func isStartSignal(inputSignal *rpc.InputSignal) bool {
	return inputSignal.GetStart() != nil
}

func outputSignal(outputBody string, contentType string, headers map[string]string) *rpc.OutputSignal {
	return &rpc.OutputSignal{
		Frame: &rpc.OutputSignal_Data{
			Data: &rpc.OutputFrame{
				Payload:     []byte(outputBody),
				ContentType: contentType,
				Headers:     headers,
			},
		},
	}
}
