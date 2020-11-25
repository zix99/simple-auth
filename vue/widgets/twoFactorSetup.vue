<template>
  <div>
    <LoadingBanner :promise="loadingPromise" title="Two Factor" :codes="errorCodes">
      Activating 2FA...
      <template v-slot:success>
        2FA activated!
      </template>
    </LoadingBanner>
    <div v-if="secret" class="qrcode">
      <img :src="`api/ui/2fa/qrcode?secret=${secret}`" />
      <span class="is-size-7">{{secret}}</span>
    </div>
    <div class="columns is-centered">
      <div class="column is-half">
        <div class="field is-grouped">
          <div class="control">
            <input class="input is-primary" type="text" placeholder="Enter 2FA Code" v-model="code" @keypress.enter="activate">
          </div>
          <div class="control">
            <button class="button is-link" @click="activate">Activate</button>
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
      secret: null,
      code: null,
      errorCodes: {
        'totp-failed': 'Invalid 2FA code. Please try again',
      },
    };
  },
  created() {
    axios.get('api/ui/2fa')
      .then((resp) => {
        this.secret = resp.data.secret;
      });
  },
  methods: {
    activate() {
      this.loadingPromise = axios.post('api/ui/2fa', { secret: this.secret, code: this.code })
        .then(() => {
          setTimeout(() => this.$emit('submitted'), 1500);
        });
    },
  },
};
</script>

<style scoped>
div.qrcode {
  text-align: center;
}
div.qrcode img {
  display: block;
  margin: 0 auto;
}
</style>
