
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

## How To Add a Plugin

Add the plugin entry extends the Scheduler and implement the execute method.

- create your collector an enricher
```
\\ plugins/myplugin/src/collector/sample.ts

@Injectable({
  scope: Scope.TRANSIENT, //set to TRANSIENT, so collector would not be single instance
})
class SampleCollector implements IExecutable<void> {
  constructor() {}
  async execute(): Promise<void> {
    return;
  }
}

export default SampleCollector;

```

declare the collectors in plugin entry

```
\\ plugins/myplugin/src/index.ts

@Injectable({
  scope: Scope.TRANSIENT,
})
@Collector({
  'Sample': SampleCollector,
})
class MyPlugin extends Scheduler<void> {
   async execute(options: any, contextId: ContextId): Promise<void> {
    const collector = await this.collectorRef.resolve(
      'Sample',
      'Jira',//name Registed in PluginModule
      contextId,
    );
    await collector.execute(options);
    return;
  }
}
```
