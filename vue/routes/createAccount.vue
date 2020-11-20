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
            ref="email"
            :class="{ 'is-danger': !validEmail }"
            type="email"
            placeholder="Email input"
            v-model="email" />
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
          <input class="input" type="text" placeholder="Text Input" v-model="username">
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
      </div>

      <ValidatedPasswordInput
        :minlength="appdata.requirements.PasswordMinLength"
        :maxlength="appdata.requirements.PasswordMaxLength"
        v-model="password"
        @valid="validPassword = $event"
        />

      <RecaptchaV2 v-if="appdata.recaptchav2.enabled" :sitekey="appdata.recaptchav2.sitekey" :theme="appdata.recaptchav2.theme" ref="recaptchav2" />

      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="!validEmail || !validPassword || !validUsername || !validUsernameCharacters">Submit</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="has-text-centered">
      <fa-icon icon="circle-notch" spin /> Creating account...
    </div>

    <div v-if="createdAccountId">
      <article class="message is-success">
        <div class="message-body">
          <fa-icon icon="check" /> Account Successfully Created!
        </div>
      </article>
      <div>
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
import ValidatedPasswordInput from '../components/validatedPasswordInput.vue';

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
      createdAccountId: null,
    };
  },
  components: {
    CenterCard,
    RecaptchaV2,
    ValidatedPasswordInput,
  },
  mounted() {
    this.$refs.email.focus();
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
  },
  methods: {
    submitClick() {
      this.error = null;

      const postData = {
        username: this.username,
        password: this.password,
        email: this.email,
      };

      if (this.appdata.recaptchav2.enabled) {
        postData.recaptchav2 = this.$refs.recaptchav2.getResponse();
        if (!postData.recaptchav2) {
          this.error = 'Please accept RECAPTCHA';
          return;
        }
      }

      this.loading = true;
      axios.post('api/ui/account', postData)
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
