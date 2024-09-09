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

import { useEffect } from 'react';
import { useNavigate, Outlet } from 'react-router-dom';

import { PageLoading } from '@/components';
import { request as requestVersion, selectVersionStatus, selectVersionError } from '@/features/version';
import { request as requestConnections, selectConnectionsStatus, selectConnectionsError } from '@/features/connections';
import { request as requestOnboard, selectOnboardStatus, selectOnboardError } from '@/features/onboard';
import { useAppDispatch, useAppSelector } from '@/hooks';
import { setUpRequestInterceptor } from '@/utils';

export const App = () => {
  const navigate = useNavigate();

  const dispatch = useAppDispatch();
  const versionStatus = useAppSelector(selectVersionStatus);
  const versionError = useAppSelector(selectVersionError);
  const connectionsStatus = useAppSelector(selectConnectionsStatus);
  const connectionsError = useAppSelector(selectConnectionsError);
  const onboardStatus = useAppSelector(selectOnboardStatus);
  const onboardError = useAppSelector(selectOnboardError);

  useEffect(() => {
    setUpRequestInterceptor(navigate);
    dispatch(requestVersion());
    dispatch(requestConnections());
    dispatch(requestOnboard());
  }, []);

  if (versionStatus === 'loading' || connectionsStatus === 'loading' || onboardStatus === 'loading') {
    return <PageLoading />;
  }

  if (versionStatus === 'failed' || connectionsStatus === 'failed' || onboardStatus === 'failed') {
    throw (versionError as any).message || (connectionsError as any).message || (onboardError as any).message;
  }

  return <Outlet />;
};
