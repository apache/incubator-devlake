import { Plugin, Context } from '../core'


export default class JiraPlugin implements Plugin {
    async execute(ctx: Context): Promise<void> {
        ctx.log('INFO >>> gitlab plugin start')

        ctx.log('INFO >>> gitlab start collect repo data')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(10)
        ctx.log('INFO >>> gitlab end collect repo data')

        ctx.log('INFO >>> gitlab start collect commits data')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(50)
        ctx.log('INFO >>> gitlab end collect commits data')

        ctx.log('INFO >>> gitlab start enricher commits data')
        await new Promise(resolve => setTimeout(resolve, 1000))
        ctx.progress(100)
        ctx.log('INFO >>> gitlab end enricher commits data')

        ctx.log('INFO >>> gitlab plugin end')
    }
}
