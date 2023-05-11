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

export const Wrapper = styled.div`
  .top {
    display: flex;
    justify-content: space-between;

    h3 {
      margin-bottom: 16px;
    }
  }

  .authentication {
    h3 {
      span 
    }
  }

  .action {
    margin-top: 36px;
    margin-bottom: 24px;
  }
`;

export const DialogTitle = styled.div`
  display: flex;
  align-items: center;

  img {
    margin-right: 8px;
    width: 24px;
  }
`;

export const DialogBody = styled.div`
  display: flex;
  align-items: center;

  .bp4-icon {
    margin-right: 8px;
    color: #f4be55;
  }
`;
