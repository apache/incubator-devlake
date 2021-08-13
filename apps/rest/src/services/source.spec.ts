import { Test } from '@nestjs/testing';
import { SourceService } from './source';
import { CreateSource } from '../types/source';
import { EntityManager } from 'typeorm';

const mockedEntityManager = new EntityManager(null);

describe('SourceService', () => {
  let sourceService: SourceService;
  let em: EntityManager;

  beforeEach(async () => {
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
      // jest issue: https://github.com/facebook/jest/issues/434
      // should to be called with `data`
      expect(save).toBeCalledWith({ type: 'jira' });
    });
  });
});
