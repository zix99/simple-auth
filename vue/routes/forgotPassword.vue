<template>
  <CenterCard title="Password Reset">
    <Message type="is-warning" v-if="!appdata.login.forgotPassword">
      Resetting password is disabled on this server.  Please reach out to your Administrator.
    </Message>

    <LoadingBanner :promise="loadingPromise" title="Password Reset">Requesting password reset...</LoadingBanner>

    <div v-if="appdata.login.forgotPassword && !sent">
      <p>Please enter your email address below, and if you have an account
      an email will be sent to you with instructions on how to reset your
      password</p>

      <div class="field">
        <label class="label">Email</label>
        <div class="control has-icons-left">
          <input class="input"
            ref="email"
            :class="{ 'is-danger': !validEmail }"
            type="email"
            placeholder="Email input"
            v-model="email" />
          <span class="icon is-small is-left">
            <i class="fas fa-envelope" />
          </span>
        </div>
      </div>
      <RecaptchaV2 v-if="appdata.recaptchav2.enabled"
        :sitekey="appdata.recaptchav2.sitekey"
        :theme="appdata.recaptchav2.theme"
        ref="recaptchav2" />
      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="submitClick" :disabled="!validEmail || (loadingPromise && loadingPromise.pending)">Submit</button>
        </div>
      </div>
    </div>

    <div v-if="sent">
      <Message type="is-success" v-if="sent">
        A reactivation email has been sent.
      </Message>
      <p>If you have an account associated with this email address, you have been sent an email with instructions on how to reset your password.</p>
    </div>
  </CenterCard>
</template>

<script>
import axios from 'axios';
import validator from 'validator';
import CenterCard from '../components/centerCard.vue';
import Message from '../components/message.vue';
import LoadingBanner from '../components/loadingBanner.vue';
import RecaptchaV2 from '../components/recaptchav2.vue';

export default {
  props: {
    appdata: null,
  },
  data() {
    return {
      email: '',
      sent: false,
      loadingPromise: null,
    };
  },
  mounted() {
    this.$refs.email.focus();
  },
  components: {
    CenterCard,
    Message,
    LoadingBanner,
    RecaptchaV2,
  },
  methods: {
    submitClick() {
      const data = {
        email: this.email,
      };
      const params = {};

      if (this.appdata.recaptchav2.enabled) {
        params.recaptchav2 = this.$refs.recaptchav2.getResponse();
        if (!params.recaptchav2) {
          return;
        }
      }

      this.loadingPromise = axios.post('api/ui/onetime', data, { params })
        .then(() => {
          this.sent = true;
        });
    },
  },
  computed: {
    validEmail() {
      return validator.isEmail(this.email);
    },
  },
};
</script>
