export default interface IExecutable {
  execute<T>(...args: any[]): Promise<T>;
}
