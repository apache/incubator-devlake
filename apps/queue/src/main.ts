import { NestFactory } from '@nestjs/core';
import { QueueModule } from './queue.module';

async function bootstrap() {
  const app = await NestFactory.createMicroservice(QueueModule);
  await app.listen();
}
bootstrap();
