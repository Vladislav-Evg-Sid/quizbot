#!/bin/bash

cd "$(dirname "$0")/.." || exit

# Критически важно: пути в Windows формате для protoc
export PATH="$PATH:/c/ProgramData/chocolatey/bin"

# Пути к Go плагинам (Windows формат!)
GO_BIN_PATH="$(go env GOPATH)/bin"
PROTOC_GEN_GO="${GO_BIN_PATH}/protoc-gen-go.exe"
PROTOC_GEN_GO_GRPC="${GO_BIN_PATH}/protoc-gen-go-grpc.exe"
PROTOC_GEN_GW="${GO_BIN_PATH}/protoc-gen-grpc-gateway.exe"
PROTOC_GEN_OPENAPI="${GO_BIN_PATH}/protoc-gen-openapiv2.exe"

echo "Checking plugins..."
ls -la "$GO_BIN_PATH"/protoc-gen-*.exe 2>/dev/null || echo "WARNING: Plugins not found in Go bin"

echo ""
echo "Generating Protobuf files..."

# 1. Генерация gRPC кода
protoc -I ./api \
  -I ./api/google/api \
  --plugin=protoc-gen-go="${PROTOC_GEN_GO}" \
  --go_out=./internal/pb --go_opt=paths=source_relative \
  --plugin=protoc-gen-go-grpc="${PROTOC_GEN_GO_GRPC}" \
  --go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
  ./api/players_api/players.proto ./api/models/player_model.proto

# 2. Генерация gRPC-Gateway (если используете)
protoc -I ./api \
  -I ./api/google/api \
  --plugin=protoc-gen-grpc-gateway="${PROTOC_GEN_GW}" \
  --grpc-gateway_out=./internal/pb \
  --grpc-gateway_opt=paths=source_relative \
  --grpc-gateway_opt=logtostderr=true \
  ./api/players_api/players.proto

# 3. Генерация OpenAPI (если используете)
protoc -I ./api \
  -I ./api/google/api \
  --plugin=protoc-gen-openapiv2="${PROTOC_GEN_OPENAPI}" \
  --openapiv2_out=./internal/pb/swagger \
  --openapiv2_opt=logtostderr=true \
  ./api/players_api/players.proto

echo ""
echo "Generation completed!"
echo "Files in: internal/pb/"
ls -la internal/pb/ 2>/dev/null || echo "No output directory"
read -p "Press Enter to continue..."