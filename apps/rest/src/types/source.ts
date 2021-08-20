import { IsNotEmpty, IsOptional } from 'class-validator';
import { PaginationRequest } from './pagination';

export class CreateSource {
  @IsNotEmpty()
  type: string;

  @IsOptional()
  name?: string;

  @IsNotEmpty()
  options: Record<string, unknown>;
}

export class UpdateSource extends CreateSource {}

export class ListSource extends PaginationRequest {
  @IsOptional()
  type?: string;
}
