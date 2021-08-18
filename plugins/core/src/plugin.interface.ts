import IExecutable from './executable.interface';

/**
 * Interface for Plugin
 */
interface Plugin extends IExecutable {
  /**
   * 当前插件的名字
   * name of plugin
   */
  name(): string;

  /**
   * 运行插件，使用指定参数自行创建collector/enricher
   * start plugin
   * It should include collector init/enricher init/or more by args
   */
  execute(...args: any[]): Promise<any>;
}

export default Plugin;
