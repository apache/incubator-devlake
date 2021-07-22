// collection schema for gitlab 
// is reference to gitlab api response
export namespace Gitlab {
    interface User {
        id: number;
        username: string;
        name: string;
        state: 'active';
        avatar_url: string;
        web_url: string;
    }

    interface CommitStat {
        additions: number;
        deletions: number;
        total: number;
    }

    interface Commit {
        id: string;
        short_id: string;
        title: string;
        message: string;
        author_name: string;
        author_email: string;
        committer_name: string;
        committer_email: string;
        created_at: Date;
        committed_date: Date;
        authored_date: Date;
        parent_ids: string[];
        stats: CommitStat;
    }

    interface MileStone { }

    interface MergeRequest {
        id: number;
        iid: number;
        project_id: number;
        state: 'opened' | 'closed' | 'locked' | 'merged';
        title: string;
        description: string;
        merged_by?: User;
        merged_at?: Date;
        closed_by?: User;
        closed_at?: Date;
        created_at: Date;
        updated_at: Date;
        target_branch: string;
        source_branch: string;
        upvotes: number;
        downvotes: number;
        author: User;
        assignee: User;
        assignees: User[];
        reviewers: User[];
        source_project_id: number;
        target_project_id: number;
        labels: string[];
        draft: boolean;
        work_in_progress: boolean;
        milestone: MileStone;
        merge_when_pipeline_succeeds: boolean;
        merge_status: 'can_be_merged';
        sha: string;
        merge_commit_sha: string;
        squash_commit_sha: string;
        user_notes_count: number;
        discussion_locked?: boolean;
        should_remove_source_branch: boolean;
        force_remove_source_branch: boolean;
        allow_collaboration: boolean;
        allow_maintainer_to_push: boolean;
        web_url: string;
        squash: boolean;
        references: {
            short: string;
            relative: string;
            full: string;
        },
        time_stats: {
            time_estimate: number,
            total_time_spent: number,
            human_time_estimate?: number,
            human_total_time_spent?: number
        },
        // https://docs.gitlab.com/ee/api/merge_requests.html#get-single-mr-commits
        commits: Commit[];
        // https://docs.gitlab.com/ee/api/merge_requests.html#get-single-mr-participants
        participants: User[];
    }

    interface Contributor {
        name: string;
        email: string;
        commits: number;
        /**
         * @deprecated addtions always return 0
         */
        addtions: number;
        /**
         * @deprecated deletions always return 0
         */
        deletions: number;
    }

    interface Repository {
        // https://docs.gitlab.com/ee/api/commits.html#list-repository-commits
        commits: Commit[];
        // https://docs.gitlab.com/ee/api/repositories.html#contributors
        contributors: Contributor[];
    }

    interface Project extends Repository {
        id: number;
        name: string;
        description?: string;
        default_branch: string;
        visibility: 'private' | 'internal' | 'public';
        // https://docs.gitlab.com/ee/api/merge_requests.html#list-merge-requests
        merge_requests: MergeRequest[];
    }
}
