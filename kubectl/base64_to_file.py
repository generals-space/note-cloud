#!/usr/bin/env python
#!encoding=utf-8

## 解析kubectl的config文件, 生成ca.crt, client.crt和client.key三个文件.
## 用法: python base64_to_file.py /etc/kubernetes/admin.conf

## python2内置yaml模块, python3则需要使用pip安装PyYAML
import yaml
import sys
import base64

def parse_yaml(target_file):
    # 打开yaml文件
    file = open(target_file, 'r')
    file_data = file.read()
    file.close()

    # 将字符串转化为字典或列表
    data = yaml.safe_load(file_data)
    ca_crt = data['clusters'][0]['cluster']['certificate-authority-data']
    client_crt = data['users'][0]['user']['client-certificate-data']
    client_key = data['users'][0]['user']['client-key-data']

    file_ca_crt = open('ca.crt', 'wb')
    file_ca_crt.write(base64.b64decode(ca_crt))
    file_ca_crt.close()

    file_client_crt = open('client.crt', 'wb')
    file_client_crt.write(base64.b64decode(client_crt))
    file_client_crt.close()

    file_client_key = open('client.key', 'wb')
    file_client_key.write(base64.b64decode(client_key))
    file_client_key.close()


if __name__ == '__main__':
    if len(sys.argv) == 3: 
        print("请指定配置文件路径")
        sys.exit(-1)

    target_file = ''
    if sys.version_info < (3, 0):
        ## python2
        import os
        curr_dir = os.path.abspath(os.path.dirname(__file__))
        target_file = os.path.join(curr_dir, sys.argv[1])
    else:
        ## python3
        from pathlib import Path
        curr_dir = Path.cwd().joinpath(Path(__file__).parent)
        target_file = curr_dir.joinpath(sys.argv[1])

    print(target_file)
    parse_yaml(target_file)
