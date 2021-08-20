import { Body, Controller, Post } from '@nestjs/common';
import { SourceTask } from '../models';
import { SourceTaskService } from '../services/sourceTask';
import { CreateSourceTask } from '../types/sourceTask';

@Controller('source')
export class SourceTaskController {
  constructor(private readonly sourceTaskService: SourceTaskService) {}

  @Post(':id/task')
  async create(@Body() task: CreateSourceTask): Promise<SourceTask> {
    return;
  }
}
