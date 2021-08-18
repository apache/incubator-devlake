import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  Query,
} from '@nestjs/common';
import { UniqueID } from '../models/base';
import Source from '../models/source';
import { PaginationResponse } from '../types/pagination';
import { CreateSource, ListSource, UpdateSource } from '../types/source';

@Controller('source')
export class SourceController {
  @Post()
  async create(@Body() source: CreateSource): Promise<Source> {
    return;
  }

  @Get()
  async list(@Query() filter: ListSource): Promise<PaginationResponse<Source>> {
    // FIXME: filter.page and filter.pagesize is of type string
    console.log(filter);
    return;
  }

  @Get(':id')
  async get(@Param('id') id: UniqueID): Promise<Source> {
    return;
  }

  @Put(':id')
  async update(
    @Param('id') id: UniqueID,
    @Body() source: UpdateSource,
  ): Promise<Source> {
    return;
  }

  @Delete(':id')
  async delete(@Param('id') id: UniqueID): Promise<Source> {
    return;
  }
}
