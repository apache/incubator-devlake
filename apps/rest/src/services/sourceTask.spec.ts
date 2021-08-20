import { Test } from '@nestjs/testing';
import { SourceTaskService } from './sourceTask';
import { CreateSourceTask, ListSourceTask } from '../types/sourceTask';
import { EntityManager, Repository } from 'typeorm';
import { SourceTask } from '../models';

const mockedEntityManager = new EntityManager(null);

describe('SourceTaskService', () => {
  let sourceTaskService: SourceTaskService;
  let em: EntityManager;

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [
        {
          provide: EntityManager,
          useFactory: () => mockedEntityManager,
        },
        SourceTaskService,
      ],
    }).compile();

    sourceTaskService = app.get(SourceTaskService);
    em = app.get(EntityManager);
  });

  beforeEach(async () => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should return sourceTask', async () => {
      const sourceId = '123'
      const data: CreateSourceTask = {
        collector: ['collector'],
        enricher: ['enricher'],
        options: {
          key: 'myOptions'
        }
      };
      const save = jest.spyOn(em, 'save').mockReturnValue(null);
      const res = await sourceTaskService.create(sourceId, data);
      expect(res).toMatchObject({
        "collector": [
          "collector",
        ],
        "enricher": [
          "enricher",
        ],
        "options": {
          "key": "myOptions",
        }
      });
      expect(save).toBeCalledWith({source_id: sourceId, ...data});
    });
  });

  describe('list', () => {
    it('should list sourceTask', async () => {
      const params: ListSourceTask = {
        page: 2,
        pagesize: 10
      };
      const sourceTaskRepository: Repository<SourceTask> = {
        count: jest.fn().mockReturnValue(Promise.resolve(0)),
        find: jest.fn().mockReturnValue(Promise.resolve([])),
      } as unknown as Repository<SourceTask>;

      jest.spyOn(em, 'getRepository').mockReturnValue(sourceTaskRepository);

      const res = await sourceTaskService.list(params);
      expect(res).toMatchObject({
        total: 0,
        offset: 10,
        data: [],
      });
      expect(sourceTaskRepository.count).toBeCalledWith({});
      expect(sourceTaskRepository.find).toBeCalledWith({
        where: {},
        skip: 10,
        take: 10,
      });
    });
  });

  // describe('get', () => {
  //   it('should get sourceTask', async () => {
  //     const mockedSourceTask = new SourceTask();
  //     mockedSourceTask.name = 'name';
  //     mockedSourceTask.type = 'whatever';
  //     mockedSourceTask.options = {
  //       host: 'https://example.com',
  //     };
  //     const sourceTaskRepository: Repository<SourceTask> = {
  //       findOneOrFail: jest.fn().mockReturnValue(mockedSourceTask),
  //     } as unknown as Repository<SourceTask>;

  //     jest.spyOn(em, 'getRepository').mockReturnValue(sourceTaskRepository);

  //     const res = await sourceTaskService.get('id');

  //     expect(sourceTaskRepository.findOneOrFail).toBeCalledWith('id');
  //     expect(res).toMatchObject(mockedSourceTask);
  //   });
  // });

  // describe('delete', () => {
  //   it('should delete sourceTask', async () => {
  //     const mockedSourceTask = new SourceTask();
  //     mockedSourceTask.name = 'name';
  //     mockedSourceTask.type = 'whatever';
  //     mockedSourceTask.options = {
  //       host: 'https://example.com',
  //     };

  //     jest
  //       .spyOn(sourceTaskService, 'get')
  //       .mockImplementation(() => Promise.resolve(mockedSourceTask));

  //     const emDelete = jest.spyOn(em, 'remove').mockReturnValue(null);
  //     expect(await sourceTaskService.delete('id')).toMatchObject(mockedSourceTask);
  //     expect(emDelete).toBeCalledWith(mockedSourceTask);
  //   });
  // });

  // describe('update', () => {
  //   it('should update sourceTask info', async () => {
  //     const mockedSourceTask = new SourceTask();
  //     mockedSourceTask.name = 'name';
  //     mockedSourceTask.type = 'whatever';
  //     mockedSourceTask.id = 'id';
  //     mockedSourceTask.options = {
  //       host: 'https://example.com',
  //     };

  //     jest
  //       .spyOn(sourceTaskService, 'get')
  //       .mockImplementation(() => Promise.resolve(mockedSourceTask));

  //     const save = jest.spyOn(em, 'save');

  //     const updateParams = {
  //       name: 'newName',
  //       type: 'whatever',
  //       options: {},
  //     };

  //     const res = await sourceTaskService.update('id', updateParams);

  //     expect(save).toBeCalledWith({
  //       id: 'id',
  //       ...updateParams,
  //     });
  //     expect(res).toMatchObject({
  //       id: 'id',
  //       ...updateParams,
  //     });
  //   });
  // });
});
