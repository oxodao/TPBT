<template>
  <div class="home">
      TODO: Add spinner
  </div>
</template>

<script>
export default {
  name: 'Home',
  components: {

  },
  mounted() {
    this.$store.dispatch("login")
  }
}
</script>
