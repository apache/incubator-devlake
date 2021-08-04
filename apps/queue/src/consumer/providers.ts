import { Provider } from '@nestjs/common';
import Jira from 'plugins/jira/src';

const JiraProvider: Provider<Jira> = { provide: 'Jira', useClass: Jira };

export const providers = [JiraProvider];
