export default interface IExecutable<T> {
  execute(...args: any[]): Promise<T>;
}
