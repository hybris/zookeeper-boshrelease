#!/bin/bash
set -e
set -x

vagrant up
vagrant ssh -c "sudo bash -c 'mkdir -p /zk-backyp/tmp && cd /zk-backyp && GOPATH=~/.go ./build.sh'"
vagrant ssh -c "sudo bash -c 'cd /zk-backyp/tmp && /zk-backyp/test/zk-backyp.sh'"