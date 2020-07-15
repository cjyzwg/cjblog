### Redis安装
- pecl install redis
- nano /etc/php/7.0/mods-available/redis.ini
extension=redis.so
Fpm 软连接
- ln -s /etc/php/7.0/mods-available/redis.ini
 /etc/php/7.0/fpm/conf.d/20-redis.ini
Cli同样也需要
- ln -s /etc/php/7.0/mods-available/redis.ini
 /etc/php/7.0/cli/conf.d/20-redis.ini
- php -m 即有
### Redis应用
[<font color=#0099ff>转发链接地址</font>](https://www.liaotaoo.cn/358.html)  

```php

<?php 
namespace app\controller;
 
 
use think\facade\Cache;
use think\facade\Db;
 
class Zan
{
    public $redis = null;
 
    //60*60*24/20=4320,每个点赞得到的分数。
    public $score = 0;
 
    //点赞增加数，或者
    public $num = 1;
 
 
    //init redis
    const COMMENT_RECORD = 'comment:record';
 
    public function __construct()
    {
        $this->redis = Cache::store('redis');
//        $this->redis->connect($this->redis_host,$this->redis_port);
//        $this->redis->auth($this->redis_pass);
    }
 
    /**
     * @param int $user_id 用户id
     * @param int $type 点击的类型 1.点like，2.点hate
     * @param int $comment_id 文章id
     * @return string json;
     */
    public function click($user_id,$type,$comment_id)
    {
        //判断是否需要更新数据
        if (!$this->ifUploadList($comment_id)) {
            return ['code'=>400,'msg'=>'文章不存在'];
        }
        $type == 1 ? $key = "like" : $key = "hate";
        //判断redis是否已经缓存了该文章数据
        //使用：分隔符对redis管理是友好的
        //这里使用redis zset-> zscore()方法
        $data = [];
        if ($this->redis->zscore("comment:like", $comment_id) || 
            $this->redis->zscore("comment:hate", $comment_id)) {
            //已经存在
            //判断点的是什么
            $rel = $this->redis->hget(self::COMMENT_RECORD, $user_id . ":" . $comment_id);
            //redis hash-> hget()
            //判断以前是否点过，点的是什么？
            switch ($rel) {
                case '':
                    //什么都没点过
                    //点赞加1
                    $this->redis->zincrby("comment:".$key, $this->num, $comment_id);
                    //记录上次操作
                    $this->redis->hset(self::COMMENT_RECORD, $user_id . ":" . $comment_id, $type);
 
                    $data = ["code" => 1,"msg" => $key."+1"];
                    break;
                case $type:
                    //点过赞了
                    //点赞减1
                    $this->redis->zincrby("comment:".$key, -($this->num), $comment_id);
                    //删除记录
                    $this->redis->hDel(self::COMMENT_RECORD, $user_id . ":" . $comment_id, $type);
                    //删除缓存 数据库再-1
                    $data = ["code" => 2,"msg" => $key."-1"];
                    break;
                case $type == 1 ? 2 : 1:
                    //点赞加1
                    $this->redis->zincrby("comment:".$key, $this->num, $comment_id);
                    //记录上次操作
                    $this->redis->hset(self::COMMENT_RECORD, $user_id . ":" . $comment_id, $type);
                    $data = ["code" => 1,"msg" => $key."+1"];
                    break;
 
            }
        } else {
            //未存在
            //点赞加一
            $this->redis->zincrby("comment:".$key, $this->num, $comment_id);
            $data = ["code" => 1,"msg" => $key."+1"];
            //记录
            $this->redis->hset(self::COMMENT_RECORD, $user_id . ":" . $comment_id, $type);
        }
        $list = $this->redis->zRange("comment:like",0,-1,1);
        foreach ($list as $k =>$v){if ($k==$comment_id){$data['like']=$v;}}
        return $data;
    }
 
    /**
     * 获取文章点赞数
     * @param $cid
     * @return mixed
     */
    public function getVodZanNum($cid)
    {
        $data = $this->redis->zRange("comment:like",0,-1,'WITHSCORES');
        if(empty($data)){
            $this->redis->zincrby("comment:like", $this->num, $cid);
            $data = Db::name('user_info_0')->where('id',$cid)->value('comic_id');
        }
        return $data;
    }
 
    /**
     * 点赞
     * @param $uid
     * @param $cid
     * @param int $type
     * @return bool
     */
    public function uidZan($uid,$cid)
    {
        $bool = $this->redis->sismember('zan'.$uid,$cid);
        if(!$bool){
            $this->redis->sadd('zan'.$uid,$cid);
            return true;
        }else{
            $this->redis->spop('zan'.$uid);
            return false;
        }
    }
 
    /**
     * 判断是否点过赞
     * @param $uid
     * @param $cid
     * @return bool
     */
    public function checkZan($uid,$cid)
    {
        $bool = $this->redis->sismember('zan'.$uid,$cid);
        if($bool) return true;
        return false;
    }
 
    /**
     * 取消点赞
     * @param $uid
     * @return bool
     */
    public function delZan($uid)
    {
        $this->redis->spop('zan'.$uid);
        return false;
    }
 
    /**
     * 判断文章是否存在
     * @param $comment_id
     * @return string
     */
    public function ifUploadList($comment_id)
    {
        date_default_timezone_set("Asia/Shanghai");
        $time = strtotime(date('Y-m-d H:i:s'));
        if(!$this->redis->sismember("comment:uploadset",$comment_id))
        {
            //文章不存在集合里，需要更新
            $this->redis->sadd("comment:uploadset",$comment_id);
            //更新到队列
            $data = ["id" => $comment_id,"time" => $time];
            $json = json_encode($data);
            $this->redis->lpush("comment:uploadlist",$json);
            return true;
        }
        return true;
    }
 
    public function dell(){
        $this->redis->flushAll();
    }
 
}
```
#### 调用
```php
$user_id = $_GET['uid'] ??25;
$comment_id= $_GET['cid'] ?? 1;
$like = 'like';
if (isset($like))
    $type = 1;
if ($user_id && $comment_id){
    $good = new Zan();
    $rel = $good->click($user_id,$type,$comment_id);
    if($rel['code'] == 1 || $rel['code'] == 7){
 
        Db::name('user_info_0')->where('id',$comment_id)->inc('comic_id')->update();
        echo '点赞成功';exit;
    }else{
        $where[] = ['id','=',$comment_id];
        $where[] = ['comic_id','>',0];
        Db::name('user_info_0')->where($where)->dec('comic_id')->update();
        $good->delZan($user_id);
        echo '取消点赞成功';
    }
 
}else{
    echo '用户ID和文章id不能为空';
}
```
