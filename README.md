# epub_online

## API

### 获取图书列表

```
GET /books
```

response body example:
```
[
     {
        "navigations":[
            {
                "title":"版权页",
                "charactor_count":3,
                "total_charactor_count":3,
                "level":1,
                "tag":"1"
            },
            {
                "title":"共产党宣言",
                "charactor_count":685,
                "total_charactor_count":40380,
                "level":1,
                "tag":"2"
            },
            {
                "title":"1872年德文版序言",
                "charactor_count":1304,
                "total_charactor_count":1304,
                "level":2,
                "tag":"2.1"
            },
            {
                "title":"1882年俄文版序言",
                "charactor_count":1303,
                "total_charactor_count":1303,
                "level":2,
                "tag":"2.2"
            },
            {
                "title":"1883年德文版序言",
                "charactor_count":699,
                "total_charactor_count":699,
                "level":2,
                "tag":"2.3"
            },
            {
                "title":"1888年英文版序言",
                "charactor_count":4033,
                "total_charactor_count":4033,
                "level":2,
                "tag":"2.4"
            },
            {
                "title":"1890年德文版序言",
                "charactor_count":3827,
                "total_charactor_count":3827,
                "level":2,
                "tag":"2.5"
            },
            {
                "title":"1892年波兰文版序言",
                "charactor_count":1161,
                "total_charactor_count":1161,
                "level":2,
                "tag":"2.6"
            },
            {
                "title":"1893年意大利文版序言",
                "charactor_count":1401,
                "total_charactor_count":1401,
                "level":2,
                "tag":"2.7"
            },
            {
                "title":"共产党宣言",
                "charactor_count":692,
                "total_charactor_count":692,
                "level":2,
                "tag":"2.8"
            },
            {
                "title":"一　资产者与无产者",
                "charactor_count":8850,
                "total_charactor_count":8850,
                "level":2,
                "tag":"2.9"
            },
            {
                "title":"二　无产者与共产党人",
                "charactor_count":6949,
                "total_charactor_count":6949,
                "level":2,
                "tag":"2.10"
            },
            {
                "title":"三　社会主义的与共产主义的著作",
                "charactor_count":15,
                "total_charactor_count":7132,
                "level":2,
                "tag":"2.11"
            },
            {
                "title":"1．反动的社会主义",
                "charactor_count":9,
                "total_charactor_count":4072,
                "level":3,
                "tag":"2.11.1"
            },
            {
                "title":"（1）封建的社会主义",
                "charactor_count":1416,
                "total_charactor_count":1416,
                "level":4,
                "tag":"2.11.1.1"
            },
            {
                "title":"（2）小资产阶级的社会主义",
                "charactor_count":822,
                "total_charactor_count":822,
                "level":4,
                "tag":"2.11.1.2"
            },
            {
                "title":"（3）德国的或“真正的”社会主义",
                "charactor_count":1825,
                "total_charactor_count":1825,
                "level":4,
                "tag":"2.11.1.3"
            },
            {
                "title":"2．保守的或资产者的社会主义",
                "charactor_count":835,
                "total_charactor_count":835,
                "level":3,
                "tag":"2.11.2"
            },
            {
                "title":"3．批判的空想的社会主义与共产主义",
                "charactor_count":2210,
                "total_charactor_count":2210,
                "level":3,
                "tag":"2.11.3"
            },
            {
                "title":"四　共产党人对各种反对党的态度",
                "charactor_count":2344,
                "total_charactor_count":2344,
                "level":2,
                "tag":"2.12"
            },
            {
                "title":"译后记",
                "charactor_count":2107,
                "total_charactor_count":2107,
                "level":1,
                "tag":"3"
            }
        ],
        "meta":{
            "contributor":"",
            "coverage":"",
            "creator":"马克思　恩格斯　著",
            "date":"1978-11",
            "description":"",
            "format":"",
            "identifier":"1001·1165",
            "language":"zh-cn",
            "meta":"images_cover_jpg",
            "publisher":"人民出版社",
            "relation":"",
            "rights":"",
            "source":"",
            "subject":"",
            "title":"共产党宣言_B_1978_000042",
            "type":"图书"
        },
        "charactor_count":42490,
        "url":"epub/gcdxy.epub"
    }
]
```

字段说明:
- navigations: 章节目录
> - title: 目录名称
> - charactor_count: 本页字数
> - total_charactor_count: 本章总字数
> - level: 目录等级
> - tag: 目录标记

- charactor_count: 图书总字数
- url: 图书下载链接