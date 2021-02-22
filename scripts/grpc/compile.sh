#!/bin/bash
# This script will compile the grpc protocol

pushd ../../service/views/grpc

protoc -I calculator/ calculator/calculator.proto --go_out=plugins=grpc:calculator
