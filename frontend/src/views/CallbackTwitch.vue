<template>
    <div class="callback">
        <img alt="Logo twitch" src="../assets/logo_twitch.png" />
        <span v-if="code === undefined || code === null || code.length === 0">Something went wrong authenticating! Contact the dev</span>
        <div v-else class="box">
            <div class="spinner"></div>
        </div>
    </div>
</template>


<script lang="ts">
    import {Component, Vue} from 'vue-property-decorator';
    import {APILogin} from '@/api/login';

    @Component({
        name: 'CallbackTwitch',
        data() {
            return {
                code: this.$route.query.code,
                scope: this.$route.query.scope,
                state: this.$route.query.state,
            };
        },
    })
    export default class CallbackTwitch extends Vue{
        mounted() {
            this.init();
        }

        public init(): void {
            APILogin(this.$store.dispatch, () => {this.$router.push({name: 'Home'}) }, this.$route.query.code.toString())
        }
    }
</script>

<style lang="scss" scoped>
    .callback {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;

        width: 100%;
        height: 100%;
    }

    img {
        width: 15em;
    }

    button {
        margin-top: 2em;

        background: #30475e;
        border: none;
        color: #ececec;

        font-size: 1.4em;
        font-weight: bold;
    }

    span {
        margin-top: 2em;
    }

    .box {
        margin-top: 2em;
    }

    /** https://codepen.io/dev_loop/pen/mgeYwo */
    .box .spinner {
        height: 40px;
        width: 40px;
        background: rgba(0, 0, 0, .2);
        border-radius: 50%;

        transform: scale(0);
        background: rgba(0, 0, 0, .8);
        opacity: 1;
        animation: spinner4 800ms linear infinite;
    }

    @keyframes spinner4 {
        to {
            transform: scale(1.5);
            opacity: 0;
        }
    }

</style>

