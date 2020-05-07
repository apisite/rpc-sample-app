// This code was autogenerated from main.proto, do not edit.
package nrpcgen

import (
	"context"
	"log"
	"time"

	"github.com/TenderPro/rpckit/app/ticker"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-rpc/nrpc"
	"github.com/prometheus/client_golang/prometheus"
	grpc "google.golang.org/grpc"

	opentracing "github.com/opentracing/opentracing-go"

	"SELF/pkg/pb"
)

// PingServiceServer is the interface that providers of the service
// PingService should implement.
type PingServiceServer interface {
	Ping(context.Context, *pb.PingRequest) (*pb.PingResponse, error)
	PingEmpty(context.Context, *pb.Empty) (*pb.PingResponse, error)
	PingError(context.Context, *pb.PingRequest) (*pb.Empty, error)
	PingList(*pb.PingRequest, pb.PingService_PingListServer) error
	//PingStream(pb.PingService_PingStreamServer) error
	// Время на сервере
	TimeService(*ticker.TimeRequest, pb.PingService_TimeServiceServer) error
}

var (
	// The request completion time, measured at client-side.
	clientRCTForPingService = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "nrpc_client_request_completion_time_seconds",
			Help:       "The request completion time for calls, measured client-side.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			ConstLabels: map[string]string{
				"service": "PingService",
			},
		},
		[]string{"method"})

	// The handler execution time, measured at server-side.
	serverHETForPingService = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "nrpc_server_handler_execution_time_seconds",
			Help:       "The handler execution time for calls, measured server-side.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			ConstLabels: map[string]string{
				"service": "PingService",
			},
		},
		[]string{"method"})

	// The counts of calls made by the client, classified by result type.
	clientCallsForPingService = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nrpc_client_calls_count",
			Help: "The count of calls made by the client.",
			ConstLabels: map[string]string{
				"service": "PingService",
			},
		},
		[]string{"method", "encoding", "result_type"})

	// The counts of requests handled by the server, classified by result type.
	serverRequestsForPingService = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nrpc_server_requests_count",
			Help: "The count of requests handled by the server.",
			ConstLabels: map[string]string{
				"service": "PingService",
			},
		},
		[]string{"method", "encoding", "result_type"})
)

// PingServiceHandler provides a NATS subscription handler that can serve a
// subscription using a given PingServiceServer implementation.
type PingServiceHandler struct {
	ctx     context.Context
	workers *nrpc.WorkerPool
	nc      nrpc.NatsConn
	server  PingServiceServer

	encodings []string
}

func NewPingServiceHandler(ctx context.Context, nc nrpc.NatsConn, s PingServiceServer) *PingServiceHandler {
	return &PingServiceHandler{
		ctx:    ctx,
		nc:     nc,
		server: s,

		encodings: []string{"protobuf"},
	}
}

func NewPingServiceConcurrentHandler(workers *nrpc.WorkerPool, nc nrpc.NatsConn, s PingServiceServer) *PingServiceHandler {
	return &PingServiceHandler{
		workers: workers,
		nc:      nc,
		server:  s,
	}
}

// SetEncodings sets the output encodings when using a '*Publish' function
func (h *PingServiceHandler) SetEncodings(encodings []string) {
	h.encodings = encodings
}

func (h *PingServiceHandler) Subject() string {
	return "pb.PingService.>"
}

func (h *PingServiceHandler) Handler(msg *nats.Msg) {
	var ctx context.Context
	if h.workers != nil {
		ctx = h.workers.Context
	} else {
		ctx = h.ctx
	}

	request := nrpc.NewRequest(ctx, h.nc, msg.Subject, msg.Reply)
	// extract method name & encoding from subject
	_, _, name, tail, err := nrpc.ParseSubject(
		"pb", 0, "PingService", 0, msg.Subject)
	if err != nil {
		log.Printf("PingServiceHanlder: PingService subject parsing failed: %v", err)
		return
	}
	var span opentracing.Span
	if len(tail) > 0 {
		addon := tail[0] // tail: [traceid[,encoding]] TODO: support for [encoding]
		ctx, span = newServerSpanFromString(ctx, addon, "/proto.TestService/"+name)
		//defer span.Finish()
		span.LogKV("event", "printlnStarted")
		tail = tail[1:]
	}
	defer finishClientSpan(span, err)
	//span.LogKV("event", "printlnStarted2")

	request.MethodName = name
	request.SubjectTail = tail

	// call handler and form response
	var immediateError *nrpc.Error
	switch name {
	case "Ping":
		_, request.Encoding, err = nrpc.ParseSubjectTail(0, request.SubjectTail)
		if err != nil {
			log.Printf("PingHanlder: Ping subject parsing failed: %v", err)
			break
		}
		var req pb.PingRequest
		if err = nrpc.Unmarshal(request.Encoding, msg.Data, &req); err != nil {
			log.Printf("PingHandler: Ping request unmarshal failed: %v", err)
			immediateError = &nrpc.Error{
				Type:    nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
			serverRequestsForPingService.WithLabelValues(
				"Ping", request.Encoding, "unmarshal_fail").Inc()
		} else {
			request.Handler = func(ctx context.Context) (proto.Message, error) {
				innerResp, err := h.server.Ping(ctx, &req)
				if err != nil {
					return nil, err
				}
				return innerResp, err
			}
		}
	case "PingEmpty":
		_, request.Encoding, err = nrpc.ParseSubjectTail(0, request.SubjectTail)
		if err != nil {
			log.Printf("PingEmptyHanlder: PingEmpty subject parsing failed: %v", err)
			break
		}
		var req pb.Empty
		if err = nrpc.Unmarshal(request.Encoding, msg.Data, &req); err != nil {
			log.Printf("PingEmptyHandler: PingEmpty request unmarshal failed: %v", err)
			immediateError = &nrpc.Error{
				Type:    nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
			serverRequestsForPingService.WithLabelValues(
				"PingEmpty", request.Encoding, "unmarshal_fail").Inc()
		} else {
			request.Handler = func(ctx context.Context) (proto.Message, error) {
				innerResp, err := h.server.PingEmpty(ctx, &req)
				if err != nil {
					return nil, err
				}
				return innerResp, err
			}
		}
	case "PingError":
		_, request.Encoding, err = nrpc.ParseSubjectTail(0, request.SubjectTail)
		if err != nil {
			log.Printf("PingErrorHanlder: PingError subject parsing failed: %v", err)
			break
		}
		var req pb.PingRequest
		if err = nrpc.Unmarshal(request.Encoding, msg.Data, &req); err != nil {
			log.Printf("PingErrorHandler: PingError request unmarshal failed: %v", err)
			immediateError = &nrpc.Error{
				Type:    nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
			serverRequestsForPingService.WithLabelValues(
				"PingError", request.Encoding, "unmarshal_fail").Inc()
		} else {
			request.Handler = func(ctx context.Context) (proto.Message, error) {
				innerResp, err := h.server.PingError(ctx, &req)
				if err != nil {
					return nil, err
				}
				return innerResp, err
			}
		}
	case "PingList":
		_, request.Encoding, err = nrpc.ParseSubjectTail(0, request.SubjectTail)
		if err != nil {
			log.Printf("PingListHanlder: PingList subject parsing failed: %v", err)
			break
		}
		var req pb.PingRequest
		if err = nrpc.Unmarshal(request.Encoding, msg.Data, &req); err != nil {
			log.Printf("PingListHandler: PingList request unmarshal failed: %v", err)
			immediateError = &nrpc.Error{
				Type:    nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
			serverRequestsForPingService.WithLabelValues(
				"PingList", request.Encoding, "unmarshal_fail").Inc()
		} else {
			request.EnableStreamedReply()
			es := pingServiceStreamingPingListServer{}
			es.request = request
			request.Handler = func(ctx context.Context) (proto.Message, error) {
				err = h.server.PingList(&req, es)
				return nil, err
			}
		}
		/*
			case "PingStream":
				_, request.Encoding, err = nrpc.ParseSubjectTail(0, request.SubjectTail)
				if err != nil {
					log.Printf("PingStreamHanlder: PingStream subject parsing failed: %v", err)
					break
				}
				var req pb.PingRequest
				if err = nrpc.Unmarshal(request.Encoding, msg.Data, &req); err != nil {
					log.Printf("PingStreamHandler: PingStream request unmarshal failed: %v", err)
					immediateError = &nrpc.Error{
						Type:    nrpc.Error_CLIENT,
						Message: "bad request received: " + err.Error(),
					}
					serverRequestsForPingService.WithLabelValues(
						"PingStream", request.Encoding, "unmarshal_fail").Inc()
				} else {
					request.EnableStreamedReply()
					es := pingServiceStreamingPingStreamServer{}
					es.request = request
					request.Handler = func(ctx context.Context) (proto.Message, error) {
						err = h.server.PingStream(&req, es)
						return nil, err
					}
				}
		*/
	case "TimeService":
		_, request.Encoding, err = nrpc.ParseSubjectTail(0, request.SubjectTail)
		if err != nil {
			log.Printf("TimeServiceHanlder: TimeService subject parsing failed: %v", err)
			break
		}
		var req ticker.TimeRequest
		if err = nrpc.Unmarshal(request.Encoding, msg.Data, &req); err != nil {
			log.Printf("TimeServiceHandler: TimeService request unmarshal failed: %v", err)
			immediateError = &nrpc.Error{
				Type:    nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
			serverRequestsForPingService.WithLabelValues(
				"TimeService", request.Encoding, "unmarshal_fail").Inc()
		} else {
			request.EnableStreamedReply()
			es := pingServiceStreamingTimeServiceServer{}
			es.request = request
			request.Handler = func(ctx context.Context) (proto.Message, error) {
				err = h.server.TimeService(&req, es)
				return nil, err
			}
		}
	default:
		log.Printf("PingServiceHandler: unknown name %q", name)
		immediateError = &nrpc.Error{
			Type:    nrpc.Error_CLIENT,
			Message: "unknown name: " + name,
		}
		serverRequestsForPingService.WithLabelValues(
			"PingService", request.Encoding, "name_fail").Inc()
	}
	request.AfterReply = func(request *nrpc.Request, success, replySuccess bool) {
		if !replySuccess {
			serverRequestsForPingService.WithLabelValues(
				request.MethodName, request.Encoding, "sendreply_fail").Inc()
		}
		if success {
			serverRequestsForPingService.WithLabelValues(
				request.MethodName, request.Encoding, "success").Inc()
		} else {
			serverRequestsForPingService.WithLabelValues(
				request.MethodName, request.Encoding, "handler_fail").Inc()
		}
		// report metric to Prometheus
		serverHETForPingService.WithLabelValues(request.MethodName).Observe(
			request.Elapsed().Seconds())
	}
	if immediateError == nil {
		if h.workers != nil {
			// Try queuing the request
			if err := h.workers.QueueRequest(request); err != nil {
				log.Printf("nrpc: Error queuing the request: %s", err)
			}
		} else {
			// Run the handler synchronously
			request.RunAndReply()
		}
	}

	if immediateError != nil {
		if err := request.SendReply(nil, immediateError); err != nil {
			log.Printf("PingServiceHandler: PingService handler failed to publish the response: %s", err)
			serverRequestsForPingService.WithLabelValues(
				request.MethodName, request.Encoding, "handler_fail").Inc()
		}
		serverHETForPingService.WithLabelValues(request.MethodName).Observe(
			request.Elapsed().Seconds())
	} else {
	}
}

type pingServiceStreamingPingListServer struct {
	grpc.ServerStream
	request *nrpc.Request
}

func (x pingServiceStreamingPingListServer) Send(m *pb.PingResponse) error {
	x.request.SendStreamReply(m)
	return nil
}
func (x pingServiceStreamingPingListServer) Context() context.Context {
	return x.request.StreamContext
}

/*
type pingServiceStreamingPingStreamServer struct {
	grpc.ServerStream
	request *nrpc.Request
}

func (x pingServiceStreamingPingStreamServer) Send(m *pb.PingResponse) error {
	x.request.SendStreamReply(m)
	return nil
}
*/
type pingServiceStreamingTimeServiceServer struct {
	grpc.ServerStream
	request *nrpc.Request
}

func (x pingServiceStreamingTimeServiceServer) Send(m *ticker.TimeResponse) error {
	x.request.SendStreamReply(m)
	return nil
}
func (x pingServiceStreamingTimeServiceServer) Context() context.Context {
	return x.request.StreamContext
}

type PingServiceClient struct {
	nc               nrpc.NatsConn
	PkgSubject       string
	PkgParaminstance string
	Subject          string
	Encoding         string
	Timeout          time.Duration
}

func NewPingServiceClient(nc nrpc.NatsConn, pkgParaminstance string) *PingServiceClient {
	return &PingServiceClient{
		nc:               nc,
		PkgSubject:       "pb",
		PkgParaminstance: pkgParaminstance,
		Subject:          "PingService",
		Encoding:         "protobuf",
		Timeout:          5 * time.Second,
	}
}

func (c *PingServiceClient) Ping(ctx context.Context, req *pb.PingRequest) (resp *pb.PingResponse, err error) {
	start := time.Now()

	subject := c.PkgSubject + "." + c.Subject + "." + "Ping"

	// call
	resp = &pb.PingResponse{}
	err = nrpc.Call(req, resp, c.nc, subject, c.Encoding, c.Timeout)
	if err != nil {
		clientCallsForPingService.WithLabelValues(
			"Ping", c.Encoding, "call_fail").Inc()
		return // already logged
	}

	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForPingService.WithLabelValues("Ping").Observe(elapsed)
	clientCallsForPingService.WithLabelValues(
		"Ping", c.Encoding, "success").Inc()

	return
}

func (c *PingServiceClient) PingEmpty(ctx context.Context, req *pb.Empty) (resp *pb.PingResponse, err error) {
	start := time.Now()

	subject := c.PkgSubject + "." + c.Subject + "." + "PingEmpty"

	// call
	resp = &pb.PingResponse{}
	err = nrpc.Call(req, resp, c.nc, subject, c.Encoding, c.Timeout)
	if err != nil {
		clientCallsForPingService.WithLabelValues(
			"PingEmpty", c.Encoding, "call_fail").Inc()
		return // already logged
	}

	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForPingService.WithLabelValues("PingEmpty").Observe(elapsed)
	clientCallsForPingService.WithLabelValues(
		"PingEmpty", c.Encoding, "success").Inc()

	return
}

func (c *PingServiceClient) PingError(ctx context.Context, req *pb.PingRequest) (resp *pb.Empty, err error) {
	start := time.Now()

	subject := c.PkgSubject + "." + c.Subject + "." + "PingError"

	// call
	resp = &pb.Empty{}
	err = nrpc.Call(req, resp, c.nc, subject, c.Encoding, c.Timeout)
	if err != nil {
		clientCallsForPingService.WithLabelValues(
			"PingError", c.Encoding, "call_fail").Inc()
		return // already logged
	}

	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForPingService.WithLabelValues("PingError").Observe(elapsed)
	clientCallsForPingService.WithLabelValues(
		"PingError", c.Encoding, "success").Inc()

	return
}

func (c *PingServiceClient) PingList(
	req *pb.PingRequest,
	stream pb.PingService_PingListServer,
) error {
	start := time.Now()
	subject := c.PkgSubject + "." + c.Subject + "." + "PingList"
	ctx := stream.Context()
	sub, err := nrpc.StreamCall(ctx, c.nc, subject, req, c.Encoding, c.Timeout)
	if err != nil {
		clientCallsForPingService.WithLabelValues(
			"PingList", c.Encoding, "error").Inc()
		return err
	}

	var res pb.PingResponse
	for {
		err = sub.Next(&res)
		if err != nil {
			break
		}
		err = stream.Send(&res)
		if err != nil {
			log.Printf("> call send error %v\n", err)
			break
		}
	}
	if err == nrpc.ErrEOS {
		err = nil
	}
	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForPingService.WithLabelValues("PingList").Observe(elapsed)
	clientCallsForPingService.WithLabelValues(
		"PingList", c.Encoding, "success").Inc()
	return err
}

/*
func (c *PingServiceClient) PingStream(
	req *pb.PingRequest,
	stream pb.PingService_PingStreamServer,
) error {
	start := time.Now()
	subject := c.PkgSubject + "." + c.Subject + "." + "PingStream"
	ctx := stream.Context()
	sub, err := nrpc.StreamCall(ctx, c.nc, subject, req, c.Encoding, c.Timeout)
	if err != nil {
		clientCallsForPingService.WithLabelValues(
			"PingStream", c.Encoding, "error").Inc()
		return err
	}

	var res pb.PingResponse
	for {
		err = sub.Next(&res)
		if err != nil {
			break
		}
		err = stream.Send(&res)
		if err != nil {
			log.Printf("> call send error %v\n", err)
			break
		}
	}
	if err == nrpc.ErrEOS {
		err = nil
	}
	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForPingService.WithLabelValues("PingStream").Observe(elapsed)
	clientCallsForPingService.WithLabelValues(
		"PingStream", c.Encoding, "success").Inc()
	return err
}
*/
func (c *PingServiceClient) TimeService(
	req *ticker.TimeRequest,
	stream pb.PingService_TimeServiceServer,
) error {
	start := time.Now()
	subject := c.PkgSubject + "." + c.Subject + "." + "TimeService"
	ctx := stream.Context()
	sub, err := nrpc.StreamCall(ctx, c.nc, subject, req, c.Encoding, c.Timeout)
	if err != nil {
		clientCallsForPingService.WithLabelValues(
			"TimeService", c.Encoding, "error").Inc()
		return err
	}

	var res ticker.TimeResponse
	for {
		err = sub.Next(&res)
		if err != nil {
			break
		}
		err = stream.Send(&res)
		if err != nil {
			log.Printf("> call send error %v\n", err)
			break
		}
	}
	if err == nrpc.ErrEOS {
		err = nil
	}
	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForPingService.WithLabelValues("TimeService").Observe(elapsed)
	clientCallsForPingService.WithLabelValues(
		"TimeService", c.Encoding, "success").Inc()
	return err
}

type Client struct {
	nc              nrpc.NatsConn
	defaultEncoding string
	defaultTimeout  time.Duration
	pkgSubject      string
	PingService     *PingServiceClient
}

/*
func NewClient(nc nrpc.NatsConn) *Client {
	c := Client{
		nc:              nc,
		defaultEncoding: "protobuf",
		defaultTimeout:  5 * time.Second,
		pkgSubject:      "pb",
	}
	c.PingService = NewPingServiceClient(nc)
	return &c
}

func (c *Client) SetEncoding(encoding string) {
	c.defaultEncoding = encoding
	if c.PingService != nil {
		c.PingService.Encoding = encoding
	}
}

func (c *Client) SetTimeout(t time.Duration) {
	c.defaultTimeout = t
	if c.PingService != nil {
		c.PingService.Timeout = t
	}
}
*/
func init() {
	// register metrics for service PingService
	prometheus.MustRegister(clientRCTForPingService)
	prometheus.MustRegister(serverHETForPingService)
	prometheus.MustRegister(clientCallsForPingService)
	prometheus.MustRegister(serverRequestsForPingService)
}
