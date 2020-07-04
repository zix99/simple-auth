<template>
  <Card title="Create Account">
    <div v-if="!loading">
      <article class="message is-danger" v-if="error">
        <div class="message-body">{{error}}</div>
      </article>

      <div class="field">
        <label class="label">Email</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input"
            :class="{ 'is-danger': !validEmail }"
            type="email"
            placeholder="Email input"
            v-model="email" />
          <span class="icon is-small is-left">
            <i class="fas fa-envelope" />
          </span>
          <span class="icon is-small is-right" v-if="!validEmail">
            <i class="fas fa-exclamation-triangle" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validEmail">Email address invalid</p>
      </div>

      <div class="field">
        <label class="label">Username</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="text" placeholder="Text Input" v-model="username">
          <span class="icon is-small is-left">
            <i class="fas fa-user" />
          </span>
          <span class="icon is-small is-right" v-if="!validUsername">
            <i class="fas fa-exclamation-triangle" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validUsername">
          Expected username to be between {{requirements.UsernameMinLength}} and {{requirements.UsernameMaxLength}} long
        </p>
      </div>

      <div class="field">
        <label class="label">Password</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="password" placeholder="Password" v-model="password1" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validPassword">
          Expected password to be between {{requirements.PasswordMinLength}} and {{requirements.PasswordMaxLength}} long,
          and a score greater than 2 (Currently {{strength.score}})
        </p>
        <p class="help" v-if="strength.feedback">
          {{strength.feedback.warning}}
          <ul>
            <li v-for="suggestion in strength.feedback.suggestions" :key="suggestion">{{suggestion}}</li>
          </ul>
        </p>
      </div>
      <div class="field">
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="password" placeholder="Re-Enter Password" v-model="password2" />
          <span class="icon is-small is-left">
            <i class="fas fa-lock" />
          </span>
        </div>
        <p class="help is-danger" v-if="!passwordMatch">Password does not match</p>
      </div>

      <RecaptchaV2 v-if="recaptchav2.Enabled" :sitekey="recaptchav2.SiteKey" :theme="recaptchav2.Theme" ref="recaptchav2" />

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="!validEmail || !validPassword || !validUsername">Submit</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="has-text-centered">
      <i class="fas fa-circle-notch fa-spin" /> Creating account...
    </div>

  </Card>
</template>

<script>
import validator from 'validator';
import zxcvbn from 'zxcvbn';
import axios from 'axios';
import Card from './components/card.vue';
import RecaptchaV2 from './components/recaptchav2.vue';

export default {
  props: {
    requirements: {
      type: Object,
      default: () => ({
        UsernameMinLength: 1,
        UsernameMaxLength: 999,
        PasswordMinLength: 1,
        PasswordMaxLength: 999,
      }),
    },
    recaptchav2: {
      type: Object,
      default: () => ({
        Enabled: false,
        SiteKey: '',
        Theme: 'light',
      }),
    },
  },
  data() {
    return {
      // input
      email: '',
      username: '',
      password1: '',
      password2: '',

      // responsive
      strength: {},
      loading: false,
      error: null,
    };
  },
  components: {
    Card,
    RecaptchaV2,
  },
  computed: {
    validEmail() {
      return validator.isEmail(this.email);
    },
    validUsername() {
      return this.username.length >= this.requirements.UsernameMinLength
        && this.username.length <= this.requirements.UsernameMaxLength;
    },
    passwordMatch() {
      return this.password1 === this.password2;
    },
    validPassword() {
      return this.password1.length >= this.requirements.PasswordMinLength
        && this.password1.length <= this.requirements.PasswordMaxLength
        && this.strength.score >= 2;
    },
  },
  watch: {
    password1() {
      this.strength = zxcvbn(this.password1);
    },
  },
  methods: {
    submitClick() {
      this.loading = true;
      this.error = null;

      const postData = {
        username: this.username,
        password: this.password1,
        email: this.email,
      };

      if (this.recaptchav2.Enabled) {
        postData.recaptchav2 = this.$refs.recaptchav2.getResponse();
      }

      axios.post('/api/ui/account', postData)
        .then((resp) => {
          if (resp.status !== 201) throw new Error('Error creating account');
          // TODO: Show form
          this.error = 'Success';
        }).catch((err) => {
          this.error = `${err.message}`;
          if (err.response && err.response.data) {
            this.error += `: ${err.response.data.message}`;
          }
        }).then(() => {
          this.loading = false;
        });
    },
  },
};
</script>
