package proxy

import (
	"context"
	"errors"
	"github.com/projectriff/streaming-http-adapter/pkg/rpc"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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
	if request.Method != http.MethodPost {
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
		accept = "application/octet-stream"
	}
	contentType := request.Header.Get("content-type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	startSignal := rpc.InputSignal{
		Frame: &rpc.InputSignal_Start{
			Start: &rpc.StartFrame{
				ExpectedContentTypes: []string{accept},
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
		Headers:     make(map[string]string, len(request.Header)),
	}
	for h, v := range request.Header {
		inputFrame.Headers[h] = v[0]
	}
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
	if _, err = writer.Write(outputSignal.GetData().Payload); err != nil {
		writeError(writer, err)
		return
	}
	writer.Header().Set("content-type", outputSignal.GetData().ContentType)
	for h, v := range outputSignal.GetData().Headers {
		writer.Header().Set(h, v)
	}
}

func writeError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Header().Set("content-type", "text/plain")
	_, _ = writer.Write([]byte(err.Error()))
}
