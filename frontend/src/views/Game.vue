<template>
    <div id="game">
        <header>
            <img :src="Picture" :alt="Username+'\'s picture'"/>
            <h4>Bienvenue {{Username}} (<span>{{Role}}</span>) <span v-if="!IsBotOnChannel" class="error">Attention ! Le bot n'est pas installé sur ta chaine.</span></h4>
            <div class="spanner"><!-- @TODO: Remove this and make it cleanly --></div>
            <button @click="openSettings">
                <img src="@/assets/cogs.png" alt="Settings"/>
            </button>
            <button>
                <img src="@/assets/sign-out.png" alt="Settings"/>
            </button>
        </header>

        <div class="game-inner">
            <div class="game-left">
                <div class="game-curr-found">
                    <h3>Personnes ayant trouvé ce round</h3>
                    <ul>
                        <Score v-for="u in Found" v-bind:key="u.Username+'/'+u.ArtistFound+'/'+u.TitleFound" :username="u.Username" :foundtitle="u.TitleFound" :foundartist="u.ArtistFound" />
                    </ul>
                </div>
                <div>
                    <h3>Leaderboard (100 personnes)</h3>
                    <ul>
                        <li v-for="sc in Leaderboard" v-bind:key="sc.Name+sc.Score">{{sc.Name}}: {{ sc.Score }} points</li>
                    </ul>
                </div>
            </div>
            <div class="game-center">
                <form>
                    <div class="simple-cb">
                        <input type="checkbox" id="hasTimer" :disabled="Game.Running" value="hasTimer" :checked="Game.IsTimed" v-model="Game.IsTimed"/>
                        <label for="hasTimer">Temps limité</label>
                    </div>
                    <div v-if="Game.IsTimed">
                        <label for="song-duration">Durée pour trouver (en secondes)</label>
                        <input type="number" id="song-duration" :disabled="Game.Running" min="1" v-model.number="Game.TimeSet"/>
                    </div>
                    <div>
                        <label for="song-title">Titre de la musique</label>
                        <input type="text" id="song-title" :disabled="Game.Running" v-model="Game.Title"/>
                    </div>

                    <div>
                        <label for="song-artist">Nom de l'artiste</label>
                        <input type="text" id="song-artist" :disabled="Game.Running" v-model="Game.Artist"/>
                    </div>

                    <div v-if="Game.IsTimed && Game.Running" class="timer">
                        <h4>Temps restant: </h4>
                        <span>{{getTimeLeft()}}</span>
                    </div>

                    <button type="button" @click="toggleTurn" v-if="!Game.Running">Lancer le tour</button>
                    <button type="button" @click="toggleTurn" v-else>Arrêter le tour</button>
                </form>
            </div>
        </div>

        <Settings/>
    </div>
</template>

<script>
    import {Component, Vue} from 'vue-property-decorator';
    import {mapGetters, mapState} from "vuex";

    import Settings from '../components/Settings';
    import Score from '../components/Score';

    @Component({
        name: 'Game',
        components: {
            Settings,
            Score,
        },
        computed: {
            ...mapState({
                IsBotOnChannel: state => state.User.IsBotInChannel,
                Leaderboard:    state => state.AppState.Leaderboard,
                Username:       state => state.User.DisplayName,
                Picture:        state => state.User.Picture,
                Found:          state => state.AppState.Found,
                Role:           state => state.User.Role,
                Game:           state => state.User.Game,
            }),
            ...mapGetters([
                'getTimeLeft'
            ])
        },
        methods: {
            openSettings() {
                this.$store.commit('toggleSettingsShown', true);
            },
            toggleTurn() {
                this.$store.dispatch('toggleTurn')
            }
        }
    })
    export default class CallbackTwitch extends Vue{
        mounted() {
            console.log(this.$store.state)
            if (!this.$store.state.IsLoggedIn) {
                this.$router.push({name:'Home'})
            }
        }
    }
</script>

<style lang="scss" scoped>
    #game {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;

        header {
            display: flex;
            flex-direction: row;
            align-items: center;

            height: 3em;
            background: darken(#222831, 5%);

            padding-left: 2em;

            & > img {
                width: 2em;
                border-radius: 50%;
                border: 2px solid white;
                margin-right: 1em;
            }

            button {
                background : darken(#222831, 8%);
                border: none;
                color: white;
                height: 100%;

                padding-left: 1em;
                padding-right: 1em;

                img {
                    height: 80%;
                }

                &:hover {
                    background: darken(#222831, 50%);
                }
            }
        }

        .game-inner {
            flex: 1;
            display: flex;
            flex-direction: row;
            min-height: 0;

            .game-left {
                width: 20%;
                display: flex;
                flex-direction: column;

                div {
                    flex: 1;
                    min-height: 0;
                    display: flex;
                    flex-direction: column;
                    text-align: center;
                    border-right: 1px solid;

                    ul {
                        flex: 1;
                        overflow-y: scroll;
                        list-style-type: none;
                        padding: 0;
                        margin: 0;
                    }
                }

                .game-curr-found {
                    border-bottom: 2px solid;
                }
            }

            .game-center {
                display: flex;
                flex: 1;
                align-items: center;
                justify-content: center;

                form {
                    height: 400px;
                    display: flex;
                    flex: 1;

                    flex-direction: column;
                    justify-content: space-around;
                    align-items: center;

                    div {
                        width: 300px;
                        & > * {
                            text-align: center;
                            display: block;
                            box-sizing: border-box;
                        }

                    }

                    label {
                        padding-top: .5em;
                    }

                    .timer {
                        padding: 1em;
                        font-size: 2em;

                        h4 {
                            margin: 0;
                        }
                    }


                    input[type="text"], input[type="number"] {
                        margin-top: 10px;
                        background: darken(#222831, 5%);
                        padding: .25em;
                        border: none;
                        color: aquamarine;
                        width: 100%;

                        font-size: 2em;
                    }

                    button {
                        background: darken(#222831, 5%);
                        border: 1px solid aquamarine;
                        padding: .25em;
                        font-size: 1.2em;
                        font-weight: bold;
                        color: aquamarine;
                    }
                }
            }
        }
    }

    .spanner {
        flex: 1;
    }

    .simple-cb {
        display: flex;
        justify-content: center;
        align-items: center;

        label {
            margin: 0 0 0 .5em !important;
            padding: 0 !important;
        }
    }

    .error {
        margin-left: 1em;
        color: #c0392b;
    }
</style>