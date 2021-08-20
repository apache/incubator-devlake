import { Job } from 'bull'
import { Plugin, Context } from './plugins/core'

// hardcoded plugins for demo, should be loaded dynamically
import Jira from './plugins/jira'
import Gitlab from './plugins/gitlab'
import Quality from './plugins/quality'
const plugins: Record<string, Plugin> = {
    jira: new Jira(),
    gitlab: new Gitlab(),
    quality: new Quality()
}

export default async function(job: Job) {
    const { plugin, args } = job.data
    const ctx = new Context(job.id, plugin, args)
    ctx.log('start pipeline job', args)
    await plugins[plugin].execute(ctx)
    ctx.log('end pipeline job', args)
}