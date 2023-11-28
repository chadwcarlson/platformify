#!/usr/bin/env bash

# laravel
DIR=laravel-example
rm -rf $DIR
export UPSUN_STACK=laravel
export UPSUN_USEDEFAULTS=1
export UPSUN_SHOWCOMMENTS=0
go build ./cmd/upsunify/
composer create-project laravel/laravel $DIR
cd $DIR
../upsunify