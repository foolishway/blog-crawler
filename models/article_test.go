package models

import (
	"testing"
)

func TestInsertCollectArticles(t *testing.T) {
	articles := []Article{
		{Title: "移动端双指缩放图片JS事件的实践心得", Author: "张鑫旭", Address: "https://www.zhangxinxu.com/wordpress/2020/06/css-gap-history/"},
		{Title: "SVG任意图形path曲线路径的面积计算", Author: "张鑫旭", Address: "https://www.zhangxinxu.com/wordpress/2020/06/css-gap-history/"},
		{Title: "CSS columns轻松实现两端对齐布局效果", Author: "张鑫旭", Address: "https://www.zhangxinxu.com/wordpress/2020/06/css-gap-history/"},
	}
	if err := InsertCollectArticles(articles); err != nil {
		t.Fatal("Insert error.")
	}
}
