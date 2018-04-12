#!/bin/sh
go build && ./codegen -cout ../rule_gen.go -tout ../rule_gen_test.go && \
    gofmt -w=true ../rule_gen.go && \
    gofmt -w=true ../rule_gen_test.go && \
    rm codegen
