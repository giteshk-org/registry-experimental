#FROM golang:1.18
#
#RUN git clone https://github.com/googleapis/api-common-protos /googleapis-common-protos
#RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.1
#
#VOLUME /workspace
#WORKDIR /workspace
#CMD ["/go/bin/protoc-gen-doc **/*.proto --proto_path='.' --proto_path=/googleapis-common-protos --doc_opt='html,api.html' --doc_out='.'"]

FROM pseudomuto/protoc-gen-doc:v1.5.1 as protoc-gen-doc

FROM ghcr.io/apigee/registry-tools:latest

WORKDIR /

COPY --from=protoc-gen-doc /go/bin/protoc-gen-doc /bin/

ARG PROTOC_VERSION=21.8
RUN wget -q -O protoc.zip "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" && \
    unzip -o protoc.zip bin/protoc -d /

VOLUME /workspace

RUN git clone https://github.com/googleapis/api-common-protos /googleapis-common-protos
