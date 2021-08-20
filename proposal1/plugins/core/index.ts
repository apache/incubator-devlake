import * as Queue from 'bull'

export const queue = new Queue('lake', 'redis://localhost:6379')

export class Context {
    constructor(private jobId: string|number, private pluginName: string, args: any) {
    }

    log(...msg: any[]): void {
        console.log(`[${process.pid}|${this.jobId}][${new Date().toJSON()}] <${this.pluginName}> ${msg[0]}`, ...msg.slice(1))
    }

    progress(percent: number): void {
        this.log('progress: %d', percent)
    }
}

export interface Plugin {
    readonly dependencies?: string[];
    execute?: (job: Context) => Promise<void>;
}

export interface TaskDesc {
    plugin: string,
    args: any
}


export async function waitAll(tasks: TaskDesc[]) {
    // add all tasks to queue
    const jobs = await Promise.all(tasks.map(t => queue.add(t, { attempts: 4 })))
    // wait for all tasks to finish
    await Promise.all(jobs.map(j => j.finished()))
}