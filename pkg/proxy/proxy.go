/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/projectriff/streaming-http-adapter/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	XHttpMethodHeader = http.CanonicalHeaderKey("x-http-method")
	XHttpPathHeader   = http.CanonicalHeaderKey("x-http-path")
	XHttpQueryHeader  = http.CanonicalHeaderKey("x-http-query")
	XHttpProtoHeader  = http.CanonicalHeaderKey("x-http-proto")
	XHttpStatusHeader = http.CanonicalHeaderKey("x-http-status")
)

type proxy struct {
	server      *http.Server
	riffClient  rpc.RiffClient
	grpcAddress string
}

func NewProxy(grpcAddress string, httpAddress string) (*proxy, error) {

	p := proxy{grpcAddress: grpcAddress}

	m := http.NewServeMux()
	m.HandleFunc("/", p.invokeGrpc)

	p.server = &http.Server{
		Addr:    httpAddress,
		Handler: m,
	}

	return &p, nil
}

func (p *proxy) Run() error {

	timeout, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	conn, err := grpc.DialContext(timeout, p.grpcAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	p.riffClient = rpc.NewRiffClient(conn)

	err = p.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	} else {
		return nil
	}
}

func (p *proxy) Shutdown(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func (p *proxy) invokeGrpc(writer http.ResponseWriter, request *http.Request) {
	// TODO relax these restriction now that we expose more http semantics to functions
	if request.Method != http.MethodPost || request.URL.Path != "/" {
		writer.WriteHeader(http.StatusNotImplemented)
		return
	}
	client, err := p.riffClient.Invoke(context.Background())
	if err != nil {
		writeError(writer, err)
		return
	}

	accept := request.Header.Get("accept")
	if accept == "" {
		accept = "*/*"
	}
	contentType := request.Header.Get("content-type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	startSignal := rpc.InputSignal{
		Frame: &rpc.InputSignal_Start{
			Start: &rpc.StartFrame{
				ExpectedContentTypes: []string{accept},
				InputNames:           []string{"in"},
				OutputNames:          []string{"out"},
			},
		},
	}
	if err := client.Send(&startSignal); err != nil {
		writeError(writer, err)
		return
	}

	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writeError(writer, err)
		return
	}
	inputFrame := rpc.InputFrame{
		ContentType: contentType,
		ArgIndex:    0,
		Payload:     bytes,
		Headers:     make(map[string]string, len(request.Header)+4),
	}
	for h, v := range request.Header {
		inputFrame.Headers[h] = v[0]
	}
	inputFrame.Headers[XHttpMethodHeader] = request.Method
	inputFrame.Headers[XHttpPathHeader] = request.URL.Path
	inputFrame.Headers[XHttpQueryHeader] = request.URL.RawQuery
	inputFrame.Headers[XHttpProtoHeader] = request.Proto
	dataSignal := rpc.InputSignal{
		Frame: &rpc.InputSignal_Data{
			Data: &inputFrame,
		},
	}
	if err := client.Send(&dataSignal); err != nil {
		writeError(writer, err)
		return
	}
	if err := client.CloseSend(); err != nil {
		writeError(writer, err)
		return
	}

	outputSignal, err := client.Recv()
	if err != nil {
		writeError(writer, err)
		return
	}
	if _, err := client.Recv(); err != io.EOF {
		writeError(writer, errors.New("expected EOF"))
		return
	}
	if status, ok := outputSignal.GetData().Headers[XHttpStatusHeader]; ok {
		code, err := strconv.Atoi(status)
		if err != nil {
			writeError(writer, fmt.Errorf("invalid status code %q", status))
			return
		}
		writer.WriteHeader(code)
	}
	writer.Header().Set("content-type", outputSignal.GetData().ContentType)
	for h, v := range outputSignal.GetData().Headers {
		writer.Header().Set(h, v)
	}
	if _, err = writer.Write(outputSignal.GetData().Payload); err != nil {
		fmt.Printf("unable to write proxy response: %s\n", err)
		return
	}
}

func writeError(writer http.ResponseWriter, err error) {
	if grpcError, ok := status.FromError(err); ok {
		writeHeaderFromGrpcError(grpcError, writer)
		writer.Header().Set("content-type", "text/plain")
		_, _ = writer.Write([]byte(grpcError.Message()))
		_, _ = writer.Write([]byte("\n"))
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("content-type", "text/plain")
		_, _ = writer.Write([]byte(err.Error()))
		_, _ = writer.Write([]byte("\n"))
	}

}

func writeHeaderFromGrpcError(grpcError *status.Status, writer http.ResponseWriter) {
	if grpcError.Code() != codes.InvalidArgument {
		writer.WriteHeader(http.StatusInternalServerError)
	} else if strings.HasPrefix(grpcError.Message(), "Invoker: Unsupported Media Type") {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
	} else if strings.HasPrefix(grpcError.Message(), "Invoker: Not Acceptable") {
		writer.WriteHeader(http.StatusNotAcceptable)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
