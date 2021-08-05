const mockRawCommits = [
  {
    id: '926ccda073f04f12eb7ab373c84cb73cff4ce238',
    short_id: '926ccda0',
    created_at: '2021-07-27T11:09:19.000-03:00',
    parent_ids: [
      'a64bb3f50c2993b4517edfecb36564afc407f197'
    ],
    title: 'goodbye message',
    message: 'goodbye message\n',
    author_name: 'Kevin Kline',
    author_email: 'kevinkline@Kevins-MacBook-Pro.local',
    authored_date: '2021-07-27T11:09:19.000-03:00',
    committer_name: 'Kevin Kline',
    committer_email: 'kevinkline@Kevins-MacBook-Pro.local',
    committed_date: '2021-07-27T11:09:19.000-03:00',
    trailers: {},
    web_url: 'https://gitlab.com/kevin-kline/test-project/-/commit/926ccda073f04f12eb7ab373c84cb73cff4ce238',
    stats: {
      additions: 1,
      deletions: 0,
      total: 1
    }
  },
  {
    id: 'a64bb3f50c2993b4517edfecb36564afc407f197',
    short_id: 'a64bb3f5',
    created_at: '2021-07-22T10:16:38.000-03:00',
    parent_ids: [
      'b42ba345c6d032ee1799545c03bc1d3dffd8b04b'
    ],
    title: 'hello',
    message: 'hello\n',
    author_name: 'Kevin Kline',
    author_email: 'kevinkline@Kevins-MacBook-Pro.local',
    authored_date: '2021-07-22T10:16:38.000-03:00',
    committer_name: 'Kevin Kline',
    committer_email: 'kevinkline@Kevins-MacBook-Pro.local',
    committed_date: '2021-07-22T10:16:38.000-03:00',
    trailers: {},
    web_url: 'https://gitlab.com/kevin-kline/test-project/-/commit/a64bb3f50c2993b4517edfecb36564afc407f197',
    stats: {
      additions: 4,
      deletions: 0,
      total: 4
    }
  },
  {
    id: 'b42ba345c6d032ee1799545c03bc1d3dffd8b04b',
    short_id: 'b42ba345',
    created_at: '2021-07-22T10:15:11.000-03:00',
    parent_ids: [],
    title: 'add README',
    message: 'add README\n',
    author_name: 'Kevin Kline',
    author_email: 'kevinkline@Kevins-MacBook-Pro.local',
    authored_date: '2021-07-22T10:15:11.000-03:00',
    committer_name: 'Kevin Kline',
    committer_email: 'kevinkline@Kevins-MacBook-Pro.local',
    committed_date: '2021-07-22T10:15:11.000-03:00',
    trailers: {},
    web_url: 'https://gitlab.com/kevin-kline/test-project/-/commit/b42ba345c6d032ee1799545c03bc1d3dffd8b04b',
    stats: {
      additions: 0,
      deletions: 0,
      total: 0
    }
  }
]

module.exports = mockRawCommits
