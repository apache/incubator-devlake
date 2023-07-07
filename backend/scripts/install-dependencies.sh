#!/bin/sh
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

validate_prerequisites() {
  go version > /dev/null >&1 &&\
  python --version > /dev/null >&1 &&\
  poetry --version > /dev/null >&1
  return $?
}

matches() {
  return $(echo "$1" | grep -q "$2" > /dev/null 2>&1)
}

while [ $# -gt 0 ]; do
  case "$1" in
    --skip-check=*)
      skip_check="${1#*=}"
      ;;
    *)
      printf "***************************\n"
      printf "* Error: Invalid argument. Allowed args:\n"
      printf " --skip-check=true/false \n"
      exit 1
  esac
  shift
done

if [ "$skip_check" != "true" ]; then
  if ! validate_prerequisites; then
    printf "Failed to validate prerequisites\n"
    exit 1
  fi
fi

if matches $(uname -s) "Linux" ; then
  #TODO detect distro and call the correct package-manager
  echo "This is Linux"
elif matches $(uname -s) "Darwin" ; then
  #TODO make brew calls
  echo "This is Mac"
elif matches $(uname -s) "MINGW" ; then
  #TODO "probably error out here..."
  echo "This is Windows"
fi