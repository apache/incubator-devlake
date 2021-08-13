import { ConfigModule } from '@nestjs/config';
import { Test } from '@nestjs/testing';
import Source from '../models/source';
import CustomTypeOrmModule from '../providers/typeorm.module';
import { SourceService } from '../services/source';
import { SourceController } from './source';

describe('SourceController', () => {
  let sourceController: SourceController;
  let sourceService: SourceService;

  beforeEach(async () => {
    const app = await Test.createTestingModule({
      imports: [
        ConfigModule.forRoot({ isGlobal: true, envFilePath: '.env.test' }),
        CustomTypeOrmModule.forRootAsync(null),
      ],
      controllers: [SourceController],
      providers: [SourceService],
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
