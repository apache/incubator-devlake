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

import { useRouteError, isRouteErrorResponse } from 'react-router-dom';

import { Logo } from '@/components';

import { ErrorEnum } from './types';
import { Offline, NeedsDBMigrate, Exception } from './components';
import * as S from './styled';

export const Error = () => {
  const error = useRouteError() as Error;

  return (
    <S.Wrapper>
      <Logo />
      <S.Inner>
        {isRouteErrorResponse(error) && error.data.error === ErrorEnum.API_OFFLINE && <Offline />}
        {isRouteErrorResponse(error) && error.data.error === ErrorEnum.NEEDS_DB_MIRGATE && <NeedsDBMigrate />}
        {!isRouteErrorResponse(error) && <Exception error={error} />}
      </S.Inner>
    </S.Wrapper>
  );
};
