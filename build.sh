#!/bin/bash
# windows
#gox -osarch="darwin/amd64" -output bin/feibor_darwin_amd64 ./main

work_path=$(cd `dirname $0`; pwd)

echo $work_path

cd $work_path

mkdir -p bin

# 编译Linux 和windows版本
#gox -osarch="darwin/amd64" -output ${work_path}/bin/fabric_darwin_amd64 ./runner
#gox -osarch="windows/amd64" -output ${work_path}/bin/fabric_win_amd64 ./runner
gox -osarch="linux/amd64" -output ${work_path}/bin/fabric_linux_amd64 ./runner

echo "done!"