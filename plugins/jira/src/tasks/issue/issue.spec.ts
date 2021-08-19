import { Test, TestingModule } from '@nestjs/testing';
import IssueCollector from './issue.collector';
import IssueEntity from './issue.entity';

describe('Jira/Issue', () => {
  let app: TestingModule;

  beforeEach(async () => {
    app = await Test.createTestingModule({
      providers: [
        IssueCollector,
        {
          provide: IssueEntity,
          useFactory: () => {
            return jest.fn();
          },
        },
      ],
    }).compile();
  });

  describe('Collector', () => {
    it('constructor', () => {
      const collector = app.get(IssueCollector);
      expect(collector).toBeDefined();
    });
  });
});
