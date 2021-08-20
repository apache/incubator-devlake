import { Test } from '@nestjs/testing';
import Source from '../models/source';
import { SourceService } from '../services/source';
import { SourceController } from './source';

const mockedSourceService = new SourceService(null);

describe('SourceController', () => {
  let sourceController: SourceController;
  let sourceService: SourceService;

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      controllers: [SourceController],
      providers: [
        {
          provide: SourceService,
          useFactory: () => mockedSourceService,
        },
      ],
    }).compile();

    sourceController = app.get(SourceController);
    sourceService = app.get(SourceService);
  });

  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should return source info', async () => {
      const reqCreateSource = {
        type: 'jira',
        options: {
          host: 'https://www.atlassian.com/',
          email: 'xx@example.com',
          auth: 'base64EncodedAuthToken',
        },
      };

      const fn = jest
        .spyOn(sourceService, 'create')
        .mockImplementation(async () => {
          const res = new Source();
          res.type = reqCreateSource.type;
          res.id = 'id';
          res.options = reqCreateSource.options;
          return res;
        });

      const response = await sourceController.create(reqCreateSource);
      expect(fn).toBeCalledWith(reqCreateSource);
      expect(response).toMatchObject({
        id: 'id',
        ...reqCreateSource,
      });
    });
  });

  describe('list', () => {
    it('should list sources', async () => {
      const reqListSource = {
        page: 1,
        pagesize: 10,
        type: 'jira',
      };

      const fn = jest
        .spyOn(sourceService, 'list')
        .mockImplementation(async () => {
          return {
            page: 1,
            offset: 0,
            data: [],
          };
        });

      await sourceController.list(reqListSource);
    });
  });

  describe('get', () => {
    it('should return target source', async () => {
      const fn = jest
        .spyOn(sourceService, 'get')
        .mockImplementation(async () => {
          return new Source();
        });

      await sourceController.get('id');
    });
  });

  describe('update', () => {
    it('should update source', async () => {
      const reqUpdateSource = {
        type: 'jira',
        options: {
          host: 'https://www.atlassian.com/',
          email: 'xx@example.com',
          auth: 'base64EncodedAuthToken',
        },
      };
      const fn = jest
        .spyOn(sourceService, 'update')
        .mockImplementation(async () => {
          return new Source();
        });

      await sourceController.update('id', reqUpdateSource);
    });
  });

  describe('delete', () => {
    it('should delete source', async () => {
      const fn = jest
        .spyOn(sourceService, 'delete')
        .mockImplementation(async () => {
          return new Source();
        });

      await sourceController.delete('id');
    });
  });
});
