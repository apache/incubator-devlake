import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import Jira from 'plugins/jira/src';
import { ConsumerModule } from './consumer';
import { ProducerModule } from './producer';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    ConsumerModule.forRoot(),
    ProducerModule.forRoot(),
  ],
  providers: [{ provide: 'Jira', useClass: Jira }],
})
export class QueueModule {}
