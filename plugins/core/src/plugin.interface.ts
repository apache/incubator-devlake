import { Type } from '@nestjs/common';
import BaseEntity from './base.entity';
import { DAG } from './dependency.resolver';
import IExecutable from './executable.interface';

/**
 * Interface for Plugin
 */
interface Plugin extends IExecutable<DAG> {
  /**
   * 运行插件，使用指定参数自行创建collector/enricher
   * start plugin
   * It should include collector init/enricher init/or more by args
   */
  execute(...args: any[]): Promise<DAG>;

  exports(): Type<BaseEntity>[];
}

export default Plugin;
