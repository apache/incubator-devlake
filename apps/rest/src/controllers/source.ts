import {
  Body,
  Controller,
  Delete,
  Get,
  Logger,
  Param,
  Post,
  Put,
  Query,
} from '@nestjs/common';
import { UniqueID } from '../models/base';
import Source from '../models/source';
import { SourceService } from '../services/source';
import { PaginationResponse } from '../types/pagination';
import { CreateSource, ListSource, UpdateSource } from '../types/source';

@Controller('source')
export class SourceController {
  private readonly logger = new Logger(SourceController.name);

  constructor(private readonly sourceService: SourceService) {}

  @Post()
  async create(@Body() source: CreateSource): Promise<Source> {
    // TODO: validator source type
    // should write registered plugins into database
    // then reject create request if target source type not exist
    return await this.sourceService.create(source);
  }

  @Get()
  async list(@Query() filter: ListSource): Promise<PaginationResponse<Source>> {
    return await this.sourceService.list(filter);
  }

  @Get(':id')
  async get(@Param('id') id: UniqueID): Promise<Source> {
    return await this.sourceService.get(id);
  }

  @Put(':id')
  async update(
    @Param('id') id: UniqueID,
    @Body() source: UpdateSource,
  ): Promise<Source> {
    return await this.sourceService.update(id, source);
  }

  @Delete(':id')
  async delete(@Param('id') id: UniqueID): Promise<Source> {
    return await this.sourceService.delete(id);
  }
}
