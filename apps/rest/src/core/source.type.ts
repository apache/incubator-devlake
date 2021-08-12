import { IsIn, IsNotEmpty } from 'class-validator';
import { SourceType, SupportedSourceType } from './source.model';

export class CreateSource {
  @IsIn(SupportedSourceType)
  type: SourceType;

  @IsNotEmpty()
  options: Record<string, unknown>;
}
