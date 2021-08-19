import { DynamicModule } from '@nestjs/common';
import { BullQueueModule } from '../bull/queue.module';
import PluginModule from 'plugins/core/src/plugins.module';
import { ConsumerService } from './service';
import { entities, migrations, plugins } from 'plugins';
import CustomTypeOrmModule from '../../../rest/src/providers/typeorm.module';

export class ConsumerModule {
  static forRoot(queue = 'default'): DynamicModule {
    return {
      module: ConsumerModule,
      imports: [
        BullQueueModule.forRoot(queue),
        PluginModule.forRootAsync(plugins),
        CustomTypeOrmModule.forRootAsync({
          entities,
          migrations,
          migrationsRun: true,
        }),
      ],
      providers: [ConsumerService],
      exports: [ConsumerService],
    };
  }
}
