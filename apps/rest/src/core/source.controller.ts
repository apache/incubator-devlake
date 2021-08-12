import { Body, Controller, Post } from '@nestjs/common';
import Source from './source.model';
import { CreateSource } from './source.type';

@Controller('source')
export class SourceController {
  @Post()
  async create(@Body() source: CreateSource): Promise<Source> {
    return;
  }
}
