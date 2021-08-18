import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, INestMicroservice } from '@nestjs/common';
import * as request from 'supertest';
import { QueueModule } from './../src/queue.module';
import { ConsumerService } from '../src/consumer/service';
import { ProducerService } from '../src/producer/service';

describe('QueueController (e2e)', () => {
  let app: INestMicroservice;

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [QueueModule],
    }).compile();

    app = moduleFixture.createNestMicroservice({});
    await app.init();
  });

  afterEach(async () => {
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
