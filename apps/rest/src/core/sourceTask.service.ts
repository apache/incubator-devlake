import { Injectable } from '@nestjs/common';
import { UniqueID } from './base.model';
import { SourceTask } from './sourceTask.model';
import { CreateSourceTask } from './sourceTask.type';

@Injectable()
export class SourceTaskService {
  async create(
    sourceId: UniqueID,
    sourceTask: CreateSourceTask,
  ): Promise<SourceTask> {
    return;
  }
}
