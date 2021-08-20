import { getQueueToken } from '@nestjs/bull';
import { Test, TestingModule } from '@nestjs/testing';
import { DAG } from 'plugins/core/src/dependency.resolver';
import { TasksModule } from './tasks.module';
import { TasksService } from './tasks.services';

jest.mock('ioredis', () => {
  return {
    default: require('ioredis-mock/jest'),
  };
});

const mockedQueue = {
  add: jest.fn(() => {
    return Promise.resolve('ut_test');
  }),
};

describe('TasksModule', () => {
  let app: TestingModule;

  beforeAll(async () => {
    app = await Test.createTestingModule({
      imports: [TasksModule],
    })
      .overrideProvider(getQueueToken('default'))
      .useValue(mockedQueue)
      .compile();
  });

  afterAll(async () => {
    await app.close();
  });

  afterEach(async () => {
    mockedQueue.add.mockClear();
  });

  describe('TasksService', () => {
    it('TasksService should initlized', () => {
      const tasksService = app.get<TasksService>(TasksService);
      expect(tasksService).toBeDefined();
    });

    it('add Task', async () => {
      const tasksService = app.get<TasksService>(TasksService);
      const taskId = await tasksService.startTask(new DAG([]));
      expect(taskId).not.toBeNull();
    });

    it('add Task to Queue', async () => {
      const tasksService = app.get<TasksService>(TasksService);
      await tasksService.startTask(new DAG([{ name: 'Test' }]));
      expect(mockedQueue.add).toBeCalled();
    });
  });
});
