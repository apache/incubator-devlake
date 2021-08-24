import { getQueueToken } from '@nestjs/bull';
import { Test, TestingModule } from '@nestjs/testing';
import { DAG } from 'plugins/core/src/dependency.resolver';
import Task from './task.model';
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

  describe('TaskModel', () => {
    it('constructor', async () => {
      const client = app.get('REDIS_TASK_CLIENT');
      const task = new Task('test', client);
      expect(task).toBeDefined();
    });

    it('init empty', async () => {
      const client = app.get('REDIS_TASK_CLIENT');
      const task = new Task('test', client);
      task.init(new DAG([]));
      const jobs = await task.next();
      expect(jobs).toHaveLength(0);
    });

    it('init with job', async () => {
      const client = app.get('REDIS_TASK_CLIENT');
      const task = new Task('test', client);
      task.init(new DAG([{ name: 'UT' }]));
      const jobs = await task.next();
      expect(jobs).toHaveLength(1);
      expect(jobs[0]).toHaveProperty('name', 'UT');
      expect(jobs[0]).toHaveProperty('id');
    });
  });
});
