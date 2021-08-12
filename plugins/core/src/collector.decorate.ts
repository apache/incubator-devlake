import { Type } from '@nestjs/common';
import { IExecutable } from './executable.interface';

export type CollectorMap = {
  [key: string]: Type<IExecutable<any>>;
};

export const COLLECTORS_METADATA = 'COLLECTORS';

export default function Collector(collectors: CollectorMap) {
  return (target) => {
    Reflect.defineMetadata(COLLECTORS_METADATA, collectors, target);
  };
}
