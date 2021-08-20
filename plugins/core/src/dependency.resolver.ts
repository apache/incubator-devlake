import { Injectable, Type } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import BaseEntity from './base.entity';
import { EXPORTS_META_KEY, PRODUCER_META_KEY } from './exports.decorator';
import Task from './task.interface';

export class DAG {
  private _tasks = [];
  async appendTask(task: Type<Task>): Promise<void> {
    this._tasks.splice(0, 0, { name: task.name });
  }

  async pushTask(task: Type<Task>): Promise<void> {
    this._tasks.push({ name: task.name });
  }

  toPipline(): any[] {
    return this._tasks;
  }
}

@Injectable()
export default class DependencyResolver {
  constructor(private moduleRef: ModuleRef) {}

  async resolve(entity: Type<BaseEntity>): Promise<DAG> {
    //TODO: fetch TASK DAG from target Entity
    const dag = new DAG();
    await this.resolveEntity(entity, dag);
    return dag;
  }

  async resolveEntity(entity: Type<BaseEntity>, dag: DAG): Promise<void> {
    const ProducerType = Reflect.getMetadata(PRODUCER_META_KEY, entity);
    dag.appendTask(ProducerType);
    const importEntities = Reflect.getMetadata(EXPORTS_META_KEY, ProducerType);
    if (importEntities) {
      for (const entityClass of importEntities) {
        this.resolveEntity(entityClass, dag);
      }
    }
  }
}
