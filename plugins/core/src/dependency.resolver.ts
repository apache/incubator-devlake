import { Injectable, Type } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import BaseEntity from './base.entity';
import { EXPORTS_META_KEY, PRODUCER_META_KEY } from './exports.decorator';
import { IMPORTS_META_KEY } from './imports.decorator';
import Task from './task.interface';

export class DAG {
  private _tasks = [];

  constructor(task: any[]) {
    this._tasks = task;
  }

  get length(): number {
    return this._tasks.length;
  }

  async appendTask(task: Type<Task>): Promise<void> {
    this._tasks.splice(0, 0, { name: task.name });
  }

  async pushTask(task: Type<Task>): Promise<void> {
    this._tasks.push({ name: task.name });
  }

  toPipline(): any[] {
    return this._tasks;
  }

  get(index: number): any {
    return this._tasks[index];
  }

  find(query: any): any {
    return this._tasks.find((task) => {
      for (const key of Object.keys(query)) {
        if (query[key] !== task[key]) {
          return false;
        }
      }
      return true;
    });
  }

  findIndex(query: any): number {
    return this._tasks.findIndex((task) => {
      for (const key of Object.keys(query)) {
        if (query[key] !== task[key]) {
          return false;
        }
      }
      return true;
    });
  }

  getPipline() {
    return this._tasks;
  }
}

@Injectable()
export default class DependencyResolver {
  async resolve(entity: Type<BaseEntity>): Promise<DAG> {
    //TODO: fetch TASK DAG from target Entity
    const dag = new DAG([]);
    await this.resolveEntity(entity, dag);
    return dag;
  }

  async resolveEntity(entity: Type<BaseEntity>, dag: DAG): Promise<void> {
    const ProducerType = Reflect.getMetadata(PRODUCER_META_KEY, entity);
    dag.appendTask(ProducerType);
    const importEntities = Reflect.getMetadata(IMPORTS_META_KEY, ProducerType);
    if (importEntities) {
      for (const entityClass of importEntities) {
        this.resolveEntity(entityClass, dag);
      }
    }
  }
}
