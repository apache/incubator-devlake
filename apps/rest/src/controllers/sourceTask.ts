import {
  Body,
  Controller,
  HttpException,
  HttpStatus,
  Param,
  Post,
} from '@nestjs/common';
import { UniqueID } from '../models/base';
import { SourceTask } from '../models/sourceTask';
import { SourceTaskService } from '../services/sourceTask';
import { CreateSourceTask } from '../types/sourceTask';

@Controller('source')
export class SourceTaskController {
  constructor(private readonly sourceTaskService: SourceTaskService) {}

  @Post(':id')
  async create(
    @Param('id') id: UniqueID,
    @Body() sourceTask: CreateSourceTask,
  ): Promise<SourceTask> {
    if (sourceTask.collector.length === 0 && sourceTask.enricher.length === 0) {
      throw new HttpException(
        `Need at least one collector or enricher`,
        HttpStatus.BAD_REQUEST,
      );
    }
    return await this.sourceTaskService.create(id, sourceTask);
  }
}
