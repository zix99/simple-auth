<template>
  <div>
    <article class="message is-danger" v-if="error">
      <div class="message-body">{{error}}</div>
    </article>

    <LoadingBanner :promise="signinPromise" :codes="errorCodes">
      Signing in...
    </LoadingBanner>

    <div v-if="state === 'login'">
      <div class="field">
        <label class="label">Username</label>
        <div class="control has-icons-left">
          <input class="input" type="text" placeholder="Text Input" v-model="username" @keypress.enter="submitClick" v-focus />
          <span class="icon is-small is-left">
            <fa-icon icon="user" />
          </span>
        </div>
      </div>

      <div class="field">
        <label class="label">Password</label>
        <div class="control has-icons-left">
          <input class="input" type="password" placeholder="Password" v-model="password" @keypress.enter="submitClick" />
          <span class="icon is-small is-left">
            <fa-icon icon="lock" />
          </span>
        </div>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="loading">Login</button>
        </div>
      </div>
      <router-link v-if="allowForgotPassword" to="forgot-password">Forgot Password?</router-link>
    </div>

    <div v-if="state === 'totp'">
      <p>
        Your account requires a two-factor login.  Please enter your token below to continue.
      </p>
      <div class="field">
        <label class="label">Token</label>
        <div class="control has-icons-left">
          <input class="input" type="text" placeholder="2FA Token" v-model="totp" @keypress.enter="submitClick" v-focus />
          <span class="icon is-small is-left">
            <fa-icon icon="key" />
          </span>
        </div>
      </div>
      <div class="field">
        <div class="control">
          <button class="button is-link" @click="submitClick">Login</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import axios from 'axios';
import LoadingBanner from '../components/loadingBanner.vue';

export default {
  components: {
    LoadingBanner,
  },
  props: {
    allowForgotPassword: null,
  },
  data() {
    return {
      error: null,
      username: '',
      password: '',
      totp: '',
      state: 'login',
      loading: false,
      signinPromise: null,
      errorCodes: {
        'totp-failed': 'Invalid 2FA Code',
        'invalid-credentials': 'Your username or password is invalid',
        'unsatisfied-stipulations': 'Your account has a hold on it',
        inactive: 'Your account is marked as inactive. Please contact an administrator if this is a mistake',
      },
    };
  },
  methods: {
    submitClick() {
      if (this.username === '' || this.password === '') {
        return;
      }

      this.loading = true;
      this.error = null;

      const postData = {
        username: this.username,
        password: this.password,
        totp: this.totp,
      };
      this.signinPromise = axios.post('api/v1/auth/session', postData)
        .then(() => {
          this.$emit('loggedIn');
        }).finally(() => {
          this.loading = false;
        }).catch((err) => {
          if (err.response.data.reason === 'totp-missing') {
            this.state = 'totp';
            return;
          }
          throw err;
        });
    },
    forgotPassword() {
      this.$emit('forgotPassword');
    },
  },
  watch: {
    state() {
      this.$emit('state', this.state);
    },
  },
};
</script>
