import { Injectable } from '@nestjs/common';
import { UniqueID } from '../models/base';
import SourceTask from '../models/sourceTask';
import { CreateSourceTask } from '../types/sourceTask';
import { InjectEntityManager } from '@nestjs/typeorm';
import { EntityManager, FindConditions, FindManyOptions } from 'typeorm';

@Injectable()
export class SourceTaskService {
  constructor(@InjectEntityManager() private em: EntityManager) {}

  // async list(filter: ListSource): Promise<PaginationResponse<Source>> {
  //   const offset = filter.pagesize * (filter.page - 1);
  //   const where: FindConditions<Source> = {};
  //   const options: FindManyOptions<Source> = {
  //     skip: offset,
  //     take: filter.pagesize,
  //   };
  //   if (filter.type) {
  //     where.type = filter.type;
  //   }
  //   options.where = where;

  //   const total = await this.em.getRepository(Source).count(where);
  //   const sources = await this.em.getRepository(Source).find(options);
  //   return {
  //     offset,
  //     total,
  //     data: sources,
  //   };
  // }

  async create(
    sourceId: UniqueID,
    data: CreateSourceTask,
  ): Promise<SourceTask> {
    const sourceTask = new SourceTask();
    sourceTask.source_id = sourceId;
    sourceTask.collector = data.collector;
    sourceTask.enricher = data.enricher;
    sourceTask.options = data.options;
    await this.em.save(sourceTask);
    return sourceTask;
  }

  // async get(id: UniqueID): Promise<Source> {
  //   return await this.em.getRepository(Source).findOneOrFail(id);
  // }

  // async delete(id: UniqueID): Promise<Source> {
  //   const target = await this.get(id);
  //   await this.em.remove(target);
  //   return target;
  // }
}
