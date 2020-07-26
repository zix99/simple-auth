<template>
  <CenterCard title="Create Account">
    <article class="message is-danger" v-if="error">
      <div class="message-body">{{error}}</div>
    </article>

    <article class="message is-warning" v-if="!appdata.login.createAccount">
      <div class="message-body">
        Account creation has been disabled on this server.
      </div>
    </article>

    <div v-if="!loading && !createdAccountId && appdata.login.createAccount">
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
          Expected username to be between {{appdata.requirements.UsernameMinLength}} and {{appdata.requirements.UsernameMaxLength}} long
        </p>
        <p class="help is-danger" v-if="!validUsernameCharacters && validUsername">
          Username contains invalid characters. Expected {{appdata.requirements.UsernameRegex}}
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
          Expected password to be between {{appdata.requirements.PasswordMinLength}} and {{appdata.requirements.PasswordMaxLength}} long,
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

      <RecaptchaV2 v-if="appdata.recaptchav2.enabled" :sitekey="appdata.recaptchav2.sitekey" :theme="appdata.recaptchav2.theme" ref="recaptchav2" />

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="!validEmail || !validPassword || !validUsername || !validUsernameCharacters">Submit</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="has-text-centered">
      <i class="fas fa-circle-notch fa-spin" /> Creating account...
    </div>

    <div v-if="createdAccountId">
      <article class="message is-success">
        <div class="message-body">
          <i class="fas fa-check"></i> Account Successfully Created!
        </div>
      </article>
      <div>
        <i class="fas fa-cog fa-spin"></i> Redirecting...
      </div>
    </div>

  </CenterCard>
</template>

<script>
import validator from 'validator';
import zxcvbn from 'zxcvbn';
import axios from 'axios';
import CenterCard from '../components/centerCard.vue';
import RecaptchaV2 from '../components/recaptchav2.vue';

export default {
  props: {
    appdata: {
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
      createdAccountId: null,
    };
  },
  components: {
    CenterCard,
    RecaptchaV2,
  },
  computed: {
    validEmail() {
      return validator.isEmail(this.email);
    },
    validUsername() {
      return this.username.length >= this.appdata.requirements.UsernameMinLength
        && this.username.length <= this.appdata.requirements.UsernameMaxLength;
    },
    validUsernameCharacters() {
      if (!this.appdata.requirements.UsernameRegex) return true;
      return this.username.match(this.appdata.requirements.UsernameRegex);
    },
    passwordMatch() {
      return this.password1 === this.password2;
    },
    validPassword() {
      return this.password1.length >= this.appdata.requirements.PasswordMinLength
        && this.password1.length <= this.appdata.requirements.PasswordMaxLength
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

      if (this.appdata.recaptchav2.enabled) {
        postData.recaptchav2 = this.$refs.recaptchav2.getResponse();
      }

      axios.post('/api/ui/account', postData)
        .then((resp) => {
          if (resp.status !== 201) throw new Error('Error creating account');
          this.createdAccountId = resp.data.id;
          setTimeout(() => {
            this.$router.push('/login-redirect');
          }, 2.5 * 1000);
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
