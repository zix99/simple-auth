<template>
  <div class="columns is-centered">
    <div class="column is-half">
      <div class="box has-text-centered">
        <h2 class="subtitle">{{appdata.company}} Account Management</h2>
        <div v-if="appdata.login.createAccount">
          <router-link to="/create" class="button is-primary is-light">Create Account</router-link>
          <p class="is-size-4">or</p>
        </div>
        <div class="has-text-left">
          <Login @loggedIn="$router.push('/login-redirect')" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Login from '../widgets/login.vue';

export default {
  components: {
    Login,
  },
  props: {
    appdata: null,
  },
  created() {
    axios.get('/api/ui/account')
      .then(() => {
        // Is logged in!
        this.$router.push('/login-redirect');
      }).catch(() => {});
  },
};
</script>
