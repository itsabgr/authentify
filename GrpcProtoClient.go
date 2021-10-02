package authentify

import (
	"context"
	"google.golang.org/grpc"
)

type GrpcProtoClient interface {
	SendOTP(ctx context.Context, in *ProtoSendOTPRequest, opts ...grpc.CallOption) (*ProtoSendOTPResponse, error)
	CheckOTP(ctx context.Context, in *ProtoCheckOTPRequest, opts ...grpc.CallOption) (*ProtoCheckOTPResponse, error)
}

type grpcProtoClient struct {
	cc *grpc.ClientConn
}

func NewGrpcProtoClient(cc *grpc.ClientConn) GrpcProtoClient {
	return &grpcProtoClient{cc}
}

func (c *grpcProtoClient) SendOTP(ctx context.Context, in *ProtoSendOTPRequest, opts ...grpc.CallOption) (*ProtoSendOTPResponse, error) {
	out := new(ProtoSendOTPResponse)
	err := c.cc.Invoke(ctx, "/Authentify.Authenticator/SendOTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcProtoClient) CheckOTP(ctx context.Context, in *ProtoCheckOTPRequest, opts ...grpc.CallOption) (*ProtoCheckOTPResponse, error) {
	out := new(ProtoCheckOTPResponse)
	err := c.cc.Invoke(ctx, "/Authentify.Authenticator/CheckOTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
