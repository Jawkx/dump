#!/bin/bash

TARGETS=(
  "linux/amd64"
  "darwin/amd64"
  "windows/amd64"
)

mkdir -p "${OUTPUT_DIR}"

# Function to build the Go binary
build_binary() {
  local GOOS="${1}"
  local GOARCH="${2}"
  local OUTPUT_NAME="${3}"

  echo "Building for ${GOOS}/${GOARCH}..."
  
  # Set environment variables
  export GOOS="${GOOS}"
  export GOARCH="${GOARCH}"

  # Build the binary
  go build -v -ldflags "-X main.version=${VERSION} -s -w" -o "${OUTPUT_DIR}/${OUTPUT_NAME}"
  
  echo "Build finished: ${OUTPUT_DIR}/${OUTPUT_NAME}"
}

# Build for each target platform
for target in "${TARGETS[@]}"; do
  os=$(echo "$target" | cut -d'/' -f1)
  arch=$(echo "$target" | cut -d'/' -f2)
  output_name="cmtbot-${os}-${arch}"
  
  if [ "$os" == "windows" ]; then
    output_name="${output_name}.exe"
  fi
  
  build_binary "${os}" "${arch}" "${output_name}"
done

echo "All builds completed."
