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

    <div v-if="!loading && !successMessage && appdata.login.createAccount">
      <div class="field">
        <label class="label">Email</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input"
            :class="{ 'is-danger': !validEmail }"
            type="email"
            placeholder="Email input"
            @keypress.enter="submitClick"
            v-model="email"
            v-focus />
          <span class="icon is-small is-left">
            <fa-icon icon="envelope" />
          </span>
          <span class="icon is-small is-right" v-if="!validEmail">
            <fa-icon icon="exclamation-triangle" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validEmail">Email address invalid</p>
      </div>

      <div class="field">
        <label class="label">Username</label>
        <div class="control has-icons-left has-icons-right">
          <input class="input" type="text" placeholder="Text Input" v-model="username" @keyup="checkUsername" @keypress.enter="submitClick">
          <span class="icon is-small is-left">
            <fa-icon icon="user" />
          </span>
          <span class="icon is-small is-right" v-if="!validUsername">
            <fa-icon icon="exclamation-triangle" />
          </span>
        </div>
        <p class="help is-danger" v-if="!validUsername">
          Expected username to be between {{appdata.requirements.UsernameMinLength}} and {{appdata.requirements.UsernameMaxLength}} long
        </p>
        <p class="help is-danger" v-if="!validUsernameCharacters && validUsername">
          Username contains invalid characters. Expected {{appdata.requirements.UsernameRegex}}
        </p>
        <p class="help" v-if="checkingUsername">
          <fa-icon icon="circle-notch" spin /> Checking username...
        </p>
        <p class="help" v-if="!checkingUsername && usernameAvailableResponse">
          <template v-if="usernameAvailableResponse.exists">
            <fa-icon icon="exclamation-triangle" /> {{usernameAvailableResponse.username}} is already taken
          </template>
          <template v-else>
            <fa-icon icon="check" /> {{usernameAvailableResponse.username}} is available
          </template>
        </p>
      </div>

      <ValidatedPasswordField
        :minlength="appdata.requirements.PasswordMinLength"
        :maxlength="appdata.requirements.PasswordMaxLength"
        v-model="password"
        @valid="validPassword = $event"
        @enter="submitClick"
        />

      <RecaptchaV2 v-if="appdata.recaptchav2.enabled" :sitekey="appdata.recaptchav2.sitekey" :theme="appdata.recaptchav2.theme" ref="recaptchav2" />

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="!validForm">Submit</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="has-text-centered">
      <fa-icon icon="circle-notch" spin /> Creating account...
    </div>

    <div v-if="successMessage">
      <article class="message is-success">
        <div class="message-body">
          <fa-icon icon="check" /> {{successMessage}}
        </div>
      </article>
      <div v-if="willRedirect">
        <fa-icon icon="cog" spin /> Redirecting...
      </div>
    </div>

  </CenterCard>
</template>

<script>
import validator from 'validator';
import axios from 'axios';
import CenterCard from '../components/centerCard.vue';
import RecaptchaV2 from '../components/recaptchav2.vue';
import ValidatedPasswordField from '../components/fields/validatedPassword.vue';
import debounce from '../lib/debounce';

const errorCodes = {
  'username-unavailable': 'The username you have selected is unavailable',
  'account-email-exists': 'The email address you have entered is already associated with an account',
};

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
      password: '',
      validPassword: false,

      // responsive
      loading: false,
      error: null,
      successMessage: null,
      willRedirect: false,

      // live-update
      checkingUsername: false,
      usernameAvailableResponse: null,
    };
  },
  components: {
    CenterCard,
    RecaptchaV2,
    ValidatedPasswordField,
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
    validForm() {
      return this.validEmail && this.validPassword && this.validUsername && this.validUsernameCharacters;
    },
  },
  methods: {
    checkUsername: debounce(function checkUsername() {
      if (this.checkingUsername || !this.validUsername || !this.validUsernameCharacters) {
        return;
      }
      this.checkingUsername = true;
      axios.post('api/v1/account/check', { username: this.username })
        .then((resp) => {
          this.usernameAvailableResponse = resp.data;
        }).finally(() => {
          this.checkingUsername = false;
          if (this.usernameAvailableResponse.username !== this.username) {
            this.checkUsername();
          }
        });
    }, 200),
    submitClick() {
      this.error = null;

      if (!this.validForm) {
        this.error = 'Invalid form values';
        return;
      }

      const postData = {
        username: this.username,
        password: this.password,
        email: this.email,
      };
      const queryData = {
        createSession: true,
      };

      if (this.appdata.recaptchav2.enabled) {
        queryData.recaptchav2 = this.$refs.recaptchav2.getResponse();
        if (!queryData.recaptchav2) {
          this.error = 'Please accept RECAPTCHA';
          return;
        }
      }

      this.loading = true;
      axios.post('api/v1/account', postData, { params: queryData })
        .then((resp) => {
          if (resp.status !== 201) throw new Error('Error creating account');
          if (resp.data.createdSession) {
            this.successMessage = 'Account Successfully Created!';
            this.willRedirect = true;
            setTimeout(() => {
              this.$router.push('/login-redirect');
            }, 2.5 * 1000);
          } else {
            this.successMessage = 'Account created, but needs email verification before logging in. Please check your email.';
          }
        }).catch((err) => {
          if (err.response.data && err.response.data.reason && errorCodes[err.response.data.reason]) {
            this.error = errorCodes[err.response.data.reason];
          } else {
            this.error = `${err.message}`;
            if (err.response && err.response.data) {
              this.error += `: ${err.response.data.message}`;
            }
          }
        }).then(() => {
          this.loading = false;
        });
    },
  },
};
</script>
