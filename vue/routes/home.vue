<template>
  <div class="columns is-centered">
    <div class="column is-half">
      <div class="box has-text-centered">
        <h2 class="subtitle" v-if="!appdata.login.continue">{{appdata.company}} Account Management</h2>

        <div v-if="appdata.login.continue" class="content">
          <h2 class="subtitle">{{appdata.company}} Login</h2>
          <p v-if="isRemoteContinue">After logging in, you will be redirected to:<br /><strong>{{appdata.login.continue}}</strong></p>
          <p v-if="isOAuth2Continue">After logging in, you will redirect back to the Application Login Page</p>
        </div>

        <div class="box has-text-left">
          <Login
            @loggedIn="$router.push('/login-redirect')"
            @state="showAltLogin=false"
            :allowForgotPassword="appdata.login.forgotPassword"
            />
        </div>

        <div v-if="showAltLogin">
          <div>
            <div v-for="oidc in appdata.oidc" :key="oidc.id" class="my-2">
              <OIDCButton :id="oidc.id" :icon="oidc.icon" class="is-info" :continue="appdata.login.continue">Continue with {{oidc.name}}</OIDCButton>
            </div>
          </div>
          <div v-if="appdata.login.createAccount">
            <p class="is-size-5">or</p>
            <div class="my-2">
              <router-link to="/create" class="button is-primary">Create Account</router-link>
            </div>
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
  data() {
    return {
      showAltLogin: true,
    };
  },
  created() {
    axios.get('api/v1/account')
      .then(() => {
        // Is logged in!
        this.$router.push('/login-redirect');
      }).catch(() => {});
  },
  computed: {
    isRemoteContinue() {
      return this.appdata.login.continue && this.appdata.login.continue.startsWith('http');
    },
    isOAuth2Continue() {
      return this.appdata.login.continue && this.appdata.login.continue.startsWith('/#/oauth2');
    },
  },
};
</script>
