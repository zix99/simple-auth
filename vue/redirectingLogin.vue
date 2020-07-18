<template>
  <div>
    <Login @loggedIn="loginSuccess" />
  </div>
</template>

<script>
import Login from './widgets/login.vue';

export default {
  properties: {
    successUrl: null,
  },
  components: {
    Login,
  },
  methods: {
    loginSuccess() {
      if (this.successUrl) {
        window.location = this.successUrl;
      } else {
        window.location = '/';
      }
    },
  },
};
</script>
