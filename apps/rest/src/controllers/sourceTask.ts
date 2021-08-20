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
import { SourceTaskService } from '../services/sourceTask';
import { SourceTask } from '../models';
import { UniqueID } from '../models/base';
import { CreateSourceTask, ListSourceTask } from '../types/sourceTask';
import { PaginationResponse } from '../types/pagination';



@Controller()
export class SourceTaskController {
  constructor(private readonly sourceTaskService: SourceTaskService) {}

  @Post('source/:id/task')
  async create(
    @Param('id') sourceId: string,
    @Body() task: CreateSourceTask,
  ): Promise<SourceTask> {
    return await this.sourceTaskService.create(sourceId, task);
  }

  @Get('/task')
  async list(@Query() filter: ListSourceTask): Promise<PaginationResponse<SourceTask>> {
    return await this.sourceTaskService.list(filter);
  }

  // @Get('source/:id/task')
  // async list(@Query() : ListTask): Promise<PaginationResponse<SourceTask>> {
  //   return await this.sourceTaskService.list(filter);
  // }

  // @Get('task/:id')
  // async get(@Param('id') id: UniqueID): Promise<Source> {
  //   return await this.sourceService.get(id);
  // }

  // @Delete('task/:id')
  // async delete(@Param('id') id: UniqueID): Promise<Source> {
  //   return await this.sourceService.delete(id);
  // }
}
