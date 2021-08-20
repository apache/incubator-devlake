import { IsNotEmpty, IsOptional } from 'class-validator';
import { PaginationRequest } from './pagination';

export class CreateSource {
  /**
   * source type is defined when you register plugin
   * now it is equal to plugin's class name, eg: Jira
   */
  @IsNotEmpty()
  type: string;

  /**
   * name is used to distinct different source with a same plugin
   * although there is a global uuid, but a proper name is human friendly.
   */
  @IsOptional()
  name?: string;

  /**
   * options will be delivered as parameters to plugin when plugin running
   */
  @IsNotEmpty()
  options: Record<string, unknown>;
}

export class UpdateSource extends CreateSource {}

export class ListSource extends PaginationRequest {
  /**
   * filter source list by type
   */
  @IsOptional()
  type?: string;
}
