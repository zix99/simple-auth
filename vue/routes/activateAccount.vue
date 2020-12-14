<template>
  <CenterCard title="Account Activation">
    <LoadingBanner :promise="loadingPromise" :codes="errorCodes">
      <template v-slot:default>
        Activating your account...
      </template>
      <template v-slot:success>
        Your account has been activated. <router-link to="/">Login Now</router-link>
      </template>
      <template v-slot:error="{ error }">
        {{error}} <router-link to="/">Return Home</router-link>
      </template>
    </LoadingBanner>
  </CenterCard>
</template>

<script>
import axios from 'axios';
import LoadingBanner from '../components/loadingBanner.vue';
import CenterCard from '../components/centerCard.vue';

export default {
  components: {
    LoadingBanner,
    CenterCard,
  },
  props: {
    token: null,
    account: null,
  },
  data() {
    return {
      loadingPromise: null,
      activated: false,
      errorCodes: {
        'internal-error': 'There was an error activating your account',
        'no-code': 'The stipulation is no longer valid',
      },
    };
  },
  created() {
    this.loadingPromise = axios.post('api/v1/stipulation', { account: this.account, token: this.token });
  },
};
</script>
