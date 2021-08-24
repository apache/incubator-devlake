import { Test, TestingModule } from '@nestjs/testing';
import ChangelogCollector from './changelog.collector';
import ChangelogEntity from './changelog.entity';

describe('Jira/Changelog', () => {
  let app: TestingModule;

  beforeEach(async () => {
    app = await Test.createTestingModule({
      providers: [
        ChangelogCollector,
        {
          provide: ChangelogEntity,
          useFactory: () => {
            return jest.fn();
          },
        },
      ],
    }).compile();
  });

  describe('Collector', () => {
    it('constructor', () => {
      const collector = app.get(ChangelogCollector);
      expect(collector).toBeDefined();
    });
  });
});
