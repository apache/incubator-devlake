import { UniqueID } from './base.model';

export class SourceTask {
  id: UniqueID;
  source_id: number;
  collector: string[];
  enricher: string[];
  options: Record<string, unknown>;
}
