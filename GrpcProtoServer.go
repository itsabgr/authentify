package authentify

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcProtoServer interface {
	SendOTP(context.Context, *ProtoSendOTPRequest) (*ProtoSendOTPResponse, error)
	CheckOTP(context.Context, *ProtoCheckOTPRequest) (*ProtoCheckOTPResponse, error)
	mustEmbedUnimplementedAuthenticatorServer()
}

type grpcProtoServer struct {
	authenticator Authenticator
	senders       SendersMap
}

func (s *grpcProtoServer) SendOTP(ctx context.Context, req *ProtoSendOTPRequest) (*ProtoSendOTPResponse, error) {
	senderName := req.GetType()
	sender, ok := s.senders[senderName]
	if false == ok {
		return nil, errors.New("unsupported type")
	}
	token, prefix, err := s.authenticator.SendCode(ctx, sender, req.GetTo())
	if err != nil {
		return nil, err
	}
	res := new(ProtoSendOTPResponse)
	res.Salt = token.Salt().Bytes()
	res.Deadline = timestamppb.New(token.Deadline())
	res.Prefix = prefix
	return res, nil
}
func (s *grpcProtoServer) CheckOTP(ctx context.Context, req *ProtoCheckOTPRequest) (*ProtoCheckOTPResponse, error) {
	token, err := s.authenticator.RetrieveToken(ctx, req.GetReceiver())
	if err != nil {
		return nil, err
	}
	res := new(ProtoCheckOTPResponse)
	res.Deadline = timestamppb.New(token.Deadline())
	res.Valid = false
	code, err := NewCode(req.GetPrefix(), req.GetCode())
	if err != nil {
		return res, nil
	}
	salt, err := AsSalt(req.GetSalt())
	if err != nil {
		return res, nil
	}
	res.Valid = token.Validate(req.GetType(), code, salt)
	return res, nil
}
func (s *grpcProtoServer) mustEmbedUnimplementedAuthenticatorServer() {}

func NewGrpcProtoServer(a Authenticator, senders SendersMap) (GrpcProtoServer, error) {
	s := &grpcProtoServer{}
	s.authenticator = a
	s.senders = senders
	return s, nil
}

func GrpcProtoServiceDesc() *grpc.ServiceDesc {
	clone := grpcProtoServiceDesc
	return &clone
}

var grpcProtoServiceDesc = grpc.ServiceDesc{
	ServiceName: "Authentify.Authenticator",
	HandlerType: (*GrpcProtoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendOTP",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(ProtoSendOTPRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.(GrpcProtoServer).SendOTP(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/Authentify.Authenticator/SendOTP",
				}
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(GrpcProtoServer).SendOTP(ctx, req.(*ProtoSendOTPRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
		{
			MethodName: "CheckOTP",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(ProtoCheckOTPRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.(GrpcProtoServer).CheckOTP(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/Authentify.Authenticator/CheckOTP",
				}
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(GrpcProtoServer).CheckOTP(ctx, req.(*ProtoCheckOTPRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/Authentify.proto",
}
