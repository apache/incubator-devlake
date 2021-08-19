import { DynamicModule } from '@nestjs/common';
import { TypeOrmModule, TypeOrmModuleOptions } from '@nestjs/typeorm';
import { ConfigModule, ConfigService } from '@nestjs/config';

/**
 * 自定义typeorm模块
 * 引入了config模块自动获取url，同时封装了以目录获取entities、migrations
 * CustomTypeOrmModule
 * import configModule to get db url
 * also add hack to entities/migrations path auto
 */
export default class CustomTypeOrmModule {
  static forRootAsync(options?: TypeOrmModuleOptions): DynamicModule {
    return TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: async (config: ConfigService) => {
        return <TypeOrmModuleOptions>{
          type: 'mysql',
          url: config.get<string>('DB_URL'),
          ...options,
        };
      },
      inject: [ConfigService],
    });
  }
}
