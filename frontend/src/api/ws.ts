import {SERVER} from "./main";

function parseMessage(commit: Function, message: any) {
    console.log(message)
    const command = JSON.parse(message);
    const args    = JSON.parse(command.Arguments)

    switch(command.Command) {
        case 'SET_USER':
            commit('setUser', { data: args })
            if (args.Game) {
                commit('setGame', args.Game);
            }
            break;
        case 'SET_GAME':
            commit('setGame', args)
            break;
        case 'SET_LEADERBOARD':
            commit('setLeaderboard', args)
            break;
        case 'SET_FOUND':
            commit('setFound', args)
            break;
        default:
            console.log("UNHANDLED COMMAND: ", command.Command);
            break;
    }
}

export default function connectWebsocket(token: string, commit: Function): WebSocket {
    const protocol = window.location.protocol === "https:" ? "s" : "";
    const url = "ws" + protocol + "://" + location.host + "/connect/" + token;
    const ws     = new WebSocket(url)

    ws.onmessage = (e) => parseMessage(commit, e.data);
    ws.onerror   = (e) => console.log("Err: ", e);
    ws.onclose   = (e) => {
        console.log("Connection closed (ws.ts:8): ", e);
        commit('closedConnection', e.reason && e.reason.length > 0 ? e.reason : 'Connection closed');
    }

    window.setInterval(function() {
        ws.send(JSON.stringify({
            Command: 'PING',
            Arguments: '{}'
        }))
    }, 5000);

    return ws;
}

