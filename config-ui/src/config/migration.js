/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import { DEVLAKE_ENDPOINT } from '@/utils/config'
const MigrationOptions = {
  apiProceedEndpoint: `${DEVLAKE_ENDPOINT}/proceed-db-migration`, // API Get Endpoint
  // NO Api.Get action required for cancel at this time
  apiCancelEndpoint: null,
  apiStatusCode: 428, // API Response Code for Migration Required
  warningId: 'DEVLAKE__MIGRATION_WARNING', // Local Storage Warning ID Key
  cancelToastMessage:
    'Migration Halted - Please downgrade manually, you will continue to receive a warning unless you proceed with migration or rollback.',
  failedToastMessage: 'Database Migration Failed! (Check Network Console)',
  AlertDialog: {
    title: 'New Migration Scripts Detected',
    cancelBtnText: 'Cancel',
    confirmBtnText: 'Proceed to Database Migration',
    confirmRetryBtnText: 'Retry Database Migration',
    continueBtnText: 'Continue'
  }
}

export { MigrationOptions }
