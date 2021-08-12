import { Type } from '@nestjs/common';
import { IExecutable } from './executable.interface';

export type EnricherMap = {
  [key: string]: Type<IExecutable<any>>;
};

export const ENRICHERS_METADATA = 'ENRICHERS';

export default function Enricher(enrichers: EnricherMap) {
  return (target) => {
    Reflect.defineMetadata(ENRICHERS_METADATA, enrichers, target);
  };
}
