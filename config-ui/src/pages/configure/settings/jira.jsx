import React, { useEffect, useState } from 'react'
import {
  useParams,
  useHistory
} from 'react-router-dom'
import MappingTag from '@/pages/configure/settings/jira/MappingTag'
import ClearButton from '@/pages/plugins/jira//ClearButton'

import { Button, MenuItem } from '@blueprintjs/core'
import { Select } from '@blueprintjs/select'

import { epicsData } from '@/pages/configure/mock-data/epics'
import { boardsData } from '@/pages/configure/mock-data/boards'
import { granularitiesData } from '@/pages/configure/mock-data/granularities'

import '@/styles/integration.scss'
import '@/styles/connections.scss'
import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

export default function JiraSettings (props) {
  const { connection, provider, isSaving, onSettingsChange } = props
  const { providerId, connectionId } = useParams()
  const history = useHistory()

  const [typeMappingBug, setTypeMappingBug] = useState()
  const [typeMappingIncident, setTypeMappingIncident] = useState()
  const [typeMappingRequirement, setTypeMappingRequirement] = useState()
  const [typeMappingAll, setTypeMappingAll] = useState()
  const [statusMappings, setStatusMappings] = useState()
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState()
  const [jiraIssueStoryCoefficient, setJiraIssueStoryCoefficient] = useState()
  const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState()

  const [selectedEpicItem, setSelectedEpicItem] = useState()
  const [epics, setEpics] = useState(epicsData)

  const [selectedGranularityItem, setSelectedGranularityItem] = useState()
  const [granularities, setGranularities] = useState(granularitiesData)

  const [selectedBoardItem, setSelectedBoardItem] = useState()
  const [boards, setBoards] = useState(boardsData)

  useEffect(() => {
    const settings = {
      JIRA_ISSUE_EPIC_KEY_FIELD: jiraIssueEpicKeyField,
      JIRA_ISSUE_TYPE_MAPPING: typeMappingAll,
      JIRA_ISSUE_STORYPOINT_COEFFICIENT: jiraIssueStoryCoefficient,
      JIRA_ISSUE_STORYPOINT_FIELD: jiraIssueStoryPointField,
      // @todo SET BOARD ID
      // JIRA_ISSUES_BOARD_ID: ??
    }
    onSettingsChange(settings)
    console.log('>> JIRA INSTANCE SETTINGS FIELDS CHANGED!', settings)
    console.log(
      typeMappingBug,
      typeMappingAll,
      typeMappingIncident,
      typeMappingRequirement,
      statusMappings,
      jiraIssueEpicKeyField,
      jiraIssueStoryPointField,
      jiraIssueStoryCoefficient,
      onSettingsChange)
  }, [
    typeMappingBug,
    typeMappingAll,
    typeMappingIncident,
    typeMappingRequirement,
    statusMappings,
    jiraIssueEpicKeyField,
    jiraIssueStoryPointField,
    jiraIssueStoryCoefficient,
    onSettingsChange
  ])

  useEffect(() => {
    if (typeMappingBug && typeMappingIncident && typeMappingRequirement) {
      const typeBug = 'Bug:' + typeMappingBug.toString() + ';'
      const typeIncident = 'Incident:' + typeMappingIncident.toString() + ';'
      const typeRequirement = 'Requirement:' + typeMappingRequirement.toString() + ';'
      const all = typeBug + typeIncident + typeRequirement
      setTypeMappingAll(all)
    }
  }, [typeMappingBug, typeMappingIncident, typeMappingRequirement])

  useEffect(() => {
    // @todo Fetch EPICS, GRANULARITES and BOARDS from API
    console.log('>> CONN SETTINGS OBJECT ', connection)
    // setEpics([])
    // setBoards([])
    // setGranularities([])

    // @todo FETCH & SET INITIAL MAPPING TYPES
    let mappings = {
      Bug: [],
      Incident: [],
      Requirement: []
    }
    if (connection && connection.ID) {
      const types = connection.JIRA_ISSUE_TYPE_MAPPING.split(';').map(t => t.split(':')[0])
      if (types.lastIndexOf('') !== -1) {
        types.pop()
      }
      const tags = connection.JIRA_ISSUE_TYPE_MAPPING.split(';').map(t => t.split(':')[1])
      types.forEach((type, idx) => {
        if (type) {
          mappings = {
            ...mappings,
            [type]: tags[idx] ? tags[idx].split(',') : []
          }
        }
      })

      console.log('>> RE-CREATED ISSUE TYPE MAPPINGS OBJ...', mappings)

      setTypeMappingRequirement(mappings.Requirement)
      setTypeMappingBug(mappings.Bug)
      setTypeMappingIncident(mappings.Incident)
      setStatusMappings([])

      // @todo FETCH & SET EPIC KEY
      const selectedEpic = epics.find(e => e.value === connection.JIRA_ISSUE_EPIC_KEY_FIELD)
      console.log('>>> EPIC ITEM = ', selectedEpic)
      setSelectedEpicItem(selectedEpic)

      // @todo FETCH & SET BOARD ID
      // setSelectedBoardItem(boards.find(b => b.value === connection.JIRA_ISSUES_BOARD_ID???))

      // @todo FETCH & SET GRANULARITY KEY
      setSelectedGranularityItem(granularities.find(g => g.value === connection.JIRA_ISSUE_STORYPOINT_FIELD))
    }
  }, [connection, epics, granularities, boards])

  return (
    <>
      <div className='headlineContainer'>
        <h3 className='headline'>Issue Type Mappings</h3>
        <p className='description'>Map your own issue types to Dev Lake's standard types</p>
      </div>

      <MappingTag
        labelName='Requirement'
        labelIntent='primary'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingRequirement}
        helperText='JIRA_ISSUE_TYPE_MAPPING'
        rightElement={<ClearButton onClick={() => setTypeMappingRequirement([])} />}
        onChange={(values) => setTypeMappingRequirement(values)}
        disabled={isSaving}
      />

      <MappingTag
        labelName='Bug'
        labelIntent='danger'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingBug}
        helperText='JIRA_ISSUE_TYPE_MAPPING'
        rightElement={<ClearButton onClick={() => setTypeMappingBug([])} />}
        onChange={(values) => setTypeMappingBug(values)}
        disabled={isSaving}
      />

      <MappingTag
        labelName='Incident'
        labelIntent='warning'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingIncident}
        helperText='JIRA_ISSUE_TYPE_MAPPING'
        rightElement={<ClearButton onClick={() => setTypeMappingIncident([])} />}
        onChange={(values) => setTypeMappingIncident(values)}
        disabled={isSaving}
      />

      <div className='headlineContainer'>
        <h3 className='headline'>Epic Key <span className='bp3-form-helper-text'>JIRA_ISSUE_EPIC_KEY_FIELD</span></h3>
        <p className=''>Choose the Jira field you’re using to represent the key of an Epic to which an issue belongs to.</p>
        <span style={{ display: 'inline-block' }}>
          <Select
            className='select-epic-key'
            inline={true}
            fill={false}
            items={epics}
            activeItem={selectedEpicItem}
            itemPredicate={(query, item) => item.title.toLowerCase().indexOf(query.toLowerCase()) >= 0}
            itemRenderer={(item, { handleClick, modifiers }) => (
              <MenuItem
                active={modifiers.active}
                key={item.value}
                label={item.value}
                onClick={handleClick}
                text={item.title}
              />
            )}
            noResults={<MenuItem disabled={true} text='No epic results.' />}
            onItemSelect={(item) => {
              // @todo SET/VERIFY ENV FIELD FOR EPIC KEY
              setJiraIssueEpicKeyField(item.value)
              setSelectedEpicItem(item)
            }}
          >
            <Button
              style={{ maxWidth: '260px' }}
              text={selectedEpicItem ? `${selectedEpicItem.title}` : epics[0].title}
              rightIcon='double-caret-vertical'
            />
          </Select>
        </span>
      </div>

      <div className='headlineContainer'>
        <h3 className='headline'>Requirement Granularity  <span className='bp3-form-helper-text'>JIRA_ISSUE_STORYPOINT_FIELD</span></h3>
        <p className=''>Choose the Jira field you’re using to represent the granularity of a requirement-type issue.</p>
        <span style={{ display: 'inline-block' }}>
          <Select
            className='select-granularity-key'
            inline={true}
            fill={false}
            items={granularities}
            activeItem={selectedGranularityItem}
            itemPredicate={(query, item) => item.title.toLowerCase().indexOf(query.toLowerCase()) >= 0}
            itemRenderer={(item, { handleClick, modifiers }) => (
              <MenuItem
                active={modifiers.active}
                key={item.value}
                label={item.value}
                onClick={handleClick}
                text={item.title}
              />
            )}
            noResults={<MenuItem disabled={true} text='No granularity results.' />}
            onItemSelect={(item) => {
              // @todo SET/VERIFY ENV FIELD FOR GRANULARITY
              setJiraIssueStoryCoefficient(item.value)
              setSelectedGranularityItem(item)
            }}
          >
            <Button
              style={{ maxWidth: '260px' }}
              text={selectedGranularityItem ? `${selectedGranularityItem.title}` : granularities[0].title}
              rightIcon='double-caret-vertical'
            />
          </Select>
        </span>
      </div>

      <div className='headlineContainer'>
        <h3 className='headline'>Board ID (Optional) <span className='bp3-form-helper-text'>JIRA_ISSUES_BOARD_ID?</span></h3>
        <p className=''>Choose the specific Jira board(s) to collect issues from.</p>
        <span style={{ display: 'inline-block' }}>
          <Select
            className='select-board-key'
            inline={true}
            fill={false}
            items={boards}
            activeItem={selectedBoardItem}
            itemPredicate={(query, item) => item.title.toLowerCase().indexOf(query.toLowerCase()) >= 0}
            itemRenderer={(item, { handleClick, modifiers }) => (
              <MenuItem
                active={modifiers.active}
                key={item.value}
                label={item.value}
                onClick={handleClick}
                text={item.title}
              />
            )}
            noResults={<MenuItem disabled={true} text='No board results.' />}
            onItemSelect={(item) => {
              // @todo SET/VERIFY ENV FIELD FOR BOARD ID
              setJiraIssueStoryPointField(item.value)
              setSelectedBoardItem(item)
            }}
          >
            <Button
              style={{ maxWidth: '260px' }}
              text={selectedBoardItem ? `${selectedBoardItem.title}` : boards[0].title}
              rightIcon='double-caret-vertical'
            />
          </Select>
        </span>
      </div>

    </>
  )
}
