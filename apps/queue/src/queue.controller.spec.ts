import { Test, TestingModule } from '@nestjs/testing';
import { QueueController } from './queue.controller';
import { QueueService } from './queue.service';

describe('QueueController', () => {
  let queueController: QueueController;

  beforeEach(async () => {
    const app: TestingModule = await Test.createTestingModule({
      controllers: [QueueController],
      providers: [QueueService],
    }).compile();

    queueController = app.get<QueueController>(QueueController);
  });

  describe('root', () => {
    it('should return "Hello World!"', () => {
      expect(queueController.getHello()).toBe('Hello World!');
    });
  });
});
