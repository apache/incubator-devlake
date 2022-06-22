/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migrationscripts

import "github.com/apache/incubator-devlake/migration"

// RegisterAll register all the migration scripts of framework
func All() []migration.Script {
	return []migration.Script{
		new(initSchemas),
		new(updateSchemas20220505), new(updateSchemas20220507), new(updateSchemas20220510),
		new(updateSchemas20220513), new(updateSchemas20220524), new(updateSchemas20220526),
		new(updateSchemas20220527), new(updateSchemas20220528), new(updateSchemas20220601),
		new(updateSchemas20220602), new(updateSchemas20220612), new(updateSchemas20220613),
		new(updateSchemas20220614), new(updateSchemas2022061402), new(updateSchemas20220616),
		new(blueprintNormalMode),
	}
}
