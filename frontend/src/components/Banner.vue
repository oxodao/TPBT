<template>
    <div v-if="!IsAllowed()" class="error">
        <h1>Vous n'avez pas accès à cette application</h1>
        <p>{{UnallowedReason()}}</p>
    </div>
</template>

<script>
    import {mapGetters} from "vuex";

    export default {
        name: "Banner",
        computed: mapGetters([
            'IsAllowed',
            'UnallowedReason'
        ]),
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