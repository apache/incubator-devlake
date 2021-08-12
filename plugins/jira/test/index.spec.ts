import { Injectable, Scope } from '@nestjs/common';
import { ContextIdFactory } from '@nestjs/core';
import { Test, TestingModule } from '@nestjs/testing';
import PluginModule from 'plugins/core/src/plugin.module';
import Jira from '../src';

//exnted new plugin class to change scope to default.so test module could resolve the instance
@Injectable({ scope: Scope.DEFAULT })
class MockedJira extends Jira {}

describe('Jira Plugin', () => {
  let testModule: TestingModule;
  beforeEach(async () => {
    testModule = await Test.createTestingModule({
      providers: PluginModule.Register([
        {
          name: 'Jira',
          schedule: MockedJira,
        },
      ]).providers,
    }).compile();
  });

  it('Jira', async () => {
    const jira = await testModule.resolve('Jira', ContextIdFactory.create(), {
      strict: false,
    });
    expect(jira).toBeDefined();
    expect(jira.execute).toBeDefined();
  });
});
