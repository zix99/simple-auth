<template>
  <div>
    <LoadingBanner :promise="loadingPromise" title="Two Factor" :codes="errorCodes">
      Deactivating...
    </LoadingBanner>
    <p>Please enter 2FA Code to Deactivate</p>
    <div class="columns is-centered">
      <div class="column is-half">
        <div class="field is-grouped">
          <div class="control">
            <input class="input is-primary" type="text" placeholder="Enter 2FA Code" v-model="code" @keypress.enter="deactivate" v-focus>
          </div>
          <div class="control">
            <button class="button is-link is-danger" @click="deactivate">Deactivate</button>
          </div>
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
      loadingPromise: null,
      code: '',
      errorCodes: {
        'totp-failed': 'Invalid 2FA Code',
      },
    };
  },
  methods: {
    deactivate() {
      this.loadingPromise = axios.delete(`api/v1/local/2fa?code=${this.code}`)
        .then(() => {
          this.$emit('submitted');
        });
    },
  },
};
</script>
