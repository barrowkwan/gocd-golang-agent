#!/bin/bash
#
# Copyright 2016 ThoughtWorks, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
set -e

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"
PROJECT_DIR="$SCRIPT_DIR/../../"
BINARY_FILE="$SCRIPT_DIR/package/usr/bin/gocd-golang-agent"
DEB_FILE="gocd-golang-agent_0.1-1_all.deb"

echo "############################"
echo "Cross compiling for linux..."
echo "############################"

cd $PROJECT_DIR
CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .
mkdir -p `dirname $BINARY_FILE`
cp gocd-golang-agent $BINARY_FILE
chmod 0755 $BINARY_FILE

echo "############################"
echo "Packaging to deb file..."
echo "############################"

cd $SCRIPT_DIR
docker build -t gocd/deb-maker .
docker run -v ${PWD}:/build gocd/deb-maker /bin/bash -c "cd /build/package && fakeroot dpkg-deb --build . ../$DEB_FILE"

echo "############################"
echo "Making sure installer can be installed..."
echo "############################"

docker run -v ${PWD}:/build gocd/deb-maker /bin/bash -c "cd /build && dpkg -i $DEB_FILE && service gocd-golang-agent start && sleep 3 && ps aux | grep -v 'grep' | grep gocd-golang-agent"

echo "All check passed."
