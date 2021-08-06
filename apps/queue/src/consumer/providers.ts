import { ClassProvider } from '@nestjs/common';
import Jira from 'plugins/jira/src';

const JiraProvider: ClassProvider<Jira> = {
  provide: 'Jira',
  useClass: Jira,
};

export const providers = [JiraProvider];
