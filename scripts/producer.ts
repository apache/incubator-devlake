import { NestFactory } from '@nestjs/core';
import { Module } from '@nestjs/common';
import { ProducerModule } from 'apps/queue/src/producer';
import { ProducerService } from 'apps/queue/src/producer/service';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [ConfigModule.forRoot({ isGlobal: true }), ProducerModule.forRoot()],
})
class CustomModule {}

async function bootstrap() {
  const appContext = await NestFactory.createApplicationContext(CustomModule);
  const producer = await appContext.get(ProducerService, { strict: false });
  await producer.addJob('Jira', {});
  await appContext.close();
}
bootstrap();
