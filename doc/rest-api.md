#omega-billing REST API

###GET `/api/v3/billing/list`
获取计量信息列表
```
curl -XGET localhost:5013/api/v3/billing/list?per_page=pcount&page=pnum&order=order&sort_by=sortby&appname=appname&cid=cid&starttime=starttime&endtime=endtime
```
