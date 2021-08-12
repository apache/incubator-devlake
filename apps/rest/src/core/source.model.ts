export const SupportedSourceType = ['jira', 'gitlab'] as const;

export type SourceType = typeof SupportedSourceType[number];

export default class Source {
  id: number;
  type: SourceType;

  options: Record<string, unknown>;
}
