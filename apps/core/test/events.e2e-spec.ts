import { Test, TestingModule } from '@nestjs/testing';
import { INestMicroservice } from '@nestjs/common';
import { EventsModule } from '../events/events.module';
import { EventsService } from '../events/events.service';

describe('EventsModule (e2e)', () => {
  let app: INestMicroservice;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [EventsModule],
    }).compile();

    app = moduleFixture.createNestMicroservice(moduleFixture);
    await app.init();
  });

  afterAll(async () => {
    await app.close();
  });

  it('Initialized', () => {
    const service = app.get(EventsService);
    expect(service).toBeDefined();
  });

  // it('Event', (done) => {
  //   const service = app.get(EventsService);
  //   const linstenFuc = (value) => {
  //     expect(value).toEqual({ ut: 'test' });
  //     //use done to make sure function be called
  //     done();
  //     return;
  //   };
  //   service.on('Custom', linstenFuc);
  //   service.emit('Custom', { ut: 'test' });
  // }, 10000);
});
