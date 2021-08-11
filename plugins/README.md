
# Plugins

## How Plugin Start

The plugins is a sequence of tasks. then entry of a plugin would organization the tasks in sequence.

All tasks should run in queue service. So we need to register plugin entry with a name in queue service

```
// <root>/apps/queue/src/consumer/index.ts
export class ConsumerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ConsumerModule,
      imports: [
        BullQueueModule.forRoot(queue),
        PluginModule.Register([{ name: 'Jira', schedule: Jira }]),
        ConsumerService,
      ],
      providers: [],
      exports: [ConsumerService],
    };
  }
}
```

```
// Call Producer Service to add job
import ProducerService from 'apps/queue/src/producer/service';

@Injectable()
export class SomeService {
    constructor(private produer: ProducerService) {}

    methodToAddJob() {
        this.producer.addJob('Jira', {projects: ['MyJiraProject']})
    }
}
```

## How To Create Your Own Plugin

Add the plugin entry extends the Scheduler and implement the execute method.

```
\\ plugins/myplugin/src/index.ts

@Injectable({
  scope: Scope.TRANSIENT,
})
class MyPlugin extends Scheduler<void> {
   async execute(options: JiraOptions): Promise<void> {
    //TODO: schedule the task here
    return;
  }
}
```
