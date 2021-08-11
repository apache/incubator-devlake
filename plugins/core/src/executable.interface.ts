export interface IExecutable<T> {
  execute(...args: any[]): Promise<T>;
}
