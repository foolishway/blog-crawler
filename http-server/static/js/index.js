$(function () {
    $(".share-img").click(function () {
            var p = $(this).parent();
            var articleId = p.siblings(".artice-id").text();
            var title = p.siblings(".title").text();
            var author = p.siblings(".author").text();
            var publishTime = p.siblings(".publish-time").text();
            var address = p.siblings(".title").find("a.address").attr("href")
            var msg = {
                articleId,
                title,
                author,
                address,
                publishTime
            };

            $.post("/share", JSON.stringify(msg), function (result) {
                p.append("<b class='shared'>已分享</b>").end().remove("img");
            }, "text")
    })
})