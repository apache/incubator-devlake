import Collector from './collector';

type PrimaryValues = Record<string, string | number | boolean>;

interface Enricher extends Collector {
  /**
   * 返回改依赖项，enricher可以依赖collector或者enricher，但不能循环依赖
   * return dependencies
   * dependencies can include collector or enricher
   * 返回格式(return format):
   * {
   *   # depend on collector/enricher in this plugin
   *   [collectorOrEnricherName: string]: primaryKeyObject
   *   # or depend on collector/enricher in other plugin
   *   'pluginName' + '/' + 'collectorOrEnricherName': primaryKeyObject
   * }
   */
  dependencies(primaryKeys: PrimaryValues): Record<string, PrimaryValues>;

  /**
   * 是否支持懒加载，即查询数据时才进行计算，如果支持，那么查询数据必须通过queryData
   * return if support lazy load for query data,
   * You must query data by queryData() when support lazy load
   * 默认为不支持 default: false
   */
  couldLazyLoad?(): boolean;

  /**
   * 使用队列开始计算数据，此时不应该请求任何外部数据
   * start collect data within queue
   * Must keep enricher offline when cal data
   * @param primaryKeys
   * @param consumerModule
   */
  calData(primaryKeys: PrimaryValues, consumerModule): Promise<void>;

  /**
   * 查询指定主键对应的数据
   * query data for these primaryKeys
   * @param primaryKeys
   */
  queryData(primaryKeys: PrimaryValues): Promise<any>;

  /**
   * 自检指定主键的数据是否正确
   * self check data for these primaryKeys
   * @param primaryKeys
   * @throws Error error of self check failure reason 自检失败原因
   */
  selfCheck(primaryKeys: PrimaryValues): Promise<boolean>;
}

export default Enricher;
