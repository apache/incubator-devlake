import { UniqueID } from './base.model';

export const SupportedSourceType = ['jira', 'gitlab'] as const;

export type SourceType = typeof SupportedSourceType[number];

export default class Source {
  id: UniqueID;
  type: SourceType;

  options: Record<string, unknown>;
}
