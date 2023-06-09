---
categories:
  - MYSQL
---
#### 到年底了，都需要统计下报表，mysql统计的时候有点问题，尤其是到12-31 到01-01 有问题
```php
select location as locationid
,if(DATE_FORMAT(date,'%v')='01' and date_format(date_add('1900-01-01',interval floor(datediff(date,'1900-01-01')/7)*7 day),'%Y')!=date_format(date_add('1900-01-01',interval floor(datediff(date,'1900-01-01')/7)*7+6 day),'%Y'),concat(date_format(date_add('1900-01-01',interval floor(datediff(date,'1900-01-01')/7)*7+6 day),'%Y'),'-',DATE_FORMAT(date,'%v')),DATE_FORMAT(date,'%X-%v')) weeks

,DATE_FORMAT(date,'%v') weeknum

,date_add('1900-01-01',interval floor(datediff(date,'1900-01-01')/7)*7 day) as week_start

, date_add('1900-01-01',interval floor(datediff(date,'1900-01-01')/7)*7+6 day) as week_end

,count(*) as count 

from resdata  where date >'2014-11-14 00:00:00' and date<'2019-12-14 23:59:59'  and location  = 5
  group by floor(datediff(date,'1900-01-01')/7),location
```

如下图  
![avatar](https://blog.hexiefamily.xin/assets/mysql1.jpg)  

#### 列出时间段中的所有周
```php
function list_middle_weeks($start,$end){
//    $start = "2014-10-14";
//    $end = "2019-12-14";

    //$first =1 表示每周星期一为开始日期 0表示每周日为开始日期
    $first=1;
    //获取当前周的第几天 周日是 0 周一到周六是 1 - 6
    $s_w=date('w',strtotime($start));
    //获取本周开始日期，如果$w是0，则表示周日，减去 6 天
    $s_week_start=date('Y-m-d',strtotime("$start -".($s_w ? $s_w - $first : 6).' days'));
    //本周结束日期
    $s_week_end=date('Y-m-d',strtotime("$s_week_start +6 days"));

    $weeknum = date('W', strtotime($start));
    // echo "week start:".$s_week_start.",week end:".$s_week_end.",weeks num is:".$weeknum;


    $e_w=date('w',strtotime($end));
    $e_week_start=date('Y-m-d',strtotime("$end -".($e_w ? $e_w - $first : 6).' days'));
    $e_week_end=date('Y-m-d',strtotime("$e_week_start +6 days"));
    $weeknum = date('W', strtotime($end));
    // echo "week start:".$e_week_start.",week end:".$e_week_end.",weeks num is:".$weeknum;


    $weeks = [];
    $startdate = $s_week_start;
    while(strtotime($startdate)<=strtotime($e_week_start)){
        $starttime = strtotime($startdate);
        $endtime =  strtotime("$startdate +6 days");
        $enddate = date("Y-m-d",$endtime);
        //startdate enddate
        $weeknum = date('W', $starttime);
        $firsty = date('Y',$starttime);
        $lasty = date('Y',$endtime);
        $year = $firsty;
        if($weeknum=="01" && $firsty!=$lasty){
            $year = $lasty;
        }
        $c = [];
        $c['week_start'] = $startdate;
        $c['week_end'] = $enddate;
        $weeks[$year."-".$weeknum] = $c;
        $startdate = date("Y-m-d",strtotime("$startdate +7 days"));
    }


//    echo json_encode($weeks);
    return $weeks;
}
```

