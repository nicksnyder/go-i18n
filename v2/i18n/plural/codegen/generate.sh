#!/bin/sh
go build && ./codegen -cout ../pluralrules_gen.go -tout ../pluralrules_gen_test.go && \
    gofmt -w=true ../pluralrules_gen.go && \
    gofmt -w=true ../pluralrules_gen_test.go && \
    rm codegen
