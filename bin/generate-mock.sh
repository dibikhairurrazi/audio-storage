#!/usr/bin/env bash

set -o pipefail

for file in `find . -name '*.go' | grep -v cmd | grep module/audio-storage | grep -v module/audio-storage/test`; do
    if `grep -q 'interface {' ${file}`; then
        dest=${file//.\/module\/payment\//}
        echo $dest
        mockgen -source ${file} -destination module/audio-storage/test/mock/${dest}
    fi
done