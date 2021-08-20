import { Test } from '@nestjs/testing';
import SourceTask from '../models/sourceTask';
import { SourceTaskService } from '../services/sourceTask';
import { SourceTaskController } from './sourceTask';

const mockedSourceTaskService = new SourceTaskService(null);

describe('SourceTaskController', () => {
  let sourceTaskController: SourceTaskController;
  let sourceTaskService: SourceTaskService;

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      controllers: [SourceTaskController],
      providers: [
        {
          provide: SourceTaskService,
          useFactory: () => mockedSourceTaskService,
        },
      ],
    }).compile();

    sourceTaskController = app.get(SourceTaskController);
    sourceTaskService = app.get(SourceTaskService);
  });

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should return sourceTask info', async () => {
      const sourceId = '123'
      const reqCreateSourceTask = {
        collector: ['my collector'],
        enricher: ['my enricher'],
        options: {}
      };

      const fn = jest
        .spyOn(sourceTaskService, 'create')
        .mockImplementation(async () => {
          const res = new SourceTask();
          res.source_id = sourceId;
          res.collector = reqCreateSourceTask.collector;
          res.enricher = reqCreateSourceTask.enricher;
          res.options = reqCreateSourceTask.options;
          return res;
        });

      const response = await sourceTaskController.create(sourceId, reqCreateSourceTask);
      expect(fn).toBeCalledWith(sourceId, reqCreateSourceTask);
      expect(response).toMatchObject({
        source_id: sourceId,
        ...reqCreateSourceTask,
      });
    });
  });

  describe('list', () => {
    it('should list sourceTasks', async () => {
      const reqListSourceTask = {
        page: 1,
        pagesize: 10,
        source_id: '123'
      };

      const fn = jest
        .spyOn(sourceTaskService, 'list')
        .mockImplementation(async () => {
          return {
            total: 1,
            offset: 0,
            page: 1,
            pagesize: 1,
            data: [],
          };
        });

      await sourceTaskController.list(reqListSourceTask);
    });
  });

  // describe('get', () => {
  //   it('should return target sourceTask', async () => {
  //     const fn = jest
  //       .spyOn(sourceTaskService, 'get')
  //       .mockImplementation(async () => {
  //         return new SourceTask();
  //       });

  //     await sourceTaskController.get('id');
  //   });
  // });

  // describe('update', () => {
  //   it('should update sourceTask', async () => {
  //     const reqUpdateSourceTask = {
  //       type: 'jira',
  //       options: {
  //         host: 'https://www.atlassian.com/',
  //         email: 'xx@example.com',
  //         auth: 'base64EncodedAuthToken',
  //       },
  //     };
  //     const fn = jest
  //       .spyOn(sourceTaskService, 'update')
  //       .mockImplementation(async () => {
  //         return new SourceTask();
  //       });

  //     await sourceTaskController.update('id', reqUpdateSourceTask);
  //   });
  // });

  // describe('delete', () => {
  //   it('should delete sourceTask', async () => {
  //     const fn = jest
  //       .spyOn(sourceTaskService, 'delete')
  //       .mockImplementation(async () => {
  //         return new SourceTask();
  //       });

  //     await sourceTaskController.delete('id');
  //   });
  // });
});
