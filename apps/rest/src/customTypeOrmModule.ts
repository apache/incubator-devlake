import { DynamicModule } from '@nestjs/common';
import { TypeOrmModule, TypeOrmModuleOptions } from '@nestjs/typeorm';
import { ConfigModule, ConfigService } from '@nestjs/config';

export declare type CustomTypeOrmModuleOptions = {
  entitiesFunc?: () => any | { keys: () => string[] };
  migrationsFunc?: () => any | { keys: () => string[] };
} & TypeOrmModuleOptions;

/**
 * 自定义typeorm模块
 * 引入了config模块自动获取url，同时封装了以目录获取entities、migrations
 * CustomTypeOrmModule
 * import configModule to get db url
 * also add hack to entities/migrations path auto
 */
export default class CustomTypeOrmModule {
  static forRootAsync(
    name?: string,
    options?: CustomTypeOrmModuleOptions,
  ): DynamicModule {
    return TypeOrmModule.forRootAsync({
      name,
      imports: [ConfigModule],
      useFactory: async (config: ConfigService) => {
        // call func to get all entities/migrations auto
        let entities = [];
        if (options.entitiesFunc) {
          const r = options.entitiesFunc();
          entities = r.keys().map((key) => r(key).default);
          delete options.entitiesFunc;
        }
        let migrations = [];
        if (options.migrationsFunc) {
          const r = options.migrationsFunc();
          migrations = r.keys().map((key) => r(key).default);
          delete options.migrationsFunc;
        }

        return <TypeOrmModuleOptions>{
          type: config.get<'postgres' | 'mysql'>('DB_TYPE', 'mysql'),
          url: config.get<string>('DB_URL'),
          migrationsRun:
            !name && config.get<boolean>('DB_AUTO_APP_MIGRATE', true),
          entities: entities,
          migrations: migrations,
          ...options,
        };
      },
      inject: [ConfigService],
    });
  }
}
