import PluginInterface from './core/src/plugin.interface';
import { Type } from '@nestjs/common';

/**
 * 导入plugins目录下的所有plugin/migration脚本/entity
 * import and export all plugins/migrations/entities
 */
export const pluginRecords: Record<string, Type<PluginInterface>> = {};
export const plugins: Type<PluginInterface>[] = [];
export const migrations = [];
export const entities = [];

if (process.env.NODE_ENV !== 'test') {
  {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const r = require.context('./', true, /src\/index\.ts$/);
    r.keys().forEach((key: string) => {
      // 获取plugin文件夹的名字
      // get plugin's path name
      const attr = key.substring(
        key.indexOf('/') + 1,
        key.indexOf('/', key.indexOf('/') + 1),
      );
      plugins.push(r(key).default);
      pluginRecords[attr] = r(key).default;
    });
  }
  {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const r = require.context(
      './',
      true,
      /(src\/.*\/migrations\/.*\.ts$)|(src\/.*\/.*\.migration\.ts$)/,
    );
    r.keys().forEach((key: string) => {
      migrations.push(r(key).default);
    });
  }
  {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const r = require.context(
      './',
      true,
      /(src\/entities\/*\.ts$)|(src\/*\.entity\.ts$)/,
    );
    r.keys().forEach((key: string) => {
      entities.push(r(key).default);
    });
  }
}
