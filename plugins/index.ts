import Plugin from './core/src/plugin';

/**
 * 导入plugins目录下的所有plugin
 * import and export all plugins
 */
export const pluginRecords: Record<string, typeof Plugin> = {};
export const plugins: typeof Plugin[] = [];

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
