import { Body, Controller, Post } from '@nestjs/common';
import Source from 'apps/model/core/source';
import { SourceService } from './source.service';
import { CreateSource } from './source.type';

@Controller('source')
export class SourceController {
  constructor(private readonly sourceService: SourceService) {}

  @Post()
  async create(@Body() source: CreateSource): Promise<Source> {
    return await this.sourceService.create(source);
  }
}
