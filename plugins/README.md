# How plugin work?

- 1. Plugin Worker is triggered by Apps. Use queue producer to push a plugin task into queue with plugin name and data
```
producer.add('Jira', {host:'', username: '', password: ''});
```

- 2. Plugin is Loaded in Queue Consumer. the nestjs IOC framework would help to get a instance of Plugin
```
moduleRef.resolve('Jira', contextId, {strict: false});
```

- 3. Plugin resolve self with a Task DAG.

- 4. A Task Service will read the Task DAG and start a Session to manage tasks execution

- 5. Task Service push task into queue service with task name and task data
```
producer.add('JiraIssueCollector', {host:'', username: '', password: ''});
```

- 6. Task Service recived the task completed events and check if has next task. the do step 5. if no task anymore finished the pipline
```
on('job:finished', (jobId, datas) => {  producer.add('JiraLeadtimeEnricher', datas); });
```

# How plugin dependency described

- Plugin is the organizer, Described what Entities would exports.

- Entities may have one Task to exports. a task is one of collector or enricher.

- The task would required some Entities, that should Import in Task

```
@Imports([RequiredEntity, AnotherReuqired.Entity])
@Exports(SampleEntity)
class SampleEnricher implements Task {
    async execute(): Promise<any> {
        //do enrich here
        return {};
    }
}
```

- Only one Entity should be Exports in Task

- Dependency Resolver would get the Task DAG for task exexution managment.

## 推荐的文件目录是什么样的
## recommended file directory
可以以如下的格式来创建一个插件，以下结构core会自动注册一个名为「example」，入口代码定义在`index.ts`的一个插件。

You can create a plug-in in the following format. <br>
The core of the following structure will automatically register a plugin named "example" whose entry code is defined in `index.ts`.
```text
plugins/
    example/
        src/
            index.ts
            migrations/
            entities/
            ……
        test/
            ……
```

- All migrates should put at /plugins/YourPlugin/src/migrations, generator your migrate file with typeorm
```typeorm migrate:create -d /plugins/YourPlugin/src/migrations -n mymigrate````
