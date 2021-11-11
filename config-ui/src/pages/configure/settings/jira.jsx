import React, { useEffect, useState, Fragment } from 'react'
// import {
//   useParams,
//   useHistory
// } from 'react-router-dom'
import {
  FormGroup,
  InputGroup
} from '@blueprintjs/core'
import MappingTag from '@/pages/configure/settings/jira/MappingTag'
import ClearButton from '@/components/ClearButton'
import '@/styles/integration.scss'
import '@/styles/connections.scss'
// import '@blueprintjs/popover2/lib/css/blueprint-popover2.css'

const MAPPING_TYPES = {
  Requirement: 'Requirement',
  Incident: 'Incident',
  Bug: 'Bug'
}

export default function JiraSettings (props) {
  const { connection, provider, isSaving, onSettingsChange } = props
  // const { providerId, connectionId } = useParams()
  // const history = useHistory()

  const [typeMappingBug, setTypeMappingBug] = useState([])
  const [typeMappingIncident, setTypeMappingIncident] = useState([])
  const [typeMappingRequirement, setTypeMappingRequirement] = useState([])
  const [typeMappingAll, setTypeMappingAll] = useState({})
  const [statusMappings, setStatusMappings] = useState()
  const [jiraIssueEpicKeyField, setJiraIssueEpicKeyField] = useState()
  const [jiraIssueStoryCoefficient, setJiraIssueStoryCoefficient] = useState()
  const [jiraIssueStoryPointField, setJiraIssueStoryPointField] = useState()
  // const [epicKey, setEpicKey] = useState()
  // const [granularityKey, setGranularityKey] = useState()
  // const [boardId, setBoardId] = useState()

  // @todo restore when re-enabling selector-based ux
  // ---------------------------------------------------------
  // const [selectedEpicItem, setSelectedEpicItem] = useState()
  // const [epics, setEpics] = useState(epicsData)

  // const [selectedGranularityItem, setSelectedGranularityItem] = useState()
  // const [granularities, setGranularities] = useState(granularitiesData)

  // const [selectedBoardItem, setSelectedBoardItem] = useState()
  // const [boards, setBoards] = useState(boardsData)

  const createTypeMapObject = (customType, standardType) => {
    return customType && standardType
      ? {
          [customType]: {
            standardType
          }
        }
      : null
  }

  const parseTypeMappings = (mappings) => {
    const GroupedMappings = {
      [MAPPING_TYPES.Requirement]: [],
      [MAPPING_TYPES.Incident]: [],
      [MAPPING_TYPES.Bug]: [],
    }
    Object.entries(mappings).forEach(([tag, typeObj]) => {
      GroupedMappings[typeObj.standardType].push(tag)
    })
    console.log('>>>> PARSED TYPE MAPPINGS ....', GroupedMappings)
    setTypeMappingRequirement(GroupedMappings[MAPPING_TYPES.Requirement])
    setTypeMappingBug(GroupedMappings[MAPPING_TYPES.Bug])
    setTypeMappingIncident(GroupedMappings[MAPPING_TYPES.Incident])
    return GroupedMappings
  }

  useEffect(() => {
    const settings = {
      epicKeyField: jiraIssueEpicKeyField,
      typeMappings: typeMappingAll,
      storyPointCoefficient: jiraIssueStoryCoefficient,
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
        ? typeMappingRequirement.toString().split(',').map(r => createTypeMapObject(r, MAPPING_TYPES.Requirement))
        : []
      const IncidentMappings = typeMappingIncident !== ''
        ? typeMappingIncident.toString().split(',').map(i => createTypeMapObject(i, MAPPING_TYPES.Incident))
        : []
      const BugMappings = typeMappingBug !== ''
        ? typeMappingBug.toString().split(',').map(b => createTypeMapObject(b, MAPPING_TYPES.Bug))
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
    // @todo Fetch EPICS, GRANULARITES and BOARDS from API
    console.log('>> CONN SETTINGS OBJECT ', connection)
    // setEpics([])
    // setBoards([])
    // setGranularities([])
    // let mappings = {
    //   Bug: [],
    //   Incident: [],
    //   Requirement: []
    // }
    if (connection && connection.ID) {
      // Parse Type Mappings (V2)
      parseTypeMappings(connection.typeMappings)

      // LEGACY TYPE MAPPINGS (Disabled)
      // const types = connection.JIRA_ISSUE_TYPE_MAPPING ? connection.JIRA_ISSUE_TYPE_MAPPING.split(';').map(t => t.split(':')[0]) : []
      // if (types.lastIndexOf('') !== -1) {
      //   types.pop()
      // }
      // const tags = connection.JIRA_ISSUE_TYPE_MAPPING ? connection.JIRA_ISSUE_TYPE_MAPPING.split(';').map(t => t.split(':')[1]) : []
      // types.forEach((type, idx) => {
      //   if (type) {
      //     mappings = {
      //       ...mappings,
      //       [type]: tags[idx] ? tags[idx].split(',') : []
      //     }
      //   }
      // })
      // console.log('>> RE-CREATED ISSUE TYPE MAPPINGS OBJ...', mappings)
      // setTypeMappingRequirement(mappings.Requirement)
      // setTypeMappingBug(mappings.Bug)
      // setTypeMappingIncident(mappings.Incident)
      setStatusMappings([])
      setJiraIssueEpicKeyField(connection.epicKeyField)
      setJiraIssueStoryCoefficient(connection.storyPointCoefficient)
      setJiraIssueStoryPointField(connection.storyPointField)

      // @todo RE-ENABLE SELECTORS!
      // @todo FETCH & SET EPIC KEY
      // const selectedEpic = epics.find(e => e.value === connection.JIRA_ISSUE_EPIC_KEY_FIELD)
      // console.log('>>> EPIC ITEM = ', selectedEpic)
      // setSelectedEpicItem(selectedEpic)

      // @todo FETCH & SET BOARD ID
      // setSelectedBoardItem(boards.find(b => b.value === connection.JIRA_ISSUES_BOARD_ID???))

      // @todo FETCH & SET GRANULARITY KEY
      // setSelectedGranularityItem(granularities.find(g => g.value === connection.JIRA_ISSUE_STORYPOINT_FIELD))
    }
  }, [connection/*, epics, granularities, boards */])

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
        rightElement={<ClearButton onClick={() => setTypeMappingIncident([])} />}
        onChange={(values) => setTypeMappingIncident(values)}
        disabled={isSaving}
      />

      <div className='headlineContainer'>
        <h3 className='headline'>
          Epic Key<span className='requiredStar'>*</span>
        </h3>
        <p className=''>Choose the Jira field you’re using to represent the key of an Epic to which an issue belongs to.</p>
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
        <h3 className='headline'>Story Point Coefficient
          <span className='requiredStar'>*</span>
        </h3>
        <p className=''>
          This is a number that can convert your jira story points to a new magnitude.&nbsp;
          IE: Convert days to hours with 8 since there are 8 working hours in a day.
        </p>
        {/* <span style={{ display: 'inline-block' }}>
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
        </span> */}
      </div>
      <div className='formContainer' style={{ maxWidth: '250px' }}>
        <FormGroup
          disabled={isSaving}
          // readOnly={['gitlab', 'jenkins'].includes(activeProvider.id)}
          label=''
          inline={true}
          labelFor='granularity-field'
          className='formGroup'
          contentClassName='formGroupContent'
        >
          <InputGroup
            id='granularity-field'
            disabled={isSaving}
            // readOnly={['gitlab', 'jenkins'].includes(activeProvider.id)}
            placeholder='eg. 2000'
            value={jiraIssueStoryCoefficient}
            onChange={(e) => {
              if (e.target.value !== '') {
                const storyPointCoefficientFloat = parseFloat(e.target.value)
                setJiraIssueStoryCoefficient(storyPointCoefficientFloat)
              } else {
                setJiraIssueStoryCoefficient('')
              }
            }}
            className='input granularity-field'
          />
        </FormGroup>
      </div>

      <div className='headlineContainer'>
        <h3 className='headline'>Story Point Field (Optional)</h3>
        <p className=''>Choose the Jira field you’re using to represent the granularity of a requirement-type issue.</p>
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
