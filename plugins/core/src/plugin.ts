interface Plugin {
  /**
   * 当前插件的名字，如果未继承默认为去掉后缀的类名
   * name of plugin
   * If you do not inherit name(),  it defaults to the class name removing the suffix
   */
  name?(): string;

  /**
   * 当前插件的版本号，从1开始的正整数
   * current version of plugin, version ∈ {1, 2, 3, ...}
   */
  version(): number;

  /**
   * 运行插件，使用指定参数自行创建collector/enricher
   * start plugin
   * It should include collector init/enricher init/or more by args
   */
  execute(...args: any[]): Promise<void>;

  /**
   * 初始化所需的环境，其中应该包括数据表的初始化、文件目录的创建等
   * migrateUp should include db migrate / path create and so on
   * @param pluginPrev db table or path prev string 数据表或者目录应该保持的前缀
   */
  migrateUp(pluginPrev: string): Promise<void>;

  /**
   * 销毁所需的环境，其中应该包括数据表、文件目录的删除等
   * migrateDown should include db table drop / path remove and so on
   * @param currentVersion 当前版本号
   */
  migrateDown(currentVersion: string): Promise<void>;
}

export default Plugin;
