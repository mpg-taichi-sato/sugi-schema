FROM golang:1.13.1  
  
RUN apt-get update && apt-get install -y unzip  
  
# Install protobuf  
# @see https://github.com/yoshi42662/go-grpc/blob/master/server/Dockerfile  
RUN mkdir -p /tmp/protoc && \  
  curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip > /tmp/protoc/protoc.zip && \  
  cd /tmp/protoc && \  
  unzip protoc.zip && \  
  mv include/* /usr/local/include/ && \
  cp /tmp/protoc/bin/protoc /usr/local/bin && \  
  chmod go+rx /usr/local/bin/protoc && \  
  cd /tmp && \  
  rm -r /tmp/protoc  
  
RUN go get -u google.golang.org/grpc && \  
  go get -u github.com/golang/protobuf/protoc-gen-go