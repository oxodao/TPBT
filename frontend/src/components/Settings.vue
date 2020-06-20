<template>
    <div id="settings" :class="shown ? 'shown' : ''">
        <div class="title">
            <h1>Paramètres</h1>
            <button @click="close">X</button>
        </div>
        <div class="content">
            <p v-if="!isBotInChannel">Pour appeler le bot sur votre chaine, cliquez sur ce bouton. Il sera présent en permanence.</p>
            <p v-else>La déinstallation du bot n'impacte pas votre compte sur ce site, le bot ne viendra simplement plus sur votre chat.</p>
            <button type="button" @click="toggleBot" v-if="!isBotInChannel">Appeler le bot</button>
            <button type="button" @click="toggleBot" v-else>Désinstaller le bot</button>
        </div>
    </div>
</template>

<script>
    import {mapState} from "vuex";

    export default {
        name: "Settings",
        computed: {
            ...mapState({
                shown: state => state.AppState.SettingsShown,
                isBotInChannel: state => state.User.IsBotInChannel
            })
        },
        methods: {
            toggleBot() {
                this.$store.dispatch('toggleBot')
            },
            close() {
                this.$store.commit('toggleSettingsShown', false);
            }
        }
    }
</script>

<style lang="scss" scoped>
    $background: #29303b;
    $border-radius: .5em;

    #settings {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translateX(-50%) translateY(-50%);
        display: none;

        width: 300px;
        height: 300px;

        z-index: 100;

        background: $background;
        border: 1px solid black;
        border-radius: $border-radius;

        &.shown {
            display: flex;
            flex-direction: column;
        }

        .title {
            display: flex;
            background: #181c22;
            border-top-left-radius: $border-radius;
            border-top-right-radius: $border-radius;

            h1 {
                padding-left: 1em;
                font-size: 1.2em;
                flex: 1;
            }

            button {
                background: transparent;
                border: none;
                padding-left: 1em;
                padding-right: 1em;
                border-top-right-radius: $border-radius;
                color: white;

                &:hover {
                    background: lighten($background, 10%);
                }
            }

        }

        .content {
            flex: 1;
            padding: 1em;

            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;

            p {
                text-align: justify;
            }

            button {
                background: #181c22;
                border: 1px solid aquamarine;
                padding: .25em;
                font-size: 1.2em;
                font-weight: bold;
                color: aquamarine;
            }
        }
    }
</style>