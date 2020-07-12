$(function () {
    $("#insert").click(function () {
        window.location.href = "/dutyedit/?id=0"
        return false;
    })

    $(".update").click(function () {
        var id = $(this).parent("td").siblings().eq(0).text()
        window.location.href = "/dutyedit/?id="+id
        return false;
    })

    $(".delete").click(function () {
        var id = $(this).parent("td").siblings().eq(0).text()
        var name = $(this).parent("td").siblings().eq(1).text()
        if (window.confirm("确定要删除" + name + "？")) {
            $.post("/delete", "id="+id, function (result) {
                window.location.reload()
            }, "text")
        }
        return false;
    })

})
function validCheck() {
    var name = $("#name").val();
    var employeeNum = $("#employeeNum").val();
    var phone = $("#phone").val();
    if (name == "" || employeeNum == "" || phone == "") {
        alert("信息不完整。")
        return false;
    }
}