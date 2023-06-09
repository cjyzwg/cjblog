---
title: 解析ini文件
date: 2020-10-08 08:15:20
categories:
  - PHP
---
### a=none 后面为none解析不出来

sample.ini

```ini

##

[s1]

a=1

b=1

c=1

#d=1

e=fa#1

h = kda dawda

ddaw=

dwada

dwaihduihawu=dwadwa #dwa

[s2]

f=1

```

```php

//would see the version，php5.3 和大于php5.3  

$ini_array = parse_ini_file("sample.ini",true);   

//parse_ini_file解析出来的数组不一致，使用以下方法可解决  

$ini_array = parseconfig($ini_array);

print_r($ini_array);

function parseconfig($arr){

// print_r($arr);

if(count($arr)>0){

	foreach ($arr as $k => $v) {

		if(is_array($v)){

			$arr[$k] = parseconfig($v);

		}else{

			$pos = stripos($v," #");

			if($pos!==false){

				$arr[$k] = substr($v, 0,$pos);

			}

			if(substr(trim($k),0,1)=='#'){

				unset($arr[$k]);

			}

		}

	}

}

return $arr;

}

```

