interface Plugin {
  execute(...args: any[]): Promise<void>;
}

export default Plugin
