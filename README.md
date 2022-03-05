# micro-client


## Generate
```sh
$ protoc -I communicate/ communicate/product.proto --go_out=plugins=grpc:communicate
$ protoc communicate/product.proto --go_out=plugins=grpc:.
```
