import { IsArray, IsNotEmpty, IsOptional } from 'class-validator';
import { PaginationRequest } from './pagination';

export class CreateSourceTask {
  @IsArray()
  collector: string[];

  @IsArray()
  enricher: string[];

  @IsNotEmpty()
  options: Record<string, unknown>;
}

export class ListSourceTask extends PaginationRequest {
  @IsOptional()
  source_id?: string;
}
