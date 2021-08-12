import { Test } from '@nestjs/testing';
import { SourceController } from './source.controller';
import { SourceService } from './source.service';

describe('SourceController', () => {
  let sourceController: SourceController;

  beforeEach(async () => {
    const app = await Test.createTestingModule({
      controllers: [SourceController],
      providers: [SourceService],
    }).compile();

    sourceController = app.get(SourceController);
  });

  describe('post', () => {
    it('should return source info without options', async () => {
      const response = await sourceController.create({
        type: 'jira',
        options: {
          host: 'https://www.atlassian.com/',
          email: 'xx@example.com',
          auth: 'base64EncodedAuthToken',
        },
      });
      expect(response).toMatchObject({
        type: 'jira',
      });
    });
  });
});
