// versions:

package whatsapp_client_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const _ = grpc.SupportPackageIsVersion7

//
type WhatsAppClientServiceApiClient interface {
	Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServiceStatus, error)
	SendWhatsApp(ctx context.Context, in *SendWhatsAppReq, opts ...grpc.CallOption) (*ServiceStatus, error)
}

type whatsAppClientServiceApiClient struct {
	cc grpc.ClientConnInterface
}

func NewWhatsAppClientServiceApiClient(cc grpc.ClientConnInterface) WhatsAppClientServiceApiClient {
	return &whatsAppClientServiceApiClient{cc}
}

func (c *whatsAppClientServiceApiClient) Status(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ServiceStatus, error) {
	out := new(ServiceStatus)
	err := c.cc.Invoke(ctx, "/influenzanet.whatsapp_client_service.WhatsAppClientServiceApi/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *whatsAppClientServiceApiClient) SendWhatsApp(ctx context.Context, in *SendWhatsAppReq, opts ...grpc.CallOption) (*ServiceStatus, error) {
	out := new(ServiceStatus)
	err := c.cc.Invoke(ctx, "/influenzanet.whatsapp_client_service.WhatsAppClientServiceApi/SendWhatsApp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type WhatsAppClientServiceApiServer interface {
	Status(context.Context, *emptypb.Empty) (*ServiceStatus, error)
	SendWhatsApp(context.Context, *SendWhatsAppReq) (*ServiceStatus, error)
	mustEmbedUnimplementedWhatsAppClientServiceApiServer()
}

type UnimplementedWhatsAppClientServiceApiServer struct {
}

func (UnimplementedWhatsAppClientServiceApiServer) Status(context.Context, *emptypb.Empty) (*ServiceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}

func (UnimplementedWhatsAppClientServiceApiServer) SendWhatsApp(context.Context, *SendWhatsAppReq) (*ServiceStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendWhatsApp not implemented")
}
func (UnimplementedWhatsAppClientServiceApiServer) mustEmbedUnimplementedWhatsAppClientServiceApiServer() {
}

type UnsafeWhatsAppClientServiceApiServer interface {
	mustEmbedUnimplementedWhatsAppClientServiceApiServer()
}

func RegisterWhatsAppClientServiceApiServer(s grpc.ServiceRegistrar, srv WhatsAppClientServiceApiServer) {
	s.RegisterService(&WhatsAppClientServiceApi_ServiceDesc, srv)
}

func _WhatsAppClientServiceApi_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WhatsAppClientServiceApiServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/influenzanet.whatsapp_client_service.WhatsAppClientServiceApi/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WhatsAppClientServiceApiServer).Status(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WhatsAppClientServiceApi_SendWhatsApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendWhatsAppReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WhatsAppClientServiceApiServer).SendWhatsApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/influenzanet.whatsapp_client_service.WhatsAppClientServiceApi/SendWhatsApp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WhatsAppClientServiceApiServer).SendWhatsApp(ctx, req.(*SendWhatsAppReq))
	}
	return interceptor(ctx, in, info, handler)
}

var WhatsAppClientServiceApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "influenzanet.whatsapp_client_service.WhatsAppClientServiceApi",
	HandlerType: (*WhatsAppClientServiceApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _WhatsAppClientServiceApi_Status_Handler,
		},
		{
			MethodName: "SendWhatsApp",
			Handler:    _WhatsAppClientServiceApi_SendWhatsApp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "whatsapp_client_service/whatsapp-client-service.proto",
}
