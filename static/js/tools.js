	function get_host_url()
	{
	    var patt1=new RegExp("http://[a-zA-Z0-9\.]+:[0-9]+");
	    var crt_href = top.location.href;
	    var r = patt1.exec(crt_href);   
	    return r[0];
	}
    /// <summary>
    /// 格式化显示日期时间
    /// </summary>
    /// <param name="x">待显示的日期时间，例如new Date()</param>
    /// <param name="y">需要显示的格式，例如yyyy-MM-dd hh:mm:ss</param>
    function date2str(x,y) {
        var z = {M:x.getMonth()+1,d:x.getDate(),h:x.getHours(),m:x.getMinutes(),s:x.getSeconds()};
        y = y.replace(/(M+|d+|h+|m+|s+)/g,function(v) {return ((v.length>1?"0":"")+eval('z.'+v.slice(-1))).slice(-2)});
        return y.replace(/(y+)/g,function(v) {return x.getFullYear().toString().slice(-v.length)});
    }
        // alert(date2str(new Date(),"yyyy-MM-dd hh:mm:ss"));
    // alert(date2str(new Date(),"yyyy-M-d h:m:s"));


Date.prototype.format = function(format)  
{  
    /* 
    * format="yyyy-MM-dd hh:mm:ss"; 
    */  
    var o = {  
    "M+" : this.getMonth() + 1,  
    "d+" : this.getDate(),  
    "h+" : this.getHours(),  
    "m+" : this.getMinutes(),  
    "s+" : this.getSeconds(),  
    "q+" : Math.floor((this.getMonth() + 3) / 3),  
    "S" : this.getMilliseconds()  
    }  
    if (/(y+)/.test(format))  
    {  
        format = format.replace(RegExp.$1, (this.getFullYear() + "").substr(4  
    - RegExp.$1.length));  
    }  
      
    for (var k in o)  
    {  
        if (new RegExp("(" + k + ")").test(format))  
    {  
        format = format.replace(RegExp.$1, RegExp.$1.length == 1  
        ? o[k]  
        : ("00" + o[k]).substr(("" + o[k]).length));  
    }  
    }  
    return format;  
}  

// exports.getCurrentTime = function(){
//     var time = new Date();  
//     return time.format('yyyy-MM-dd hh:mm:ss');

// };