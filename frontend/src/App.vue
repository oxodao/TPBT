<template>
  <div id="app">
    <router-view/>

    <Banner />
    <Error/>
  </div>
</template>

<script>
  import {Component, Vue} from 'vue-property-decorator';
  import Error from './components/Error';
  import Banner from './components/Banner';

  @Component({
    name: 'App',
    components: {
      Banner,
      Error,
    }
  })
  export default class App extends Vue{}
</script>

<style lang="scss">
  body {
    margin: 0;
    width: 100vw;
    height: 100vh;
  }

  #app {
    margin: 0;
    width: 100%;
    height: 100%;

    display: flex;

    background: #222831;
    color: #ececec;
  }
</style>
