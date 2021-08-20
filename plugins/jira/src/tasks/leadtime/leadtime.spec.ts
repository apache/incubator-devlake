import { Test, TestingModule } from '@nestjs/testing';
import LeadTimeEnricher from './leadtime.enricher';
import IssueEntity from '../issue/issue.entity';
import IssueLeadTimeEntity from './leadtime.entity';

describe('Jira/Issue', () => {
  let app: TestingModule;

  beforeEach(async () => {
    app = await Test.createTestingModule({
      providers: [
        LeadTimeEnricher,
        {
          provide: IssueEntity,
          useFactory: () => {
            return jest.fn();
          },
        },
        {
          provide: IssueLeadTimeEntity,
          useFactory: () => {
            return jest.fn();
          },
        },
      ],
    }).compile();
  });

  describe('Collector', () => {
    it('constructor', () => {
      const enricher = app.get(LeadTimeEnricher);
      expect(enricher).toBeDefined();
    });
  });
});
