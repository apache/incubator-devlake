import { Injectable } from '@nestjs/common';
import Collector, { PrimaryValues } from '../../../core/src/collector';

@Injectable()
class IssueCollector implements Collector {
  dependencies(primaryKeys: PrimaryValues): Record<string, PrimaryValues> {
    return {};
  }

  async cleanData(primaryKeys?: PrimaryValues): Promise<boolean> {
    return false;
  }

  async collectData(primaryKeys: PrimaryValues, consumerModule): Promise<void> {
    console.log('collectData');
    return;
  }

  async isDataPrepared(primaryKeys: PrimaryValues): Promise<boolean> {
    return false;
  }
}

export default IssueCollector;
