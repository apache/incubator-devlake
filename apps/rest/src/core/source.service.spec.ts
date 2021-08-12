import { Test } from '@nestjs/testing';
import { SourceService } from './source.service';
import { CreateSource } from './source.type';

describe('SourceService', () => {
  let sourceService: SourceService;

  beforeEach(async () => {
    const app = await Test.createTestingModule({
      providers: [SourceService],
    }).compile();

    sourceService = app.get(SourceService);
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
      const res = await sourceService.create(data);
      expect(res).toMatchObject({ type: 'jira' });
    });
  });
});
