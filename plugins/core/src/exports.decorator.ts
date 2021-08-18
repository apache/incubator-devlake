import BaseEntity from './base.entity';

export const EXPORTS_META_KEY = 'EXPORTS_ENTITIES';

export default function Exports(entity: typeof BaseEntity) {
  return (target) => {
    Reflect.defineMetadata(EXPORTS_META_KEY, entity, target);
  };
}
