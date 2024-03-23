#!/usr/bin/env bash

# How to run manually:
# docker build --pull --rm -f "contrib/devtools/Dockerfile" -t cosmossdk-proto:latest "contrib/devtools"
# docker run --rm -v $(pwd):/workspace --workdir /workspace cosmossdk-proto sh ./scripts/protocgen.sh

set -e

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./tokenfactory -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    # this regex checks if a proto file has its go_package set to cosmossdk.io/api/...
    # gogo proto files SHOULD ONLY be generated if this is false
    # we don't want gogo proto to run for proto files which are natively built for google.golang.org/protobuf
   if grep go_package $file &>/dev/null; then
      echo "Current directory: $(pwd)"
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..


# move proto files to the right places
cp -r github.com/osmosis-labs/tokenfactory/* ./
rm -rf github.com

go mod tidy


# echo "Formatting protobuf files"
# find ./ -name "*.proto" -exec clang-format -i {} \;