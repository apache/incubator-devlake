import { Injectable } from '@nestjs/common';
import { InjectEntityManager } from '@nestjs/typeorm';
import { EntityManager } from 'typeorm';
import Source from '../models/source';
import { CreateSource } from '../types/source';

@Injectable()
export class SourceService {
  constructor(@InjectEntityManager() private em: EntityManager) {}

  async create(data: CreateSource): Promise<Source> {
    const source = new Source();
    source.type = data.type;
    source.options = data.options;
    await this.em.save(source);
    // ignore options
    delete source.options;
    return source;
  }
}
