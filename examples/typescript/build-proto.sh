#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building Proto files for TypeScript...${NC}"

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo -e "${RED}node_modules not found. Running npm install...${NC}"
    npm install
fi

# Clean old generated files
echo "Cleaning old generated files..."
find ./proto -name '*.js' -o -name '*.d.ts' -o -name '*_pb.ts' | xargs rm -f 2>/dev/null

# Generate JavaScript files with CommonJS
echo "Generating JavaScript files..."
./node_modules/.bin/grpc_tools_node_protoc \
  --js_out=import_style=commonjs,binary:./proto \
  --grpc_out=grpc_js:./proto \
  --plugin=protoc-gen-grpc=./node_modules/.bin/grpc_tools_node_protoc_plugin \
  -I ./proto \
  ./proto/events.proto \
  ./proto/publisher.proto

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ JavaScript files generated successfully${NC}"
else
    echo -e "${RED}✗ Failed to generate JavaScript files${NC}"
    exit 1
fi

# Generate TypeScript definitions (optional)
echo "Generating TypeScript definitions..."
if [ -f "./node_modules/.bin/protoc-gen-ts" ]; then
    npx protoc \
      --plugin=./node_modules/.bin/protoc-gen-ts \
      --ts_out=service=grpc-node,mode=grpc-js:./proto \
      -I ./proto \
      ./proto/events.proto \
      ./proto/publisher.proto 2>/dev/null

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ TypeScript definitions generated${NC}"
    else
        echo -e "${YELLOW}⚠ TypeScript definitions generation skipped${NC}"
    fi
else
    echo -e "${YELLOW}⚠ ts-protoc-gen not found, skipping TypeScript definitions${NC}"
fi

echo -e "${GREEN}✓ Proto build complete!${NC}"

# List generated files
echo ""
echo "Generated files:"
find ./proto -type f \( -name "*.js" -o -name "*.d.ts" -o -name "*_pb.ts" \) | sort

# Fix imports in generated files (if needed)
echo ""
echo "Fixing imports in generated files..."

# Fix the import paths in grpc files
if [ -f "./proto/publisher_grpc_pb.js" ]; then
    # Fix relative imports for macOS and Linux
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        sed -i '' "s|require('publisher_pb.js')|require('./publisher_pb.js')|g" ./proto/publisher_grpc_pb.js
        sed -i '' "s|require('events_pb.js')|require('./events_pb.js')|g" ./proto/publisher_grpc_pb.js
    else
        # Linux
        sed -i "s|require('publisher_pb.js')|require('./publisher_pb.js')|g" ./proto/publisher_grpc_pb.js
        sed -i "s|require('events_pb.js')|require('./events_pb.js')|g" ./proto/publisher_grpc_pb.js
    fi
    echo -e "${GREEN}✓ Fixed imports${NC}"
fi