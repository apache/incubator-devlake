export namespace Git {
    type Hash = string;
    type ShortHash = string;
    type Email = string;

    interface User {
        name: string;
        email: Email;
    }

    interface Remote {
        alias: string;
        protocol: 'file' | 'http' | 'ssh' | 'git';
        url: string;
    }

    // git log --shortstat
    interface CommitStats {
        additions: number;
        deletions: number;
        total: number;
    }

    // commit data should get from git repository (https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History)
    // or third-party platform api (eg: gitlab, https://docs.gitlab.com/ee/api/commits.html#get-a-single-commit)
    interface Commit {
        id: Hash;
        short_id: ShortHash;
        title: string; // the first line of commit message
        message: string; // the rest part of commit message
        parent_ids: Hash[];
        // difference between author and committer
        // https://stackoverflow.com/a/18754896/5760975
        author: User;
        author_date: Date;

        committer: User;
        committer_date: Date;

        stats: CommitStats;
    }

    interface Contributor extends User {
        commits: number;
        addtions: number;
        deletions: number;
    }

    interface Branch {
        name: string;
        ref: Hash;
    }

    /**
     * Repository represent a git repository
     */
    interface Repository {
        remotes: Remote[];
        commits: Commit[];
        branches: Branch[];
        // contributors can be enriched by commits data
        // or get from third-party api (eg: gitlab, https://docs.gitlab.com/ee/api/repositories.html#contributors)
        contributors: Contributor[];
    }
}
