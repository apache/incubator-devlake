import BaseEntity from './base.entity';

export const EXPORTS_META_KEY = 'EXPORTS_ENTITIES';
export const PRODUCER_META_KEY = 'PRODUCER_TASK';

export default function Exports(entity: typeof BaseEntity) {
  return (target) => {
    Reflect.defineMetadata(EXPORTS_META_KEY, entity, target);
    Reflect.defineMetadata(PRODUCER_META_KEY, target, entity);
  };
}
