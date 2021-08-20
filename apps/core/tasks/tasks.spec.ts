import { Test, TestingModule } from '@nestjs/testing';
import { TasksModule } from './tasks.module';
import { TasksService } from './tasks.services';

jest.mock('ioredis', () => {
  return {
    default: require('ioredis-mock/jest'),
  };
});

describe('TasksModule', () => {
  let app: TestingModule;

  beforeAll(async () => {
    app = await Test.createTestingModule({
      imports: [TasksModule],
    }).compile();
  });

  afterAll(async () => {
    await app.close();
  });

  describe('TasksService', () => {
    it('TasksService should initlized', () => {
      const tasksService = app.get<TasksService>(TasksService);
      expect(tasksService).toBeDefined();
    });
  });
});
