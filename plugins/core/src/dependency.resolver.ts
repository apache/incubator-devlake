import { Injectable } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import BaseEntity from './base.entity';

export class DAG {}//TEMP TYPE

@Injectable()
export default class DependencyResolver {
  constructor(private moduleRef: ModuleRef) {}

  async resolve(entity: typeof BaseEntity): Promise<DAG> {
    //TODO: fetch TASK DAG from target Entity
    return {};
  }
}
