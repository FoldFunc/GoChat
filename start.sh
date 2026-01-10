#!/bin/bash

cd server

# Build Docker image
docker build -t foldfunc-server .

# Run container, using PORT env (443/80 not needed in dev or Railway)
docker run \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  -e PORT=8080 \
  foldfunc-server

