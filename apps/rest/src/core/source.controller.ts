import { Body, Controller, Delete, Get, Post, Put } from '@nestjs/common';
import Source from './source.model';
import { SourceService } from './source.service';
import { CreateSource } from './source.type';

@Controller('source')
export class SourceController {
  constructor(private readonly sourceService: SourceService) {}

  @Post()
  async create(@Body() source: CreateSource): Promise<Source> {
    return await this.sourceService.create(source);
  }

  @Get()
  async list(): Promise<Source[]> {
    return [];
  }

  @Get()
  async get(): Promise<Source> {
    return;
  }

  @Put()
  async update(): Promise<Source> {
    return;
  }

  @Delete()
  async delete(): Promise<Source> {
    return;
  }
}
