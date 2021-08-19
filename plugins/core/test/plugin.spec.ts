import { Test, TestingModule } from '@nestjs/testing';
import { DAG } from '../src/dependency.resolver';
import Plugin from '../src/plugin.interface';
import PluginModule from '../src/plugins.module';


class TestPlugin implements Plugin {
  async execute(...args: any[]): Promise<DAG> {
    return {};
  }
}

describe('AppController', () => {
  let app: TestingModule;

  beforeAll(async () => {
    app = await Test.createTestingModule({
      imports: [PluginModule.forRootAsync([TestPlugin])],
    }).compile();
  });

  it('Get Plugin By Name', async () => {
    const test = app.get('TestPlugin');
    expect(test).toBeDefined();
  });
});
