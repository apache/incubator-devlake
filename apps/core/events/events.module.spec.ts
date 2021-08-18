import { Test, TestingModule } from '@nestjs/testing';
import { EventsModule } from './events.module';
import { EventsService } from './events.service';

jest.mock('ioredis', () => {
  return {
    default: require('ioredis-mock/jest'),
  };
});

describe('EventModule', () => {
  let app: TestingModule;

  beforeAll(async () => {
    app = await Test.createTestingModule({
      imports: [EventsModule],
    }).compile();
  });

  afterAll(async () => {
    await app.close();
  });

  describe('EventServices', () => {
    it('EventServices should initlized', () => {
      const eventsService = app.get<EventsService>(EventsService);
      expect(eventsService).toBeDefined();
    });
    it('Event', (done) => {
      const sub = app.get('REDIS_SUB_CLIENT');
      const subSyncClient = sub.createConnectedClient();
      const pub = app.get('REDIS_PUB_CLIENT');
      jest.spyOn(pub, 'publish').mockImplementation((env, message) => {
        return subSyncClient.publish(env, message);
      });
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
});
