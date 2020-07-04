<template>
  <div class="wrapper" ref="recaptchav2"></div>
</template>

<script>
const recaptchaState = {
  state: {
    api: null,
  },
  setApi(api) {
    this.state.api = api;
  },
};

window.recaptchaV2OnLoad = function recaptchaLoad() {
  recaptchaState.setApi(window.grecaptcha);
};

let injectedVendorApi = false;
function injectVendorApi() {
  if (!injectedVendorApi) {
    injectedVendorApi = true;
    const script = document.createElement('script');
    script.setAttribute('src', 'https://www.google.com/recaptcha/api.js?render=explicit&onload=recaptchaV2OnLoad');
    script.setAttribute('async', true);
    script.setAttribute('defer', true);
    document.head.appendChild(script);
  }
}

export default {
  props: {
    sitekey: { type: String, required: true },
    theme: { type: String, required: false, default: 'light' },
  },
  data() {
    return {
      recaptcha: recaptchaState.state,
      setup: false,
    };
  },
  mounted() {
    // In case the recaptcha value is already set
    injectVendorApi();
    this.setupRecaptcha();
  },
  computed: {
    recaptchaApi() {
      return this.recaptcha.api;
    },
  },
  watch: {
    recaptchaApi() {
      this.setupRecaptcha();
    },
  },
  methods: {
    setupRecaptcha() {
      if (this.recaptcha.api && !this.setup) {
        this.setup = true;
        this.recaptcha.api.render(this.$refs.recaptchav2, {
          sitekey: this.sitekey,
          theme: this.theme,
        });
      }
    },
    getResponse() {
      if (this.setup) {
        return this.recaptcha.api.getResponse();
      }
      return null;
    },
  },
};
</script>

<style scoped>
.wrapper {
  margin: 16px 0;
}
</style>
