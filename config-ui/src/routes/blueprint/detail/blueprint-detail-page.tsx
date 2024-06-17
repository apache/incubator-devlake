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

import { useParams } from 'react-router-dom';
import { Helmet } from 'react-helmet';

import { PageHeader } from '@/components';
import { PATHS } from '@/config';

import { FromEnum } from '../types';

import { BlueprintDetail } from './blueprint-detail';

const brandName = import.meta.env.DEVLAKE_BRAND_NAME ?? 'DevLake';

export const BlueprintDetailPage = () => {
  const { id } = useParams() as { id: string };

  return (
    <PageHeader
      breadcrumbs={[
        { name: 'Advanced', path: PATHS.BLUEPRINTS() },
        { name: 'Blueprints', path: PATHS.BLUEPRINTS() },
        { name: id, path: PATHS.BLUEPRINT(id) },
      ]}
    >
      <Helmet>
        <title>
          {`Blueprints:${id}`} - {brandName}
        </title>
      </Helmet>
      <BlueprintDetail id={id} from={FromEnum.blueprint} />
    </PageHeader>
  );
};
