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

import styled from 'styled-components';

const Wrapper = styled.div`
  padding: 24px;
  background: #f0f4fe;
  border: 1px solid #bdcefb;
  border-radius: 4px;
`;

interface Props {
  style?: React.CSSProperties;
  content?: React.ReactNode;
  children?: React.ReactNode;
}

export const Alert = ({ style, content, children }: Props) => {
  return <Wrapper style={style}>{content ?? children}</Wrapper>;
};
