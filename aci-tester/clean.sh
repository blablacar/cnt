#!/bin/bash
set -x
set -e
dir=$( dirname $0 )
target=${dir}/../dist/aci-tester
rm -Rf ${target}
rm -f ${dir}/../dist/bindata/aci-tester.aci
