#!/bin/sh -ex

rm -f Gopkg.lock Gopkg.toml
rm -rf vendor
dep init -v
slimdep -r -v -a github.com/frankbraun/slimdep
