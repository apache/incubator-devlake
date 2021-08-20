import { waitAll } from './plugins/core'


async function main() {
    console.log(new Date().toJSON(), 'start pipeline')

    // hardcoded pipeline for demo, should be generated dynamically
    // pipeline step 1
    await waitAll([
        { plugin: 'jira', args: { boardId: 8, force: false } },
        { plugin: 'gitlab', args: { projectId: 800012, force: false } }
    ])
    // pipeline step 2
    await waitAll([
        { plugin: 'quality', args: { boardId: 8, projectId: 800012, force: false } }
    ])

    console.log(new Date().toJSON(), 'end pipeline')
    process.exit()
}

main()