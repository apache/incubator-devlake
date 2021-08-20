import { Test } from '@nestjs/testing';
import { SourceService } from './source';
import { CreateSource, ListSource } from '../types/source';
import { EntityManager, Repository } from 'typeorm';
import { Source } from '../models';

const mockedEntityManager = new EntityManager(null);

describe('SourceService', () => {
  let sourceService: SourceService;
  let em: EntityManager;

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      providers: [
        {
          provide: EntityManager,
          useFactory: () => mockedEntityManager,
        },
        SourceService,
      ],
    }).compile();

    sourceService = app.get(SourceService);
    em = app.get(EntityManager);
  });

  beforeEach(async () => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should return source', async () => {
      const data: CreateSource = {
        type: 'jira',
        options: {
          option1: 'value1',
          option2: 'value2',
        },
      };
      const save = jest.spyOn(em, 'save').mockReturnValue(null);
      const res = await sourceService.create(data);
      expect(res).toMatchObject({ type: 'jira' });
      expect(save).toBeCalledWith(data);
    });
  });

  describe('list', () => {
    it('should list source', async () => {
      const params: ListSource = {
        page: 2,
        pagesize: 10,
        type: 'pluginType',
      };
      const sourceRepository: Repository<Source> = {
        count: jest.fn().mockReturnValue(Promise.resolve(0)),
        find: jest.fn().mockReturnValue(Promise.resolve([])),
      } as unknown as Repository<Source>;

      jest.spyOn(em, 'getRepository').mockReturnValue(sourceRepository);

      const res = await sourceService.list(params);
      expect(res).toMatchObject({
        total: 0,
        offset: 10,
        data: [],
      });
      expect(sourceRepository.count).toBeCalledWith({ type: 'pluginType' });
      expect(sourceRepository.find).toBeCalledWith({
        where: { type: 'pluginType' },
        skip: 10,
        take: 10,
      });
    });
  });

  describe('get', () => {
    it('should get source', async () => {
      const mockedSource = new Source();
      mockedSource.name = 'name';
      mockedSource.type = 'whatever';
      mockedSource.options = {
        host: 'https://example.com',
      };
      const sourceRepository: Repository<Source> = {
        findOneOrFail: jest.fn().mockReturnValue(mockedSource),
      } as unknown as Repository<Source>;

      jest.spyOn(em, 'getRepository').mockReturnValue(sourceRepository);

      const res = await sourceService.get('id');

      expect(sourceRepository.findOneOrFail).toBeCalledWith('id');
      expect(res).toMatchObject(mockedSource);
    });
  });

  describe('delete', () => {
    it('should delete source', async () => {
      const mockedSource = new Source();
      mockedSource.name = 'name';
      mockedSource.type = 'whatever';
      mockedSource.options = {
        host: 'https://example.com',
      };

      jest
        .spyOn(sourceService, 'get')
        .mockImplementation(() => Promise.resolve(mockedSource));

      const emDelete = jest.spyOn(em, 'remove').mockReturnValue(null);
      expect(await sourceService.delete('id')).toMatchObject(mockedSource);
      expect(emDelete).toBeCalledWith(mockedSource);
    });
  });

  describe('update', () => {
    it('should update source info', async () => {
      const mockedSource = new Source();
      mockedSource.name = 'name';
      mockedSource.type = 'whatever';
      mockedSource.id = 'id';
      mockedSource.options = {
        host: 'https://example.com',
      };

      jest
        .spyOn(sourceService, 'get')
        .mockImplementation(() => Promise.resolve(mockedSource));

      const save = jest.spyOn(em, 'save');

      const updateParams = {
        name: 'newName',
        type: 'whatever',
        options: {},
      };

      const res = await sourceService.update('id', updateParams);

      expect(save).toBeCalledWith({
        id: 'id',
        ...updateParams,
      });
      expect(res).toMatchObject({
        id: 'id',
        ...updateParams,
      });
    });
  });
});
