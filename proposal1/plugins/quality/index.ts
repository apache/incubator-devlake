import { Plugin, Context } from '../core'
import Issue from '../jira/entities/issue'
import Commit from '../gitlab/entities/commit'


export default class JiraPlugin implements Plugin {

    dependencies: ['jira', 'gitlab']

    async execute(ctx: Context): Promise<void> {
        ctx.log('INFO >>> quality plugin start')

        ctx.log('INFO >>> we can use entities from parent plugins now, like: ', typeof Issue, typeof Commit)

        ctx.log('INFO >>> quality start calculating BUGS COUNT PER 1K LOC')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(50)
        ctx.log('INFO >>> quality end calculating BUGS COUNT PER 1K LOC')

        ctx.log('INFO >>> quality start calculating INCIDENTS COUNT PER 1K LOC')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(100)
        ctx.log('INFO >>> quality end calculating INCIDENTS COUNT PER 1K LOC')

        ctx.log('INFO >>> quality plugin end')
    }
}
