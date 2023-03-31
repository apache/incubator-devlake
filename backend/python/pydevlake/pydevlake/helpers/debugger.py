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

import os

from pydevlake import logger


def init():
    debugger = os.getenv("USE_PYTHON_DEBUGGER", default="").lower()
    if debugger == "":
        return
    # The hostname of the machine from which you're debugging (e.g. your IDE's host).
    host = os.getenv("PYTHON_DEBUG_HOST", default="localhost")
    # The port of the machine from which you're debugging (e.g. your IDE's host)
    port = int(os.getenv("PYTHON_DEBUG_PORT", default=32000))
    print("========== Enabling remote debugging on ", host, ":", port, " ==========")
    if debugger == "pycharm":
        try:
            import pydevd_pycharm as pydevd
            try:
                pydevd.settrace(host=host, port=port, suspend=False, stdoutToServer=True, stderrToServer=True)
                logger.info("Pycharm remote debugger successfully connected")
            except TimeoutError as e:
                logger.error(f"Failed to connect to pycharm debugger on {host}:{port}. Make sure it is running")
        except ImportError as e:
            logger.error("Pycharm debugger library is not installed")
    else:
        logger.error(f"Unsupported Python debugger specified: {debugger}")
