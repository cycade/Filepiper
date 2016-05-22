# FilePiper
## 0.6
### 新建数据库表(5.16完成)
- 原来的数据库表弄的很混乱而且里面还有好多测试文件，所以决定抛弃历史包袱，删掉重新来过。
#### Filepiper(Database)
1. fs.files(GridFS上传自动生成，是GridFS中存储文件的索引)
2. fs.chunks(GridFS上传自动生成，都是一些256k大小的分片文件)
3. metafiles(提取码与GridFS文件之间的对应）
#### 其中metafiles中应有一下字段：
1. ecode:提取码
2. filename:文件名称
3. md5:文件MD5
4. uploadDate:上传时间(根据时间判断是否失效)
5. downloadTimes:下载次数
6. isValid:是否失效

## 0.61
### 修改上传下载过程(5.17完成)
- 上传过程中前端表单上传文件，先连接到Filepiper数据库中，将上传文件存入GridFS中，生成四位提取码并确定metafiles中没有生成的提取码，再在metafiles中填充上传文件信息，并将其中的四位提取码返回给前端。
- 下载过程中前端收到四位提取码，首先对四位提取码进行校验，确定其为string类型并且string == 4，再在metafiles中进行查询，如没有则返回“提取码错误或已经失效”，如有则提取md5字段，并在fs.files中查找md5字段进行下载。

## 0.62
### 用Ajax重写前端与后端的数据交互(5.21完成)
- 包括但不限于，Ajax传递提取码，提取码错误值Ajax返回，正常下载文件。
- Ajax不能输出文件流的形式，所以验证码不合法的情况下返回错误信息，验证码正确的情况下应创建html表单，并提交进行下载。
- Ajax请求可以用jQuery发起。
