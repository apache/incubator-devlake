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

import { Tour } from 'antd';

import { done as doneFuc, selectOnboard } from '@/features/onboard';
import { useAppDispatch, useAppSelector } from '@/hooks';

interface Props {
  nameRef: React.RefObject<HTMLInputElement>;
  connectionRef: React.RefObject<HTMLInputElement>;
  configRef: React.RefObject<HTMLInputElement>;
}

export const OnboardTour = ({ nameRef, connectionRef, configRef }: Props) => {
  const dispatch = useAppDispatch();
  const { step, done } = useAppSelector(selectOnboard);

  const steps = [
    {
      title: 'This is the project you just created.',
      description: 'Project is the basic management unit',
      target: nameRef.current,
    },
    {
      title: 'A connection is automatically created and associated with the project.',
      description: 'The full connection list can be found at the Connections menu.',
      target: connectionRef.current,
    },
    {
      title: 'Click here to configure project',
      description:
        'You can adjust the data scope,  time range and sync frequency of the project. You can also add scope config to transform the raw data before writing to the database.',
      target: configRef.current,
    },
  ];

  if (step !== 4 || done) {
    return null;
  }

  return <Tour steps={steps} onClose={() => dispatch(doneFuc())} />;
};
