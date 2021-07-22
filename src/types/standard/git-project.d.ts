import { Base } from "./base";
import { Git } from './git';
export namespace GitProject {

    interface User extends Git.User {
        id: string; // id is under project scope
        username: string;
    }

    interface Comment extends Base.TimeProps {
        author: User;
    }

    interface MergeRequest extends Base.TimeProps {
        title: string;
        description: string;
        stage: 'merged' | 'open' | 'closed';
        source_branch: string;
        target_branch: string;
        merge_time?: Date;
        upvotes?: number;
        downvotes?: number;
        author: User;
        comments: Comment[];
        commits: Git.Commit[];
        // participants can be enriched by comments
        // or get from third-party api (eg: https://docs.gitlab.com/ee/api/merge_requests.html#get-single-mr-participants)
        participants: User[];
    }

    interface Issue extends Base.TimeProps { }

    /**
     * Project = Repository & Some ProjectManager Data
     * eg: Merge Request, Code Review, etc.
     */
    interface Project extends Git.Repository {
        merge_requests: MergeRequest[];
        issues: Issue[];
    }
}
