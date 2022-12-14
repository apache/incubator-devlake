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

import React, { useMemo, useState } from 'react'
import { Link, useHistory } from 'react-router-dom'
import { Button, InputGroup, Checkbox, Intent } from '@blueprintjs/core'

import NoData from '@/images/no-data.svg'
import { PageHeader, Table, ColumnType, Dialog } from '@/components'

import { useProject } from './use-project'
import * as S from './styled'

type ProjectItem = {
  name: string
}

export const ProjectHomePage = () => {
  const [isOpen, setIsOpen] = useState(false)
  const [name, setName] = useState('')
  const [enableDora, setEnableDora] = useState(true)

  const history = useHistory()

  const handleShowDialog = () => setIsOpen(true)
  const handleHideDialog = () => setIsOpen(false)

  const { loading, operating, projects, onSave } = useProject<ProjectItem>({
    name,
    enableDora,
    onHideDialog: handleHideDialog
  })

  const columns = useMemo(
    () =>
      [
        {
          title: 'Project Name',
          dataIndex: 'name',
          key: 'name',
          render: (name: string) => (
            <Link to={`/projects/${name}`} style={{ color: '#292b3f' }}>
              {name}
            </Link>
          )
        },
        {
          title: '',
          dataIndex: 'name' as const,
          align: 'right' as const,
          key: 'action',
          render: (name: any) => (
            <Button
              outlined
              intent={Intent.PRIMARY}
              icon='cog'
              onClick={() => history.push(`/projects/${name}`)}
            />
          )
        }
      ] as ColumnType<ProjectItem>,
    []
  )

  return (
    <PageHeader
      breadcrumbs={[{ name: 'Projects', path: '/projects' }]}
      extra={
        projects.length ? (
          <Button
            intent={Intent.PRIMARY}
            icon='plus'
            text='New Project'
            onClick={handleShowDialog}
          />
        ) : null
      }
    >
      <S.Container>
        {!projects.length ? (
          <S.Inner>
            <div className='logo'>
              <img src={NoData} alt='' />
            </div>
            <div className='desc'>
              <p>
                Add new projects to see engineering metrics based on projects.
              </p>
            </div>
            <div className='action'>
              <Button
                intent={Intent.PRIMARY}
                icon='plus'
                text='New Project'
                onClick={handleShowDialog}
              />
            </div>
          </S.Inner>
        ) : (
          <Table loading={loading} columns={columns} dataSource={projects} />
        )}
        <Dialog
          isOpen={isOpen}
          title='Create a New Project'
          style={{
            top: -100,
            width: 820
          }}
          okText='Save'
          okDisabled={!name}
          okLoading={operating}
          onCancel={handleHideDialog}
          onOk={onSave}
        >
          <S.DialogWrapper>
            <div className='block'>
              <h3>Project Name *</h3>
              <p>Give your project a unique name.</p>
              <InputGroup
                placeholder='Your Project Name'
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>
            <div className='block'>
              <h3>Project Settings</h3>
              <div className='checkbox'>
                <Checkbox
                  label='Enable DORA Metrics'
                  checked={enableDora}
                  onChange={(e) =>
                    setEnableDora((e.target as HTMLInputElement).checked)
                  }
                />
                <p>
                  DORA metrics are four widely-adopted metrics for measuring
                  software delivery performance.
                </p>
              </div>
            </div>
          </S.DialogWrapper>
        </Dialog>
      </S.Container>
    </PageHeader>
  )
}
