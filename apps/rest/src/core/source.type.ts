import { SourceType, SupportedSourceType } from 'apps/model/core/source';
import { IsIn, IsNotEmpty } from 'class-validator';

export class CreateSource {
  @IsIn(SupportedSourceType)
  type: SourceType;

  @IsNotEmpty()
  options: Record<string, unknown>;
}
