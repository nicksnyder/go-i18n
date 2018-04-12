#!/bin/sh
go build && ./codegen -cout ../rules_gen.go -tout ../rules_gen_test.go && \
    gofmt -w=true ../rules_gen.go && \
    gofmt -w=true ../rules_gen_test.go && \
    rm codegen
