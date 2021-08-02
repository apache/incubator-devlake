import { Injectable } from "@nestjs/common";
import Plugin from "plugins/core/src";

export type JiraOptions = {}

@Injectable()
class Jira implements Plugin {
  async execute(options: JiraOptions) : Promise<void> {
    //TODO: Add jira collector and enrichment
  }
}

export default Jira;
