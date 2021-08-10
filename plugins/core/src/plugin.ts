interface Plugin {
  /**
   * 当前插件的名字
   * name of plugin
   */
  name(): string;

  /**
   * 当前插件的版本号，从1开始的正整数
   * current version of plugin, version ∈ {1, 2, 3, ...}
   */
  version(): number;

  /**
   * 注册插件，其中应该包括各种组件的定义
   * register plugin
   * It should include collector define/enricher define/or more
   */
  register(): Promise<void>;

  /**
   * 注销插件，其中可以包含各种服务的停止
   * unregister plugin
   * It can iniclude collector stop/enricher stop/or more
   */
  unregister(): Promise<void>;

  /**
   * 初始化所需的环境，其中应该包括数据表的初始化、文件目录的创建等
   * migrateUp should include db migrate / path create and so on
   * @param pluginPrev db table or path prev string 数据表或者目录应该保持的前缀
   * @param oldVersion 旧版本号
   * @return newVersion 新版本号
   */
  migrateUp(pluginPrev: string, oldVersion: string): Promise<string>;

  /**
   * 销毁所需的环境，其中应该包括数据表、文件目录的删除等
   * migrateDown should include db table drop / path remove and so on
   * @param currentVersion 当前版本号
   */
  migrateDown(currentVersion: string): Promise<void>;
}

export default Plugin;
