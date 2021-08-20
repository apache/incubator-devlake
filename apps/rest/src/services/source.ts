import { Injectable } from '@nestjs/common';
import { InjectEntityManager } from '@nestjs/typeorm';
import { EntityManager, FindManyOptions } from 'typeorm';
import { UniqueID } from '../models/base';
import Source from '../models/source';
import { PaginationResponse } from '../types/pagination';
import { CreateSource, ListSource, UpdateSource } from '../types/source';

@Injectable()
export class SourceService {
  constructor(@InjectEntityManager() private em: EntityManager) {}

  async list(filter: ListSource): Promise<PaginationResponse<Source>> {
    const offset = filter.pagesize * (filter.page - 1);
    const options: FindManyOptions<Source> = {
      skip: offset,
      take: filter.pagesize,
    };
    if (filter.type) {
      options.where = {
        type: filter.type,
      };
    }
    const sources = await this.em.getRepository(Source).find(options);
    return {
      offset,
      page: filter.page,
      data: sources,
    };
  }

  async create(data: CreateSource): Promise<Source> {
    const source = new Source();
    source.type = data.type;
    source.options = data.options;
    source.name = data.name;
    await this.em.save(source);
    return source;
  }

  async get(id: UniqueID): Promise<Source> {
    return await this.em.getRepository(Source).findOneOrFail(id);
  }

  async delete(id: UniqueID): Promise<Source> {
    const target = await this.get(id);
    await this.em.remove(target);
    return target;
  }

  async update(id: UniqueID, data: UpdateSource): Promise<Source> {
    const target = await this.get(id);
    target.type = data.type;
    target.options = data.options;
    target.name = data.name;
    await this.em.save(target);
    return target;
  }
}
