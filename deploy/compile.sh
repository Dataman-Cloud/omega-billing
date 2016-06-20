#!/bin/bash


# 获取 编译代码环境的镜像(比如go环境，python环境)
code_compile_image=`tail -n 1 $SERVICE/dockerfiles/Dockerfile_compile_env`
code_compile_image=${code_compile_image#*#}
# 生成 compile.sh
cat > /data/build/compile.sh << EOF
#!/bin/bash
mkdir -p /usr/local/go/src/github.com/Dataman-Cloud
rm -rf /usr/local/go/src/github.com/Dataman-Cloud/$SERVICE
cp -r $SERVICE /usr/local/go/src/github.com/Dataman-Cloud/
cd /usr/local/go/src/github.com/Dataman-Cloud/$SERVICE
make build
cp $SERVICE /data/build/$SERVICE/
EOF

chmod +x /data/build/compile.sh
