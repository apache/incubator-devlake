import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import Jira from 'plugins/jira/src';
import { ConsumerModule } from './consumer';

@Module({
  imports: [ConfigModule.forRoot({ isGlobal: true }), ConsumerModule.forRoot()],
  providers: [{ provide: 'Jira', useClass: Jira }],
})
export class QueueModule {}
