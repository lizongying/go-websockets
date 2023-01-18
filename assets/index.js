Window.websocket = new WebSocket('ws://127.0.0.1:1234/echo')
Window.websocket.onmessage = ({data}) => {
    const event = JSON.parse(data)
    console.log('received', event)

    // do something
    event.msg = 'ok'

    console.log('send', event)
    Window.websocket.send(JSON.stringify(event))
}
Window.websocket.onopen = () => {
    Window.websocket.send(JSON.stringify({
        from: 'chrome',
    }))
}