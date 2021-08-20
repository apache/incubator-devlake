import { Injectable, Scope } from '@nestjs/common';
import { ContextIdFactory } from '@nestjs/core';
import { Test, TestingModule } from '@nestjs/testing';
import PluginModule from 'plugins/core/src/plugins.module';
import Jira from '../src';

// Add Mocked Jira extends with Jira to Change injectable scope. so testing module could load it
@Injectable({ scope: Scope.DEFAULT })
class MockedJira extends Jira {}

describe('Jira Plugin', () => {
  let app: TestingModule;

  beforeAll(async () => {
    app = await Test.createTestingModule({
      imports: [PluginModule.forRootAsync([MockedJira])],
    }).compile();
  });

  it('Jira', async () => {
    const jira = await app.resolve<Plugin>(
      'MockedJira',
      ContextIdFactory.create(),
      {
        strict: false,
      },
    );
    expect(jira).toBeDefined();
  });
});
