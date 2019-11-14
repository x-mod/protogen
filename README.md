# protogen
Generator for ProtoBuffer files using templates

````
go get github.com/golang/protobuf/protoc-gen-go
```

test:

````
go run main.go -i examples -o examples -p code.subscriber.one/subscriber/protogen/examples http  server
````