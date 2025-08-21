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

import { message } from 'antd';

export type OperateConfig = {
  setOperating?: (success: boolean) => void;
  formatMessage?: () => string;
  formatReason?: (err: unknown) => string;
  hideToast?: boolean;
};

/**
 *
 * @param request -> a request
 * @param config
 * @param config.setOperating -> Control the status of the request
 * @param config.formatMessage -> Show the message for the success
 * @param config.formatReason -> Show the reason for the failure
 * @param config.hideToast -> Hide all the toast
 * @returns
 */
export const operator = async <T>(request: () => Promise<T>, config?: OperateConfig): Promise<[boolean, any?]> => {
  const { setOperating, formatMessage, formatReason } = config || {};

  try {
    setOperating?.(true);
    const res = await request();
    const content = formatMessage?.() ?? 'Operation successfully completed';
    if (!config?.hideToast) {
      message.success(content);
    }
    return [true, res];
  } catch (err) {
    console.error('Operation failed.', err);
    const reason = formatReason?.(err) ?? (err as any).response?.data?.message ?? 'Operation failed.';
    message.error(reason);
    return [false, err];
  } finally {
    setOperating?.(false);
  }
};
