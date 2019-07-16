Overview
    Testing project for combine js and gRPC service.

Structure

grpc-login-ec2<br/>
&nbsp;&nbsp;&nbsp;&nbsp;api<br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;proto<br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;auth.proto              => Protobuf 功能定義<br/>
&nbsp;&nbsp;&nbsp;&nbsp;auth-vue                        => Vue 登入介面<br/>
&nbsp;&nbsp;&nbsp;&nbsp;envoy                           => Envoy proxy<br/>
&nbsp;&nbsp;&nbsp;&nbsp;go-client                       => debug 測試用 client<br/>
&nbsp;&nbsp;&nbsp;&nbsp;server                          => grpc bankend service<br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;internal<br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;service<br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;v1<br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;auth-service.go => 實作 auth server 功能<br/>
