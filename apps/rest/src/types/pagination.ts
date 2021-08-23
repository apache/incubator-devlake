import { Type } from 'class-transformer';
import { Max, Min } from 'class-validator';

export class PaginationRequest {
  @Min(1)
  @Type(() => Number)
  page = 1;

  @Max(100)
  @Type(() => Number)
  pagesize = 20;
}

export class PaginationResponse<T> {
  total: number;
  offset: number;
  page: number;
  pagesize: number;
  data: T[];
}
