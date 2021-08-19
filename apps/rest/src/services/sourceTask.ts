import { Injectable } from '@nestjs/common';
import { UniqueID } from '../models/base';
import SourceTask from '../models/sourceTask';
import { CreateSourceTask } from '../types/sourceTask';

@Injectable()
export class SourceTaskService {
  async create(
    sourceId: UniqueID,
    sourceTask: CreateSourceTask,
  ): Promise<SourceTask> {
    return;
  }
}
