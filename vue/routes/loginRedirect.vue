<template>
  <div></div>
</template>

<script>
export default {
  props: {
    appdata: null,
  },
  created() {
    const continueUrl = this.appdata.login.continue;
    if (!continueUrl) {
      this.$router.push('manage');
    } else if (continueUrl.startsWith('#')) {
      this.$router.push(continueUrl.substring(1));
    } else {
      window.location = continueUrl;
    }
  },
};
</script>
