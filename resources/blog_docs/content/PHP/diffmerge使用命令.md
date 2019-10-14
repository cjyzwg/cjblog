### diffmerge使用
> 以mac为例  
http://www.sourcegear.com/diffmerge/  
osx pkg版本  

##### 命令用法
>$ diffmerge fold1 fold2  

##### compare.sh
```shell script
#!/bin/bash
GetKey(){    
   section=$(echo $1 | cut -d '.' -f 1)    
    key=$(echo $1 | cut -d '.' -f 2)    
    sed -n "/\[$section\]/,/\[.*\]/{    
    /^\[.*\]/d    
    /^[ \t]*$/d    
    /^$/d    
    /^#.*$/d    
    s/^[ \t]*$key[ \t]*=[ \t]*\(.*\)[ \t]*/\1/p    
    }" config.ini    
}    
    
#读取实例  
fold1=$(GetKey "compare.fold1")
fold2=$(GetKey "compare.fold2")
echo "========================"
printf "%s\n" $fold1
printf "%s\n" $fold2
echo "========================"
diffmerge $fold1 $fold2
echo done
```
##### config.ini
```markdown
[compare]
fold1=/Home/cj/fold1
fold2=/Home/cj/fold2
```

##### 如果想比较完第一个文件夹，再比较第二个文件夹，diffmerge提供一个退出状态,在shell脚本下追加以下命令即可。
```shell script
echo "========================"
diffmerge $afold1 $afold2
if [[ $? == 0 ]]
then
    diffmerge $bfold1 $bfold2
else
    echo "a 比较失败"
fi
if [[ $? == 0 ]]
then
    diffmerge $cfold1 $cfold2
else
    echo "b 比较失败"
fi
echo done
```
