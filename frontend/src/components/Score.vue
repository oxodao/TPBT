<template>
    <li>
        <img v-if="foundtitle" src="@/assets/tag.png" alt="Title found"/>
        <img v-if="foundartist" src="@/assets/user.png" alt="Artist found"/>
        <span>{{ username }}</span>
    </li>
</template>

<script>
    export default {
        name: "score",
        props: [ 'username', 'foundartist', 'foundtitle' ]
    }
</script>

<style lang="scss" scoped>
    li {
        display: flex;
        flex-direction: row;

        justify-content: center;
        align-items: center;

        img {
            height: 1em;
            margin-right: 1em;
        }
    }
</style>