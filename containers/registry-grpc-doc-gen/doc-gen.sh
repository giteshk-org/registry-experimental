#!/bin/bash
set -euo pipefail
#
#mkdir -p /workspace/protos
#mkdir -p /workspace/out
#
args=("$@")
#if [ "${#args[@]}" -lt 2 ]; then args+=(protos/**/*.proto); fi
#
#exec protoc -I/usr/include -Iprotos --doc_out=/workspace/out "${args[@]}"

mkdir -p /workspace/args[0]
mkdir -p /workspace/args[0]/protos

MIMETYPE=$(registry get args[0] | jq -r .mimeType)
FILENAME=$(registry get args[0] | jq -r .filename)
registry get arg[0] --contents > /workspace/args[0]/$FILENAME

if [$MIMETYPE eq "application/x.protobuf+gzip"]
then
  tar xvfz /workspace/args[0]/$FILENAME /workspace/args[0]/protos
fi

protoc /workspace/args[0]/protos/**/*.proto -Iworkspace/args[0]/protos --proto_path="/googleapis-common-protos" --doc_out=/workspace/args[0] --doc_opt=html,index.html