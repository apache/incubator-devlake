var config = {
  github: {repo: 'lake', owner: 'merico-dev', baseUrl: 'https://api.github.com'}
}

module.exports = [
  { name: 'github-commits', url: `${config.github.baseUrl}/repos/${config.github.owner}/${config.github.repo}/commits`, paging: true },
  { name: 'github-issues', url: `${config.github.baseUrl}/repos/${config.github.owner}/${config.github.repo}/issues`, paging: true }
]
