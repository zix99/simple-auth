<template>
  <div>
    <Login @loggedIn="loginSuccess"/>
  </div>
</template>

<script>
import Login from './widgets/login.vue';

export default {
  components: {
    Login,
  },
  methods: {
    loginSuccess() {
      console.log("Logged in");
    },
  },
};
</script>
