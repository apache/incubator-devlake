import { Test, TestingModule } from '@nestjs/testing';
import { INestMicroservice } from '@nestjs/common';
import { QueueModule } from './../src/queue.module';
import { ConsumerService } from '../src/consumer/service';
import { ProducerService } from '../src/producer/service';

jest.mock('ioredis', () => {
  return {
    default: require('ioredis-mock/jest'),
  };
});

describe('QueueController (e2e)', () => {
  let app: INestMicroservice;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [QueueModule],
    }).compile();

    app = moduleFixture.createNestMicroservice({});
    await app.init();
  });

  afterAll(async () => {
    await app.close();
  });

  it('Consumer', () => {
    const consumer = app.get(ConsumerService);
    expect(consumer).toBeDefined();
  });

  it('Producer', () => {
    const producer = app.get(ProducerService);
    expect(producer).toBeDefined();
  });
});
