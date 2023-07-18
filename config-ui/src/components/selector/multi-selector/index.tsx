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
import { MenuItem, Checkbox, Intent } from '@blueprintjs/core';
import { MultiSelect2 } from '@blueprintjs/select';

interface Props<T> {
  placeholder?: string;
  loading?: boolean;
  items: T[];
  disabledItems?: T[];
  getKey?: (item: T) => string | number;
  getName?: (item: T) => string;
  getIcon?: (item: T) => string;
  selectedItems?: T[];
  onChangeItems?: (selectedItems: T[]) => void;
  noResult?: string;
  onQueryChange?: (query: string) => void;
}

export const MultiSelector = <T,>({
  placeholder,
  loading = false,
  items,
  disabledItems = [],
  getKey = (it) => it as string,
  getName = (it) => it as string,
  getIcon,
  onChangeItems,
  noResult,
  onQueryChange,
  ...props
}: Props<T>) => {
  const [selectedItems, setSelectedItems] = useState<T[]>([]);

  useEffect(() => {
    setSelectedItems(props.selectedItems ?? []);
  }, [props.selectedItems]);

  const tagRenderer = (item: T) => {
    const name = getName(item);
    return <span>{name}</span>;
  };

  const itemRenderer = (item: T, { handleClick }: any) => {
    const key = getKey(item);
    const name = getName(item);
    const icon = getIcon?.(item);
    const selected = !!selectedItems.find((it) => getKey(it) === key);
    const disabled = !!disabledItems.find((it) => getKey(it) === key);

    return (
      <MenuItem
        key={key}
        disabled={selected || disabled}
        onClick={(e) => {
          e.preventDefault();
          handleClick();
        }}
        labelElement={icon ? <img src={icon} width={16} alt="" /> : null}
        text={<Checkbox disabled={selected || disabled} checked={selected || disabled} readOnly label={name} />}
      />
    );
  };

  const handleItemSelect = (item: T) => {
    const newSelectedItems = [...selectedItems, item];
    if (onChangeItems) {
      onChangeItems(newSelectedItems);
    } else {
      setSelectedItems(newSelectedItems);
    }
  };

  const handleItemRemove = (item: T) => {
    const newSelectedItems = selectedItems.filter((it) => getKey(it) !== getKey(item));
    if (onChangeItems) {
      onChangeItems(newSelectedItems);
    } else {
      setSelectedItems(newSelectedItems);
    }
  };

  return (
    <MultiSelect2
      fill
      placeholder={placeholder ?? 'Select...'}
      items={items}
      // https://github.com/palantir/blueprint/issues/3596
      // set activeItem to null will fixed the scrollBar to top when the selectedItems changed
      activeItem={null}
      selectedItems={selectedItems}
      itemRenderer={itemRenderer}
      tagRenderer={tagRenderer}
      tagInputProps={{
        tagProps: {
          intent: Intent.PRIMARY,
          minimal: true,
        },
      }}
      onItemSelect={handleItemSelect}
      onRemove={handleItemRemove}
      onQueryChange={onQueryChange}
      noResults={<MenuItem disabled={true} text={loading ? 'Fetching...' : noResult} />}
    />
  );
};
