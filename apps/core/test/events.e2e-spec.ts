import { Test, TestingModule } from '@nestjs/testing';
import { INestMicroservice } from '@nestjs/common';
import { EventsModule } from '../events/events.module';
import { EventsService } from '../events/events.service';

describe('EventsModule (e2e)', () => {
  let app: INestMicroservice;

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [EventsModule],
    }).compile();

    app = moduleFixture.createNestMicroservice(moduleFixture);
    await app.init();
  });

  afterEach(async () => {
    await app.close();
  });

  it('Initialized', () => {
    const service = app.get(EventsService);
    expect(service).toBeDefined();
  });

  it('Event', (done) => {
    const service = app.get(EventsService);
    const mockedFuc = jest.fn().mockImplementation((value) => {
      expect(value).toEqual({ ut: 'test' });
      //use done to make sure function be called
      done();
      return;
    });
    service.on('Custom', mockedFuc);
    service.emit('Custom', { ut: 'test' });
  });
});
