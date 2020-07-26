<template>
  <div>
    <article class="message is-danger" v-if="error">
      <div class="message-body">{{error}}</div>
    </article>

    <LoadingBanner v-if="loading">
      Signing in...
    </LoadingBanner>

    <div v-if="state === 'login'">
      <div class="field">
        <label class="label">Username</label>
        <div class="control has-icons-left">
          <input class="input" type="text" placeholder="Text Input" v-model="username">
          <span class="icon is-small is-left">
            <i class="fas fa-user" />
          </span>
        </div>
      </div>

      <div class="field">
        <label class="label">Password</label>
        <div class="control has-icons-left">
          <input class="input" type="password" placeholder="Password" v-model="password" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="loading">Login</button>
        </div>
      </div>
      <a href="#" @click.prevent="forgotPassword">Forgot Password?</a>
    </div>

    <div v-if="state === 'totp'">
      <p>Your account requires a two-factor login.  Please enter your token below to continue.</p>
      <div class="field">
        <label class="label">2FA</label>
        <div class="control has-icons-left">
          <input class="input" type="text" placeholder="2FA Token" v-model="totp" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
      </div>
      <div class="field is-grouped">
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
  data() {
    return {
      error: null,
      username: '',
      password: '',
      totp: '',
      state: 'login',
      loading: false,
    };
  },
  methods: {
    submitClick() {
      this.loading = true;
      this.error = null;

      const postData = {
        username: this.username,
        password: this.password,
        totp: this.totp,
      };
      axios.post('api/ui/login', postData)
        .then(() => {
          this.$emit('loggedIn');
        }).catch((err) => {
          this.error = err.message;
        }).then(() => {
          this.loading = false;
        });
    },
    forgotPassword() {
      this.$emit('forgotPassword');
    },
  },
};
</script>
