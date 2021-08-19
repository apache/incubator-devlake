export const migrations = [];

if (process.env.NODE_ENV !== 'test') {
  {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const r = require.context('./', true, /src\/migrations\/*\.ts$/);
    r.keys().forEach((key: string) => {
      migrations.push(r(key).default);
    });
  }
}
