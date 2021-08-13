import { IsIn, IsNotEmpty } from 'class-validator';
import { SourceType, SupportedSourceType } from '../models/source';

export class CreateSource {
  @IsIn(SupportedSourceType)
  type: SourceType;

  @IsNotEmpty()
  options: Record<string, unknown>;
}
