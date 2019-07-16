protoc --proto_path=api/proto --js_out=import_style=commonjs,binary:auth-vue/node_modules/ --grpc-web_out=import_style=commonjs,mode=grpcwebtext:auth-vue/node_modules/ auth.proto
