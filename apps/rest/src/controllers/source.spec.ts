import { Test } from '@nestjs/testing';
import Source from '../models/source';
import { SourceService } from '../services/source';
import { SourceController } from './source';

const mockedSourceService = new SourceService(null);

describe('SourceController', () => {
  let sourceController: SourceController;
  let sourceService: SourceService;

  beforeEach(async () => {
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

  describe('post', () => {
    it('should return source info without options', async () => {
      const fn = jest
        .spyOn(sourceService, 'create')
        .mockImplementation(async () => {
          const res = new Source();
          res.type = 'jira';
          res.id = 'id';
          return res;
        });

      const reqCreateSource = {
        type: 'jira' as const,
        options: {
          host: 'https://www.atlassian.com/',
          email: 'xx@example.com',
          auth: 'base64EncodedAuthToken',
        },
      };
      const response = await sourceController.create(reqCreateSource);
      expect(fn).toBeCalledWith(reqCreateSource);
      expect(response).toMatchObject({
        type: 'jira',
        id: 'id',
      });
    });
  });
});
