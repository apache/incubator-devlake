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
cd "${0%/*}" # make sure we're in the correct dir

for test_dir in $(find . -type f -name "*_test.py" | xargs dirname | sort -u); do
  printf "Running Python tests in $test_dir\n"
  cd $test_dir
  poetry run pytest
  if [ $? != 0 ]; then
    exit 1
  fi
  cd -
done