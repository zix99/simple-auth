<template>
  <div class="columns is-centered">
    <div class="column is-half">
      <div class="box has-text-centered">
        <h2 class="subtitle" v-if="!appdata.login.continue">{{appdata.company}} Account Management</h2>
        <div v-if="appdata.login.continue">
          <h2 class="subtitle">{{appdata.company}} Login</h2>
          <p v-if="isRemoteContinue">After logging in, you will be redirected to:<br /><strong>{{appdata.login.continue}}</strong></p>
        </div>

        <div class="has-text-left">
          <Login @loggedIn="$router.push('/login-redirect')" />
        </div>
        <div>
          <div v-for="oidc in appdata.oidc" :key="oidc.id" class="my-2">
            <OIDCButton :id="oidc.id" :icon="oidc.icon" class="is-info" :continue="appdata.login.continue">Continue with {{oidc.name}}</OIDCButton>
          </div>
        </div>
        <div v-if="appdata.login.createAccount">
          <p class="is-size-4">or</p>
          <div class="my-2">
            <router-link to="/create" class="button is-primary">Create Account</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import Login from '../widgets/login.vue';
import OIDCButton from '../components/oidcButton.vue';

export default {
  components: {
    Login,
    OIDCButton,
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
  computed: {
    isRemoteContinue() {
      return this.appdata.login.continue && this.appdata.login.continue.startsWith('http');
    },
  },
};
</script>
