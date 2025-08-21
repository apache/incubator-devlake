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

import { useState, useEffect } from 'react';
import { Modal, Input } from 'antd';

interface Props {
  name: string;
  operating: boolean;
  onCancel: () => void;
  onSubmit: (name: string) => void;
}

export const UpdateNameDialog = ({ operating, onCancel, onSubmit, ...props }: Props) => {
  const [name, setName] = useState('');

  useEffect(() => {
    setName(props.name);
  }, [props.name]);

  return (
    <Modal
      open
      title="Change Blueprint Name"
      okText="Save"
      okButtonProps={{
        disabled: !name || name === props.name,
        loading: operating,
      }}
      onCancel={onCancel}
      onOk={() => onSubmit(name)}
    >
      <h3>Blueprint Name</h3>
      <p>Give your Blueprint a unique name to help you identify it in the future.</p>
      <Input value={name} onChange={(e) => setName(e.target.value)} />
    </Modal>
  );
};
