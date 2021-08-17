import { IsIn, IsNotEmpty, IsOptional } from 'class-validator';
import { SourceType, SupportedSourceType } from '../models/source';
import { PaginationRequest } from './pagination';

export class CreateSource {
  @IsIn(SupportedSourceType)
  type: SourceType;

  @IsNotEmpty()
  options: Record<string, unknown>;
}

export class UpdateSource extends CreateSource {}

export class ListSource extends PaginationRequest {
  @IsIn(SupportedSourceType)
  @IsOptional()
  type?: SourceType;
}
