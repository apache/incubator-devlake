import { Test, TestingModule } from '@nestjs/testing';
import { HttpStatus, INestApplication } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from '../src/app.module';
import { ignoreEntityProps, truncateTableForTest } from './utils';
import * as uuid from 'uuid';
import { PaginationResponse } from '../src/types/pagination';
import Source from '../src/models/source';
import * as _ from 'lodash';

describe('SourceController (e2e)', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    await app.init();
  });

  afterAll(async () => {
    await app.close();
  });

  beforeEach(async () => {
    await truncateTableForTest(['source']);
  });

  describe('/source (POST)', () => {
    it('should return source type', () => {
      const newSource = {
        type: 'jira',
        options: {
          host: 'https://www.atlassian.com/',
          email: 'xx@example.com',
          auth: 'base64EncodedAuthToken',
        },
      };
      return request(app.getHttpServer())
        .post('/source')
        .send(newSource)
        .expect(201)
        .expect((res) => {
          expect(ignoreEntityProps(res.body)).toMatchSnapshot();
          expect(uuid.validate(res.body.id)).toBe(true);
        });
    });

    it('should return validate error', () => {
      const sourceWithoutType = {
        options: {},
      };
      return request(app.getHttpServer())
        .post('/source')
        .send(sourceWithoutType)
        .expect(400)
        .expect((res) => {
          expect(res.body).toMatchSnapshot();
        });
    });

    it('should return validate error 2', () => {
      const sourceWithoutType = {
        type: 'github',
        options: {},
      };
      return request(app.getHttpServer())
        .post('/source')
        .send(sourceWithoutType)
        .expect(400)
        .expect((res) => {
          expect(res.body).toMatchSnapshot();
        });
    });

    it('should return validate error 3', () => {
      const sourceWithoutType = {
        type: 'gitlab',
      };
      return request(app.getHttpServer())
        .post('/source')
        .send(sourceWithoutType)
        .expect(400)
        .expect((res) => {
          expect(res.body).toMatchSnapshot();
        });
    });
  });

  describe('/source/:id (GET)', () => {
    it('should return not found if record not exists', async () => {
      return request(app.getHttpServer())
        .get('/source/id-not-exist')
        .expect(404)
        .expect((res) => {
          expect(res.body).toMatchSnapshot();
        });
    });
  });

  describe('/source (CRUD)', () => {
    it('mock a CRUD progress', async () => {
      const server = request(app.getHttpServer());
      // create a jira source
      const jiraSource = {
        type: 'jira',
        options: {
          host: 'https://www.atlassian.com/',
          username: 'lake@merico.com',
          token: 'guess what?',
        },
      };
      let response = await server.post('/source').send(jiraSource);
      expect(response.status).toBe(HttpStatus.CREATED);

      const jira = response.body;
      expect(ignoreEntityProps(jira)).toMatchSnapshot();

      //  create a gitlab source
      const gitlabSource = {
        type: 'gitlab',
        options: {
          host: 'https://www.gitlab.com',
          username: 'lake@merico.com',
          token: 'gitlab token',
        },
      };

      response = await server.post('/source').send(gitlabSource);
      expect(response.status).toBe(HttpStatus.CREATED);

      const gitlab = response.body;
      expect(ignoreEntityProps(gitlab)).toMatchSnapshot();

      // list sources
      response = await server.get('/source').query({
        page: 1,
        pagesize: 10,
      });

      expect(response.status).toBe(HttpStatus.OK);
      let data: PaginationResponse<Source> = response.body;
      expect(_.omit(data, 'data')).toMatchSnapshot();
      expect(
        data.data
          .map((item) => ignoreEntityProps(item))
          .sort((a, b) => (a.type < b.type ? 1 : -1)),
      ).toMatchSnapshot();

      // filter sources
      response = await server.get('/source').query({
        page: 1,
        pagesize: 10,
        type: 'jira',
      });
      data = response.body;
      expect(_.omit(data, 'data')).toMatchSnapshot();
      expect(
        data.data.map((item) => ignoreEntityProps(item)),
      ).toMatchSnapshot();

      // update source
      await server.put(`/source/${gitlab.id}`).send({
        options: {
          token: 'updated token',
        },
      });
      // get source
      response = await server.get(`/source/${gitlab.id}`);
      expect(ignoreEntityProps(response.body)).toMatchSnapshot();

      // delete source
      await server.delete(`/source/${gitlab.id}`);

      response = await server.get(`/source/${gitlab.id}`);
      expect(response.status).toEqual(HttpStatus.NOT_FOUND);
    });
  });
});
