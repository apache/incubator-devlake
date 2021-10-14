import MappingTag from './MappingTag'

const MappingTagStatus = ({reqValue, resValue, envName, clearBtnReq, clearBtnRes, onChangeReq, onChangeRes}) => {
  return <>
    <MappingTag
      labelName="Rejected"
      labelIntent="danger"
      placeholderText="Add Issue Status..."
      values={reqValue}
      helperText={envName}
      rightElement={clearBtnReq}
      onChange={onChangeReq}
    />
    <MappingTag
      labelName="Resolved"
      labelIntent="success"
      placeholderText="Add Issue Status..."
      values={resValue}
      helperText={envName}
      rightElement={clearBtnRes}
      onChange={onChangeRes}
    />
  </>
}

export default MappingTagStatus
