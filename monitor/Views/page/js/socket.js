var cpusocket = null;
function initSocket(){
    cpusocket = new WebSocket('ws://'+window.location.host+'/monitorCpu');

    cpusocket.onerror = function (event) {
        onError(event);
    };
    cpusocket.onopen = function (event) {
        onOpen(event);
    };
    cpusocket.onmessage = function (event) {
        alert(event.data);
        onMessage(event);
    };
}

function onError(event){

}

function onOpen(event){
    cpusocket.send();//看后台需要接收什么信息才能握手成功
}

function onMessage(event){
    console.log(event);
}