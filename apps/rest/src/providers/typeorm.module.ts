import { DynamicModule } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { TypeOrmModule, TypeOrmModuleOptions } from '@nestjs/typeorm';

/**
 * 自定义typeorm模块
 * CustomTypeOrmModule
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
          type: 'mysql',
          url: config.get<string>('DB_URL'),
          ...options,
        };
      },
      inject: [ConfigService],
    });
  }
}
