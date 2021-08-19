import { Test, TestingModule } from '@nestjs/testing';
import PluginModule from 'plugins/core/src/plugins.module';
import Task from 'plugins/core/src/task.interface';
import Jira, { JiraOptions } from '../src';

describe('Jira Plugin', () => {
  let app: TestingModule;

  beforeAll(async () => {
    app = await Test.createTestingModule({
      imports: [PluginModule.forRootAsync([Jira])],
    }).compile();
  });

  it('Jira', async () => {
    const jira = await app.resolve<Task>('Jira');
    expect(jira).toBeDefined();
    expect(jira.execute).toBeDefined();
  });
});
