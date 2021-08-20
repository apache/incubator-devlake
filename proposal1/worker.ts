import * as path from 'path'
import { queue } from './plugins/core'


console.log('initializing worker')
queue.process(5, path.join(__dirname, 'worker-process.ts'))
console.log('worker is ready')