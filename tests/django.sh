#!/usr/bin/env bash

# django
DIR=upsun_django
rm -rf $DIR
export UPSUN_STACK=django
export UPSUN_USEDEFAULTS=1
export UPSUN_SHOWCOMMENTS=0
go build ./cmd/upsunify/
cookiecutter https://github.com/cookiecutter/cookiecutter-django --default-config --no-input -o $DIR
cd $DIR
../upsunify