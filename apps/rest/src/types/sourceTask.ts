import { IsArray, IsNotEmpty } from 'class-validator';

export class CreateSourceTask {
  @IsArray()
  collector: string[];

  @IsArray()
  enricher: string[];

  @IsNotEmpty()
  options: Record<string, unknown>;
}
