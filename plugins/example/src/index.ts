import { Module } from '@nestjs/common';
import ExamplePlugin from './examplePlugin';
import { ScheduleModule } from '@nestjs/schedule';
import IssueCollector from './runners/IssueCollector';
import CustomTypeOrmModule from '../../../apps/rest/src/customTypeOrmModule';

@Module({
  imports: [
    CustomTypeOrmModule.forRootAsync('exampleModuleDb', {
      entityPrefix: 'plugin_example_',
      entitiesFunc: () => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        return require.context('./entities', true, /\.ts/);
      },
      migrationsFunc: () => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        return require.context('./migrations', true, /\.ts/);
      },
    }),
    ScheduleModule.forRoot(),
  ],
  providers: [ExamplePlugin, IssueCollector],
  exports: [ExamplePlugin],
})
class ExampleModule {}

export default ExampleModule;
