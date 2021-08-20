import { Plugin, Context } from '../core'


export default class JiraPlugin implements Plugin {
    async execute(ctx: Context): Promise<void> {
        ctx.log('INFO >>> jira plugin start')

        ctx.log('INFO >>> jira start collect board data')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(10)
        ctx.log('INFO >>> jira end collect board data')

        ctx.log('INFO >>> jira start collect issues data')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(50)
        ctx.log('INFO >>> jira end collect issues data')

        ctx.log('INFO >>> jira start enricher issues data')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(100)
        ctx.log('INFO >>> jira end enricher issues data')

        ctx.log('INFO >>> jira plugin end')
    }
}
