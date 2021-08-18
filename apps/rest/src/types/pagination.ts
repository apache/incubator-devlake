import { Type } from 'class-transformer';
import { Max, Min } from 'class-validator';

export class PaginationRequest {
  @Min(1)
  @Type(() => Number)
  page: number;

  @Max(100)
  @Type(() => Number)
  pagesize: number;
}

export class PaginationResponse<T> {
  page: number;
  offset: number;
  data: T[];
}
