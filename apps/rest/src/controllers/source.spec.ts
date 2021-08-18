import { Test } from '@nestjs/testing';
import Source from '../models/source';
import { SourceController } from './source';

describe('SourceController', () => {
  let sourceController: SourceController;

  beforeAll(async () => {
    const app = await Test.createTestingModule({
      controllers: [SourceController],
      providers: [],
    }).compile();

    sourceController = app.get(SourceController);
  });

  describe('create', () => {
    it('should return source info', async () => {
      const reqCreateSource = {
        type: 'jira' as const,
        options: {
          host: 'https://www.atlassian.com/',
          email: 'xx@example.com',
          auth: 'base64EncodedAuthToken',
        },
      };

      await sourceController.create(reqCreateSource);
    });
  });

  describe('list', () => {
    it('should list sources', async () => {
      const reqListSource = {
        page: 1,
        pagesize: 10,
        type: 'jira' as const,
      };
      await sourceController.list(reqListSource);
    });
  });

  describe('get', () => {
    it('should return target source', async () => {
      await sourceController.get('id');
    });
  });

  describe('update', () => {
    it('should update source', async () => {
      const reqUpdateSource = {
        type: 'jira' as const,
        options: {
          host: 'https://www.atlassian.com/',
          email: 'xx@example.com',
          auth: 'base64EncodedAuthToken',
        },
      };
      await sourceController.update('id', reqUpdateSource);
    });
  });

  describe('delete', () => {
    it('should delete source', async () => {
      await sourceController.delete('id');
    });
  });
});
