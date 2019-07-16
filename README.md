Overview
    Testing project for combine js and gRPC service.

Structure

grpc-login-ec2
    api
        proto
            auth.proto              => Protobuf 功能定義
    auth-vue                        => Vue 登入介面
    envoy                           => Envoy proxy
    go-client                       => debug 測試用 client
    server                          => grpc bankend service
        internal
            service
                v1
                    auth-service.go => 實作 auth server 功能