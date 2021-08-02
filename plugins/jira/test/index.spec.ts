import Jira, {JiraOptions} from "../src";


describe('Jira Plugin', () => {
  it('Jira', async () => {
    const jira = new Jira();
    expect(jira).toBeDefined();
    expect(jira.execute).toBeDefined();
  });
})