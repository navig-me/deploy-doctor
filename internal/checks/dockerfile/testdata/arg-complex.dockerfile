# syntax=docker/dockerfile:1.7
ARG BASE_IMAGE=ubuntu:22.04
FROM ${BASE_IMAGE}
ARG APP_SECRET=dev-secret
RUN apt-get update && apt-get install -y curl \
  && rm -rf /var/lib/apt/lists/*
ENV APP_MODE=prod
COPY . .
CMD ["bash", "-lc", "echo ready"]
