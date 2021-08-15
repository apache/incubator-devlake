import { DynamicModule } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { TypeOrmModule, TypeOrmModuleOptions } from '@nestjs/typeorm';

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
    options?: TypeOrmModuleOptions,
  ): DynamicModule {
    return TypeOrmModule.forRootAsync({
      name,
      imports: [ConfigModule],
      useFactory: async (config: ConfigService) => {
        return <TypeOrmModuleOptions>{
          type: config.get<string>('DB_TYPE'),
          url: config.get<string>('DB_URL'),
          ...options,
        };
      },
      inject: [ConfigService],
    });
  }
}
