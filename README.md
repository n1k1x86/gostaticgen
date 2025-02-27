# Here is my simple static page generator using .md files
## How to start ?
1. You need to install golang on your machine [Golang Official Site](https://go.dev/doc/install)
2. Then you need to clone main branch
## How to launch project
1. For the first, you have to create config dir with dirs (their name will become web pages' names)
2. The next, you need to create .md files with content and ONE yaml file with configs
3. The structure of YAML config (filename is a .md file name and each .md is a div block of web page):
```yaml
  content:
    -filename: block1
```
4. Then you need to start, you need to go to the main.go file and launch:
```go
  go run .\main.go -yamls "C:\Users\UserName\Desktop\config" -out_dir "C:\Users\UserName\Desktop\output" --create_out true
  params defenition:
  -yamls: dir with config
  -out_dir: dir with webpages
  -create_out: bool param, should programm create output dir or return an error
```
5. Examples
   
![image](https://github.com/user-attachments/assets/4cb3a577-635a-43de-8c4d-a880515de3e3)

![image](https://github.com/user-attachments/assets/afa6ff49-af0b-4505-a674-bbf5d57d0e61)
