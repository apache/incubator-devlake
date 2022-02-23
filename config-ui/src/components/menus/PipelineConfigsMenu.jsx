import React from 'react'
import { Menu } from '@blueprintjs/core'
import { defaultConfig as samplePipelineConfig } from '@/data/pipeline-config-samples/default'
import { refdiffConfig as sampleRefdiffPipelineConfig } from '@/data/pipeline-config-samples/refdiff'
import { gitextractorConfig as sampleGitextractorPipelineConfig } from '@/data/pipeline-config-samples/gitextractor'
import { githubConfig as sampleGithubPipelineConfig } from '@/data/pipeline-config-samples/github'
import { gitlabConfig as sampleGitlabPipelineConfig } from '@/data/pipeline-config-samples/gitlab'
import { jiraConfig as sampleJiraPipelineConfig } from '@/data/pipeline-config-samples/jira'
import { jenkinsConfig as sampleJenkinsPipelineConfig } from '@/data/pipeline-config-samples/jenkins'
import { feishuConfig as sampleFeishuPipelineConfig } from '@/data/pipeline-config-samples/feishu'

const PipelineConfigsMenu = (props) => {
  const {
    setRawConfiguration = () => {},
    advancedMode = false
  } = props
  return (
    <Menu className='pipeline-configs-menu'>
      <label style={{
        fontSize: '10px',
        fontWeight: 800,
        fontFamily: '"Montserrat", sans-serif',
        textTransform: 'uppercase',
        padding: '6px 8px',
        display: 'block'
      }}
      >SAMPLE PIPELINE CONFIGURATIONS
      </label>
      <Menu.Item
        icon='group-objects' text='Load General Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(samplePipelineConfig, null, '  '))}
      />
      <Menu.Item
        icon='group-objects'
        text='Load RefDiff Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(sampleRefdiffPipelineConfig, null, '  '))}
      />
      <Menu.Item
        icon='group-objects' text='Load GitExtractor Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(sampleGitextractorPipelineConfig, null, '  '))}
      />
      <Menu.Item
        icon='group-objects' text='Load GitHub Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(sampleGithubPipelineConfig, null, '  '))}
      />
      <Menu.Item
        icon='group-objects' text='Load GitLab Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(sampleGitlabPipelineConfig, null, '  '))}
      />
      <Menu.Item
        icon='group-objects' text='Load JIRA Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(sampleJiraPipelineConfig, null, '  '))}
      />
      <Menu.Item
        icon='group-objects' text='Load Jenkins Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(sampleJenkinsPipelineConfig, null, '  '))}
      />
      <Menu.Item
        icon='group-objects' text='Load Feishu Configuration'
        onClick={() => setRawConfiguration(JSON.stringify(sampleFeishuPipelineConfig, null, '  '))}
      />
    </Menu>
  )
}

export default PipelineConfigsMenu
