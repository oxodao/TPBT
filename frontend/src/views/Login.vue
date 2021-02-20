<template>
  <div class="login">
    <img src="../assets/logo_twitch.png" alt="Twitch logo" />
    <h3>Twitch Plays Blind test</h3>
    <a href="/auth/login">Login</a>
  </div>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator'

@Component({
  name: 'Login'
})
export default class Login extends Vue{}
</script>

<style lang="scss" scoped>
  .login {
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

  a {
    padding: .25em;
    text-decoration: none;
    margin-top: 2em;

    background: #30475e;
    border: none;
    color: #ececec;

    font-size: 1.4em;
    font-weight: bold;
  }

</style>
