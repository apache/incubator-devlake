import os
from git import Repo

repo_path = '/path/to/your/local/repo'
file_path = os.path.join(repo_path, 'test.txt')

# 写入文件内容
with open(file_path, 'w') as file:
    file.write('hello\nhelloworld\n')

# 打开本地仓库
repo = Repo(repo_path)

# 添加文件到 Git
repo.index.add([file_path])

# 提交更改
repo.index.commit('Add test file with hello, helloworld content')

#我是增加1行注释
# 推送到远程仓库
origin = repo.remote(name='origin')
origin.push()

print('File committed and pushed to GitHub.')

#还删除了两行
