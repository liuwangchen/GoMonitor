
/**
 * Author wzy
 * Date 2018/6/4
 * Desc 百分化行宽行高工具类
 */
function scaleHeight(percent) {
    return (document.body.clientHeight) * percent;
}

function scaleWidth(percent) {
    return (document.body.clientWidth - 5) * percent;
}
