import BaseEntity from './base.entity';

export const IMPORTS_META_KEY = 'IMPORTS_ENTITIES';

export default function Imports(entities: typeof BaseEntity[]) {
  return (target) => {
    Reflect.defineMetadata(IMPORTS_META_KEY, entities, target);
  };
}
