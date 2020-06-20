import Vue                from 'vue'
import Vuex               from 'vuex'
import Router             from '../router/index';
import {APICheckValidity} from "@/api/login";
import connectWebsocket   from "@/api/ws";

Vue.use(Vuex)

export default new Vuex.Store({
    state:     {
        IsLoggedIn: false,
        User:       {
            DisplayName:    "",
            Picture:        "",
            IsBotInChannel: false,
            Role:           "",
            AppToken:       "",
            Game:           {
                Title:    "",
                Artist:   "",
                Running:  false,
                IsTimed:  false,
                TimeLeft: 120,
                TimeSet: 120,
            }
        },
        AppState:   {
            SettingsShown: false,
            Websocket:     (null as WebSocket|null), // Ffs typefuckingscript
            ClosedReason:  null,
            Leaderboard:   [],
            Found: [],
        }
    },
    mutations: {
        closedConnection:    (state, message?) => {
            state.AppState.ClosedReason = message;
        },
        setUser:             (state, {data, ws}) => {
            if (ws) {
                state.AppState.Websocket  = ws;
            }
            state.IsLoggedIn          = true;

            state.User.Role           = data.Role;
            state.User.AppToken       = data.AppToken;
            state.User.Picture        = data.Picture;
            state.User.DisplayName    = data.DisplayName;
            state.User.IsBotInChannel = data.IsBotInChannel;


            if (state.User.AppToken) {
                localStorage.setItem('token', state.User.AppToken);
            } else {
                alert("Impossible de récupérer les informations de connexion!")
            }
        },
        setGame: (state, payload) => {
            state.User.Game.Title    = payload.Title;
            state.User.Game.Artist   = payload.Artist;
            state.User.Game.Running  = payload.Running;
            state.User.Game.IsTimed  = payload.IsTimed;
            state.User.Game.TimeLeft = payload.TimeLeft;
            state.User.Game.TimeSet  = payload.TimeSet;
        },
        toggleSettingsShown: (state, payload: boolean) => {
            state.AppState.SettingsShown = payload;
        },
        setLeaderboard: (state, scores) => {
            Vue.set(state.AppState, 'Leaderboard', scores);
        },
        setFound: (state, scores) => {
            Vue.set(state.AppState, 'Found', scores);
        }
    },
    actions:   {
        login() {
            const token = localStorage.getItem('token') || '';
            if (token.length > 0) {
                APICheckValidity(token, (code: number, data?: any) => {
                    if (code === 200) {
                        this.dispatch('connect', data)
                    } else if (code === 401) {
                        localStorage.removeItem('token');
                        Router.push({name: 'Login'});
                    } else if (code === -1) {
                        alert("Impossible d'accéder au serveur")
                    } else {
                        alert("Something went wrong!")
                    }
                });
            } else {
                Router.push({name: 'Login'})
            }
        },
        connect({commit}, data) {
            commit('setUser', {
                data, ws: connectWebsocket(data.AppToken, commit)
            });

            if (data.Game) {
                commit('setGame', data.User.Game);
            }

            Router.push({name: 'Game'})
        },
        toggleBot({state}) {
            console.log("Toggle bot")
            state.AppState.Websocket!.send(JSON.stringify({
                Command: 'TOGGLE_BOT',
                Arguments: '{}'
            }))
        },
        toggleTurn({state, getters}){
            if (!getters.IsAllowed) { return }

            let args = {};

            if (!state.User.Game.Running) {
                console.log(`Setting the turn to ${state.User.Game.Title} - ${state.User.Game.Artist} (${state.User.Game.IsTimed ? "Timed, " + state.User.Game.TimeSet : ""})`)
                args = state.User.Game;
            }

            state.AppState.Websocket!.send(JSON.stringify({
                Command: 'TOGGLE_TURN',
                Arguments: JSON.stringify(args)
            }))
        }
    },
    modules:   {},
    getters: {
        getTimeLeft: (state) => () => {
            const hours   = ~~(state.User.Game.TimeLeft / 3600);
            const minutes = ~~((state.User.Game.TimeLeft % 3600) / 60);
            const seconds = ~~state.User.Game.TimeLeft % 60;

            return `${(hours < 10 ? "0" : "") + hours}:${(minutes < 10 ? "0" : "") + minutes}:${(seconds < 10 ? "0" : "") + seconds}`;
        },
        IsAllowed: (state) => () => {
            return state.User.Role.toLowerCase() !== 'visiteur' &&  state.User.Role.toLowerCase() !== 'bannis';
        },
        UnallowedReason: (state) => () => {
            switch (state.User.Role.toLowerCase()) {
                case "visiteur":
                    return "Votre compte n'est pas activé (Visiteur). Contactez un admin pour avoir la permission d'utiliser ce logiciel."
                case "bannis":
                    return "Votre compte à été bannis. Vous n'êtes plus en mesure d'utiliser ce logiciel."
                default:
                    return "";
            }
        }
    }
})
