<template>
    <div v-if="ClosedReason" class="error">
        <h1>Une erreur est survenue !</h1>
        <p>{{ClosedReason}}</p>
        <button @click="reload">Relancer</button>
    </div>
</template>

<script>
    import {mapState} from "vuex";

    export default {
        name: "Error",
        computed: mapState({
            ClosedReason: state => state.AppState.ClosedReason,
        }),
        methods: {
            reload() {
                localStorage.removeItem('token');
                window.location.assign("/");
            }
        }
    }
</script>

<style lang="scss" scoped>
    .error {
        position: absolute;

        top: 0;
        left: 0;
        right: 0;
        bottom: 0;

        z-index: 999999;

        background: rgba(0, 0, 0, .8);

        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
</style>