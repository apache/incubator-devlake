import { Injectable } from '@nestjs/common';
import Source from './source.model';
import { CreateSource } from './source.type';

@Injectable()
export class SourceService {
  async create(data: CreateSource): Promise<Source> {
    const source = new Source();
    source.type = data.type;
    source.options = data.options;
    // TODO: save source to db
    // ignore options
    delete source.options;
    return source;
  }
}
