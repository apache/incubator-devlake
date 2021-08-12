import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import ExamplePlugin from './examplePlugin';
import { ScheduleModule } from '@nestjs/schedule';
import IssueCollector from './runners/IssueCollector';
import { User } from './entities/user';
import { ConfigModule, ConfigService } from '@nestjs/config';

@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: (config: ConfigService) => ({
        type: config.get<'postgres' | 'mysql'>('DB_TYPE', 'mysql'),
        url: config.get<string>('DB_URL'),
        entityPrefix: 'plugin_example_',
        entities: [User,/*'./plugins/jira/src/entities/*.{js,ts}'*/],
        migrations: ['./plugins/jira/src/migration/*.{js,ts}'],
      }),
      inject: [ConfigService],
    }),
    ScheduleModule.forRoot(),
  ],
  providers: [ExamplePlugin, IssueCollector],
  exports: [ExamplePlugin],
})
class ExampleModule {}

export default ExampleModule;
