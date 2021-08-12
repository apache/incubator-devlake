import { Injectable } from '@nestjs/common';
import { ContextId, ModuleRef } from '@nestjs/core';
import { IExecutable } from './executable.interface';

@Injectable()
class CollectorRef {
  constructor(private moduleRef: ModuleRef) {}

  get(name: string, plugin: string): IExecutable<any> {
    return this.moduleRef.get(`${plugin}/collector/${name}`, { strict: false });
  }

  resolve(
    name: string,
    plugin: string,
    contextId?: ContextId,
  ): Promise<IExecutable<any>> {
    return this.moduleRef.resolve(`${plugin}/collector/${name}`, contextId, {
      strict: false,
    });
  }
}

export default CollectorRef;
