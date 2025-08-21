# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


import sys
import logging

CRITICAL = logging.CRITICAL
FATAL = logging.FATAL
ERROR = logging.ERROR
WARNING = logging.WARNING
INFO = logging.INFO
DEBUG = logging.DEBUG
NOTSET = logging.NOTSET

# mappings from main Go server to Python logging levels
log_levels = {
    'debug': DEBUG,
    'info': INFO,
    'warn': WARNING,
    'error': ERROR,
}

stdout_handler = logging.StreamHandler(sys.stdout)
stdout_handler.addFilter(lambda rec: rec.levelno < logging.ERROR)

stderr_handler = logging.StreamHandler(sys.stderr)
stderr_handler.addFilter(lambda rec: rec.levelno >= logging.ERROR)

logging.basicConfig(
    level=INFO, # default
    format='%(levelname)s: %(message)s',
    handlers=[stdout_handler, stderr_handler]
)

logger = logging.getLogger()
