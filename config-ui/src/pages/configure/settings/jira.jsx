import React, { useEffect, useState, Fragment } from 'react'
import request from '@/utils/request'
import {
  FormGroup,
  InputGroup,
  MenuItem,
  Button,
  Intent,
  Icon,
  Colors
} from '@blueprintjs/core'
import useJIRA from '@/hooks/useJIRA'
import { MultiSelect } from '@blueprintjs/select'
import MappingTag from '@/pages/configure/settings/jira/MappingTag'
import ClearButton from '@/components/ClearButton'
import '@/styles/integration.scss'
import '@/styles/connections.scss'

const MAPPING_TYPES = {
  Requirement: 'Requirement',
  Incident: 'Incident',
  Bug: 'Bug'
}

export default function JiraSettings (props) {
  const { connection, provider, isSaving, onSettingsChange } = props
  // const { providerId, connectionId } = useParams()
  // const history = useHistory()

  const API_PROXY_ENDPOINT = `/plugins/jira/sources/${connection?.ID}/proxy/rest`
  const ISSUE_TYPES_ENDPOINT = `${API_PROXY_ENDPOINT}/api/3/issuetype`
  const ISSUE_FIELDS_ENDPOINT = `${API_PROXY_ENDPOINT}/api/3/field`

  const { fetchIssueTypes, fetchFields, issueTypes, fields } = useJIRA({
    apiProxyPath: API_PROXY_ENDPOINT,
    issuesEndpoint: ISSUE_TYPES_ENDPOINT,
    fieldsEndpoint: ISSUE_FIELDS_ENDPOINT
  })

  const [typeMappingBug, setTypeMappingBug] = useState([])
  const [typeMappingIncident, setTypeMappingIncident] = useState([])
  const [typeMappingRequirement, setTypeMappingRequirement] = useState([])
  const [typeMappingAll, setTypeMappingAll] = useState({})
  const [statusMappings, setStatusMappings] = useState()
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState()
  const [jiraIssueStoryCoefficient, setJiraIssueStoryCoefficient] = useState(1)
  const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState()

  const [requirementTags, setRequirementTags] = useState([])
  const [bugTags, setBugTags] = useState([])
  const [incidentTags, setIncidentTags] = useState([])

  // @todo: remove tags list initial state mock data
  const [requirementTagsList, setRequirementTagsList] = useState([
    { id: 0, title: 'REQ TAG 100', value: 'REQ-100' },
    { id: 1, title: 'REQ TAG 200', value: 'REQ-200' },
    { id: 2, title: 'REQ TAG 300', value: 'REQ-300' }
  ])
  const [bugTagsList, setBugTagsList] = useState([
    { id: 0, title: 'BUG-100', value: 'BUG-100' },
    { id: 1, title: 'BUG-200', value: 'BUG-200' },
    { id: 2, title: 'BUG-300', value: 'BUG-300' }
  ])
  const [incidentTagsList, setIncidentTagsList] = useState([
    { id: 0, title: 'INCIDENT TAG 100', value: 'INCIDENT-100' },
    { id: 1, title: 'INCIDENT TAG 200', value: 'INCIDENT-200' },
    { id: 2, title: 'INCIDENT TAG 300', value: 'INCIDENT-300' }
  ])

  const [fieldsList, setFieldsList] = useState([
    { id: 0, title: 'development', value: 'development', type: 'Dev Summary Custom Field' },
    { id: 1, title: 'Epic Color', value: 'Epic Color', type: 'Color of Epic' },
    { id: 2, title: 'Epic Link', value: 'Epic Link', type: 'Epic Link Relationship' },
    { id: 3, title: 'Epic Status', value: 'Epic Status', type: 'Status of Epic' },
    { id: 4, title: 'Flagged', value: 'Flagged', type: 'Checkboxes' },
    { id: 5, title: 'Sprint', value: 'Sprint', type: 'Jira Sprint Field' },
    { id: 6, title: 'Team', value: 'Team', type: 'Team' },
    { id: 7, title: 'uuid', value: 'Flagged', type: 'UUID Field' },
    { id: 8, title: 'Rank', value: 'Rank', type: 'Global Rank' },
  ])

  const createTypeMapObject = (customType, standardType) => {
    return customType && standardType
      ? {
          [customType]: {
            standardType
          }
        }
      : null
  }

  const parseTypeMappings = (mappings = []) => {
    const GroupedMappings = {
      [MAPPING_TYPES.Requirement]: [],
      [MAPPING_TYPES.Incident]: [],
      [MAPPING_TYPES.Bug]: [],
    }
    Object.entries(mappings).forEach(([tag, typeObj]) => {
      GroupedMappings[typeObj.standardType].push(tag)
    })
    console.log('>>>> PARSED TYPE MAPPINGS ....', GroupedMappings)
    // @todo: fix parsed type mappings w/ list filters
    setTypeMappingRequirement(GroupedMappings[MAPPING_TYPES.Requirement])
    setTypeMappingBug(GroupedMappings[MAPPING_TYPES.Bug])
    setTypeMappingIncident(GroupedMappings[MAPPING_TYPES.Incident])
    setRequirementTags(requirementTagsList.filter(t => GroupedMappings[MAPPING_TYPES.Requirement].includes(t.value)))
    setBugTags(bugTagsList.filter(t => GroupedMappings[MAPPING_TYPES.Bug].includes(t.value)))
    setIncidentTags(incidentTagsList.filter(t => GroupedMappings[MAPPING_TYPES.Incident].includes(t.value)))
    return GroupedMappings
  }

  useEffect(() => {
    const settings = {
      epicKeyField: jiraIssueEpicKeyField,
      typeMappings: typeMappingAll,
      storyPointField: jiraIssueStoryPointField,
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
      // LEGACY MAPPING FORMAT (DISABLED)
      // const typeBug = 'Bug:' + typeMappingBug.toString() + ';'
      // const typeIncident = 'Incident:' + typeMappingIncident.toString() + ';'
      // const typeRequirement = 'Requirement:' + typeMappingIncident.toString() + ';'
      // const all = typeBug + typeIncident + typeRequirement
      // setTypeMappingAll(all)
      const RequirementMappings = typeMappingRequirement !== ''
        ? typeMappingRequirement.map(r => createTypeMapObject(r.value, MAPPING_TYPES.Requirement))
        : []
      const IncidentMappings = typeMappingIncident !== ''
        ? typeMappingIncident.map(i => createTypeMapObject(i.value, MAPPING_TYPES.Incident))
        : []
      const BugMappings = typeMappingBug !== ''
        ? typeMappingBug.map(b => createTypeMapObject(b.value, MAPPING_TYPES.Bug))
        : []
      const CombinedMappings = [...RequirementMappings, ...IncidentMappings, ...BugMappings].filter(m => m !== null)
      const MappingTypeObjects = CombinedMappings.reduce((pV, cV) => { return { ...cV, ...pV } }, {})
      setTypeMappingAll(MappingTypeObjects)
      console.log('>> INCIDENT TYPE MAPPING OBJECTS....', RequirementMappings, IncidentMappings, BugMappings)
      console.log('>> ALL MAPPINGS COMBINED...', CombinedMappings)
      console.log('>> FINAL MAPPING OBJECTS FOR API REQUEST...', MappingTypeObjects)
    }
  }, [typeMappingBug, typeMappingIncident, typeMappingRequirement])

  useEffect(() => {
    console.log('>> CONN SETTINGS OBJECT ', connection)
    if (connection && connection.ID) {
      // Parse Type Mappings (V2)
      parseTypeMappings(connection.typeMappings)
      setStatusMappings([])
      setJiraIssueEpicKeyField(connection.epicKeyField)
      setJiraIssueStoryPointField(connection.storyPointField)
    }
  }, [connection/*, epics, granularities, boards */])

  useEffect(() => {
    setTypeMappingRequirement(requirementTags)
  }, [requirementTags])

  useEffect(() => {
    setTypeMappingBug(bugTags)
  }, [bugTags])

  useEffect(() => {
    setTypeMappingIncident(incidentTags)
  }, [incidentTags])

  useEffect(() => {
    // @todo: set issue types lists from proxy api response data
    // @todo: set issue fields list from proxy api response data
    fetchIssueTypes()
    fetchFields()
  }, [])

  return (
    <>
      <div className='headlineContainer'>
        <h3 className='headline'>Issue Type Mappings</h3>
        <p>Map your own issue types to <strong>DevLake's</strong> standard types</p>
      </div>

      <div className='issue-type-multiselect' style={{ display: 'flex', marginBottom: '10px' }}>
        <div className='issue-type-label' style={{ minWidth: '150px', paddingRight: '10px', paddingTop: '3px' }}>
          <span
            className='bp3-tag tag-requirement'
            style={{ float: 'right' }}
          ><span className='bp3-fill bp3-text-overflow-ellipsis'>Requirement</span>
          </span>
        </div>
        <div className='issue-type-multiselect-selector' style={{ minWidth: '200px', width: '50%' }}>
          <MultiSelect
            placeholder='< Select one or more Requirement Tags >'
            popoverProps={{ usePortal: false, minimal: true, fill: true, style: { width: '100%' } }}
            className='multiselector-requirement-type'
            inline={true}
            fill={true}
            items={requirementTagsList}
            selectedItems={requirementTags}
            activeItem={null}
            itemPredicate={(query, item) => item.title.toLowerCase().indexOf(query.toLowerCase()) >= 0}
            itemRenderer={(item, { handleClick, modifiers }) => (
              <MenuItem
                active={modifiers.active || requirementTags.includes(item)}
                disabled={requirementTags.includes(item)}
                key={item.value}
                label={item.value}
                onClick={handleClick}
                text={requirementTags.includes(item) ? (<>{item.title} <Icon icon='small-tick' color={Colors.GREEN5} /></>) : item.title}
                style={{ marginBottom: '2px', fontWeight: requirementTags.includes(item) ? 'normal' : '700' }}
              />
            )}
            tagRenderer={(item) => item.title}
            tagInputProps={{
              tagProps: {
                intent: Intent.NONE,
                color: Colors.RED3,
                minimal: true
              },
            }}
            noResults={<MenuItem disabled={true} text='No results.' />}
            onRemove={(item) => {
              setRequirementTags((rT) => rT.filter(t => t.id !== item.id))
            }}
            onItemSelect={(item) => {
              setRequirementTags((rT) => !rT.includes(item) ? [...rT, item] : [...rT])
            }}
          />
        </div>
        <div className='multiselect-clear-action' style={{ marginLeft: '5px' }}>
          <ClearButton
            disabled={requirementTags.length === 0}
            intent={Intent.PRIMARY} minimal={false} onClick={() => setRequirementTags([])}
          />
        </div>
      </div>

      <div className='issue-type-multiselect' style={{ display: 'flex', marginBottom: '10px' }}>
        <div className='issue-type-label' style={{ minWidth: '150px', paddingRight: '10px', paddingTop: '3px' }}>
          <span
            className='bp3-tag tag-bug'
            style={{ float: 'right' }}
          ><span className='bp3-fill bp3-text-overflow-ellipsis'>Bug</span>
          </span>
        </div>
        <div className='issue-type-multiselect-selector' style={{ minWidth: '200px', width: '50%' }}>
          <MultiSelect
            placeholder='< Select one or more Bug Tags >'
            popoverProps={{ usePortal: false, minimal: true }}
            className='multiselector-bug-type'
            inline={true}
            fill={true}
            items={bugTagsList}
            selectedItems={bugTags}
            activeItem={null}
            itemPredicate={(query, item) => item.title.toLowerCase().indexOf(query.toLowerCase()) >= 0}
            itemRenderer={(item, { handleClick, modifiers }) => (
              <MenuItem
                active={modifiers.active || bugTags.includes(item)}
                disabled={bugTags.includes(item)}
                key={item.value}
                label={item.value}
                onClick={handleClick}
                text={bugTags.includes(item) ? (<>{item.title} <Icon icon='small-tick' color={Colors.GREEN5} /></>) : item.title}
                style={{ marginBottom: '2px', fontWeight: bugTags.includes(item) ? 'normal' : '700' }}
              />
            )}
            tagRenderer={(item) => item.title}
            tagInputProps={{
              tagProps: {
                intent: Intent.NONE,
                color: Colors.RED3,
                minimal: true
              },
            }}
            noResults={<MenuItem disabled={true} text='No results.' />}
            onRemove={(item) => {
              setBugTags((rT) => rT.filter(t => t.id !== item.id))
            }}
            onItemSelect={(item) => {
              setBugTags((rT) => !rT.includes(item) ? [...rT, item] : [...rT])
            }}
          />
        </div>
        <div className='multiselect-clear-action' style={{ marginLeft: '5px' }}>
          <ClearButton
            disabled={bugTags.length === 0}
            intent={Intent.PRIMARY} minimal={false} onClick={() => setBugTags([])}
          />
        </div>
      </div>

      <div className='issue-type-multiselect' style={{ display: 'flex', marginBottom: '10px' }}>
        <div className='issue-type-label' style={{ minWidth: '150px', paddingRight: '10px', paddingTop: '3px' }}>
          <span
            className='bp3-tag tag-incident'
            style={{ float: 'right' }}
          ><span className='bp3-fill bp3-text-overflow-ellipsis'>Incident</span>
          </span>
        </div>
        <div className='issue-type-multiselect-selector' style={{ minWidth: '200px', width: '50%' }}>
          <MultiSelect
            placeholder='< Select one or more Incident Tags >'
            popoverProps={{ usePortal: false, minimal: true }}
            className='multiselector-incident-type'
            inline={true}
            fill={true}
            items={incidentTagsList}
            selectedItems={incidentTags}
            activeItem={null}
            itemPredicate={(query, item) => item.title.toLowerCase().indexOf(query.toLowerCase()) >= 0}
            itemRenderer={(item, { handleClick, modifiers }) => (
              <MenuItem
                active={modifiers.active || incidentTags.includes(item)}
                disabled={incidentTags.includes(item)}
                key={item.value}
                label={item.value}
                onClick={handleClick}
                text={incidentTags.includes(item) ? (<>{item.title} <Icon icon='small-tick' color={Colors.GREEN5} /></>) : item.title}
                style={{ marginBottom: '2px', fontWeight: incidentTags.includes(item) ? 'normal' : '700' }}
              />
            )}
            tagRenderer={(item) => item.title}
            tagInputProps={{
              tagProps: {
                intent: Intent.NONE,
                color: Colors.RED3,
                minimal: true
              },
            }}
            noResults={<MenuItem disabled={true} text='No results.' />}
            onRemove={(item) => {
              setIncidentTags((rT) => rT.filter(t => t.id !== item.id))
            }}
            onItemSelect={(item) => {
              setIncidentTags((rT) => !rT.includes(item) ? [...rT, item] : [...rT])
            }}
          />
        </div>
        <div className='multiselect-clear-action' style={{ marginLeft: '5px' }}>
          <ClearButton
            disabled={incidentTags.length === 0}
            intent={Intent.PRIMARY} minimal={false} onClick={() => setIncidentTags([])}
          />
        </div>
      </div>

      {/* <MappingTag
        labelName='Requirement'
        classNames='tag-requirement'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingRequirement}
        rightElement={<ClearButton onClick={() => setTypeMappingRequirement([])} />}
        onChange={(values) => setTypeMappingRequirement(values)}
        disabled={isSaving}
      /> */}

      {/* <MappingTag
        labelName='Bug'
        classNames='tag-bug'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingBug}
        rightElement={<ClearButton onClick={() => setTypeMappingBug([])} />}
        onChange={(values) => setTypeMappingBug(values)}
        disabled={isSaving}
      /> */}

      {/* <MappingTag
        labelName='Incident'
        classNames='tag-incident'
        typeOrStatus='type'
        placeholderText='Add Issue Types...'
        values={typeMappingIncident}
        rightElement={<ClearButton onClick={() => setTypeMappingIncident([])} />}
        onChange={(values) => setTypeMappingIncident(values)}
        disabled={isSaving}
      /> */}

      <div className='headlineContainer'>
        <h3 className='headline'>
          Epic Key<span className='requiredStar'>*</span>
        </h3>
        <p className=''>Choose the JIRA field you’re using to represent the key of an Epic to which an issue belongs to.</p>
        {/* <span style={{ display: 'inline-block' }}>
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
        </span> */}
      </div>
      <div className='formContainer' style={{ maxWidth: '250px' }}>
        <FormGroup
          disabled={isSaving}
          // readOnly={['gitlab', 'jenkins'].includes(activeProvider.id)}
          label=''
          inline={true}
          labelFor='epic-key-field'
          // helperText='NAME'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <InputGroup
            id='epic-key-field'
            disabled={isSaving}
            // readOnly={['gitlab', 'jenkins'].includes(activeProvider.id)}
            placeholder='eg. 1000'
            value={jiraIssueEpicKeyField}
            onChange={(e) => setJiraIssueEpicKeyField(e.target.value)}
            className='input epic-key-field'
          />
        </FormGroup>
      </div>
      <div className='headlineContainer'>
        <h3 className='headline'>Story Point Field (Optional)</h3>
        <p className=''>Choose the JIRA field you’re using to represent the granularity of a requirement-type issue.</p>
        {/* <span style={{ display: 'inline-block' }}>
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
        </span> */}
      </div>
      <div className='formContainer' style={{ maxWidth: '250px' }}>
        <FormGroup
          disabled={isSaving}
          label=''
          inline={true}
          labelFor='board-id-field'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <InputGroup
            id='board-id'
            disabled={isSaving}
            placeholder='eg. 3000'
            value={jiraIssueStoryPointField}
            onChange={(e) => setJiraIssueStoryPointField(e.target.value)}
            className='input board-id'
          />
        </FormGroup>
      </div>
    </>
  )
}
