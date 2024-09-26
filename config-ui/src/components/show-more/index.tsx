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

import { useState } from 'react';
import { CaretDownFilled, CaretRightFilled } from '@ant-design/icons';
import { Flex, Button } from 'antd';

interface Props {
  text?: React.ReactNode;
  btnText: string;
  children: React.ReactNode;
}

export const ShowMore = ({ text, btnText, children }: Props) => {
  const [show, setShow] = useState(false);

  return (
    <div>
      <Flex align="center">
        {text}
        <Button
          type="link"
          size="small"
          icon={show ? <CaretDownFilled /> : <CaretRightFilled />}
          iconPosition="end"
          onClick={() => setShow(!show)}
        >
          {btnText}
        </Button>
      </Flex>
      {show && children}
    </div>
  );
};
