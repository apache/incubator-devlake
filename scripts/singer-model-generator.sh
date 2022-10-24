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

#===================================== constants =======================================

time_format='
  {
    "type": "time.Time",
    "imports": ["time"]
  }
'

#===================================== functions =======================================

json_path=$1 # e.g. "./config/singer/github.json"
plugin_path=$2 # e.g. "./plugins/github_singer"
tap_stream=$3 # e.g. "issues" or "--all" to generate all streams

print_json() {
  jq -r < "$1" # for debugging
}

handle_error() {
    exitcode=$1
    if [ "$exitcode" != 0 ]; then
      exit "$exitcode"
    fi
}

generate() {
  tap_stream=$1
  package="generated"
  file_name="$tap_stream".go
  output_path="$plugin_path/models/generated/$file_name"
  tmp_dir=$(mktemp -d -t schema-XXXXX)
  json_schema_path="$tmp_dir"/"$tap_stream"
  # add, as necessary, more elif blocks for additional transformations
  modified_schema=$(jq --argjson tf "$time_format" '
      .streams[] |
      select(.stream=="'"$tap_stream"'").schema |
        . += { "$schema": "http://json-schema.org/draft-07/schema#" } |
        walk(
          if type == "object" and .format == "date-time" then
            . += { "goJSONSchema": ($tf) }
          elif "place_holder" == "" then
            empty
          else . end
        )
  ' < "$json_path")
  handle_error $?
  # additional cleanup
  modified_schema=$(echo "$modified_schema" | sed -r "/\"null\",/d")
  modified_schema=$(echo "$modified_schema" | sed -r "/.*additionalProperties.*/d")
  echo "$modified_schema" > "$json_schema_path" &&\
  gojsonschema -v -p "$package" "$json_schema_path" -o "$output_path"
  handle_error $?
  echo "$output_path"
  # prepend the license text to the generated files
  cp "$output_path" "$output_path".bak
  license_header="$(printf "/*\n%s\n*/\n" "$(cat .golangci-goheader.template)")"
  echo "$license_header" > "$output_path"
  cat "$output_path".bak >> "$output_path"
  rm "$output_path".bak
}

#======================================================================================

if [ $# != 3 ]; then
  printf "not enough args. Usage: <json_path> <tap_stream> <output_path>: e.g.\n    \"./config/singer/github.json\" \"issues\" \"./plugins/github_singer\"\n"
  exit 1
fi

if [ "$tap_stream" = "--all" ]; then
  for stream in $(jq -r '.streams[].stream' < "$json_path"); do
    generate "$stream"
    handle_error $?
  done
else
  generate "$tap_stream"
  handle_error $?
fi

