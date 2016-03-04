#!/bin/sh

mkdir -p gen

for i in assets/*.js ; do
    uglifyjs $i --compress --mangle -o gen/`basename $i`
done

cp assets/*.txt gen/
