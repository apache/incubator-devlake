export interface IExecutable<T> {
  execute(...args: any[]): Promise<T>;
}

export type Executable = {
  name: string;
  executable: IExecutable<any>;
};
