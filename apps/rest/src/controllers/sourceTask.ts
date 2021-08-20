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
import { CreateSourceTask } from '../types/sourceTask';
import { SourceTask } from '../models';

@Controller('source')
export class SourceTaskController {
  constructor(private readonly sourceTaskService: SourceTaskService) {}

  @Post(':id/task')
  async create(
    @Param('id') sourceId: string,
    @Body() task: CreateSourceTask,
  ): Promise<SourceTask> {
    return await this.sourceTaskService.create(sourceId, task);
  }
}
